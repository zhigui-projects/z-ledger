/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
)

var logger = flogging.MustGetLogger("stateormdb")

const savePointKey = "savepoint"

type SysState struct {
	ID         string
	Value      string
	VerAndMeta string `json:"-"`
}

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBInstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	//openCounts    uint64
	cache              *statedb.Cache
	redoLoggerProvider *redoLoggerProvider
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	ormDBInstance *ormdb.ORMDBInstance
	metadataDB    *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName     string                        // The name of the chain/channel.
	namespaceDBs  map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	redoLog       *redoLogger
	mux           sync.RWMutex
	cache         *statedb.Cache
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider(config *ormdbconfig.ORMDBConfig, metricsProvider metrics.Provider, cache *statedb.Cache) (*VersionedDBProvider, error) {
	logger.Debug("constructing ORMDB VersionedDBProvider")
	instance, err := ormdb.NewORMDBInstance(config, metricsProvider)
	if err != nil {
		logger.Errorf("create ormdb instance failed [%v]", err)
		return nil, errors.WithMessage(err, "create ormdb instance failed")
	}
	p, err := newRedoLoggerProvider(config.RedoLogPath)
	if err != nil {
		logger.Errorf("create redologger failed [%v]", err)
		return nil, errors.WithMessage(err, "create redo log failed")
	}
	return &VersionedDBProvider{ormDBInstance: instance, databases: make(map[string]*VersionedDB), redoLoggerProvider: p, cache: cache}, nil
}

// GetDBHandle gets the handle to a named database
func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()
	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(provider.ormDBInstance, dbName, provider.redoLoggerProvider.newRedoLogger(dbName), provider.cache)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// Close closes the underlying db instance
func (provider *VersionedDBProvider) Close() {
	provider.redoLoggerProvider.close()
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(ormDBInstance *ormdb.ORMDBInstance, dbName string, redoLog *redoLogger, cache *statedb.Cache) (*VersionedDB, error) {
	chainName := dbName
	metaDBName := dbName + "_meta"

	metadataDB, err := ormdb.CreateORMDatabase(ormDBInstance, metaDBName)
	if err != nil {
		logger.Errorf("create meta database failed [%v]", err)
		return nil, errors.WithMessage(err, "create meta database failed")
	}
	savepoint := &SysState{}
	if !metadataDB.DB.HasTable(savepoint) {
		err = metadataDB.DB.CreateTable(savepoint).Error
		if err != nil {
			logger.Errorf("create meta table failed [%v]", err)
			return nil, errors.WithMessage(err, "create meta table failed")
		}
	}
	namespaceDBMap := make(map[string]*ormdb.ORMDatabase)
	vdb := &VersionedDB{
		ormDBInstance: ormDBInstance,
		metadataDB:    metadataDB,
		chainName:     chainName,
		namespaceDBs:  namespaceDBMap,
		redoLog:       redoLog,
		cache:         cache,
	}

	logger.Debugf("chain [%s]: checking for redolog record", chainName)
	redologRecord, err := redoLog.load()
	if err != nil {
		logger.Errorf("redolog load failed [%v]", err)
		return nil, errors.WithMessage(err, "redolog load failed")
	}
	latestSavepoint, err := vdb.GetLatestSavePoint()
	if err != nil {
		logger.Errorf("get latest savepoint failed [%v]", err)
		return nil, errors.WithMessage(err, "get latest savepoint failed")
	}

	// in normal circumstances, redolog is expected to be either equal to the last block
	// committed to the statedb or one ahead (in the event of a crash). However, either of
	// these or both could be nil on first time start (fresh start/rebuild)
	if redologRecord == nil || latestSavepoint == nil {
		logger.Debugf("chain [%s]: No redo-record or save point present", chainName)
		return vdb, nil
	}

	logger.Debugf("chain [%s]: save point = %#v, version of redolog record = %#v",
		chainName, savepoint, redologRecord.Version)

	if redologRecord.Version.BlockNum-latestSavepoint.BlockNum == 1 {
		logger.Debugf("chain [%s]: Re-applying last batch", chainName)
		if err := vdb.applyUpdates(redologRecord.UpdateBatch, redologRecord.Version); err != nil {
			return nil, err
		}
	}

	return vdb, nil
}

// GetState implements method in VersionedDB interface
func (v *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	logger.Debugf("GetState(). ns=%s, key=%s", namespace, key)

	// (1) read the KV from the cache if available
	cacheEnabled := v.cache.Enabled(namespace)
	if cacheEnabled {
		cv, err := v.cache.GetState(v.chainName, namespace, key)
		if err != nil {
			logger.Errorf("get state from cache failed [%v]", err)
			return nil, errors.WithMessage(err, "get state from cache failed")
		}
		if cv != nil {
			vv, err := constructVersionedValue(cv)
			if err != nil {
				logger.Errorf("construct versioned value from cache failed [%v]", err)
				return nil, errors.WithMessage(err, "construct versioned value from cache failed")
			}
			return vv, nil
		}
	}

	// (2) read from the database if cache miss occurs
	vv, err := v.readFromDB(namespace, key)
	if err != nil {
		logger.Errorf("get state from database failed [%v]", err)
		return nil, errors.WithMessage(err, "get state from database failed")
	}
	if vv == nil {
		return nil, nil
	}

	// (3) if the value is not nil, store in the cache
	if cacheEnabled {
		cacheValue := constructCacheValue(vv)
		if err := v.cache.PutState(v.chainName, namespace, key, cacheValue); err != nil {
			logger.Errorf("put state to cache failed [%v]", err)
			return nil, errors.WithMessage(err, "put state to cache failed")
		}
	}

	return vv, nil
}

func (v *VersionedDB) readFromDB(namespace, key string) (*statedb.VersionedValue, error) {
	db, err := v.getNamespaceDBHandle(namespace)
	if err != nil {
		logger.Errorf("get namespaced database failed [%v]", err)
		return nil, errors.WithMessage(err, "get namespaced database failed")
	}

	keys, err, isORMKey := parseKey(key)
	if err != nil {
		logger.Errorf("parse ormdb key failed [%v]", err)
		return nil, errors.WithMessage(err, "parse ormdb key failed")
	}

	if !isORMKey {
		sysState := &SysState{ID: key}
		if err = db.DB.Find(sysState).Error; err != nil {
			return nil, errors.WithMessage(err, "query sys state failed")
		}

		returnVersion, returnMetadata, err := DecodeVersionAndMetadata(sysState.VerAndMeta)
		if err != nil {
			logger.Errorf("marshal version and meta failed [%v]", err)
			return nil, errors.WithMessage(err, "marshal version and meta failed")
		}
		value, err := base64.StdEncoding.DecodeString(sysState.Value)
		if err != nil {
			logger.Errorf("decode sys state value failed [%v]", err)
			return nil, errors.WithMessage(err, "decode sys state value failed")
		}
		return &statedb.VersionedValue{Version: returnVersion, Metadata: returnMetadata, Value: value}, nil
	} else {
		entityName := keys[0]
		entityId := keys[1]

		db.RWMutex.RLock()
		entity := db.ModelTypes[entityName].Interface()
		id, exist := db.ModelTypes[entityName].FieldByName("ID")
		if !exist {
			return nil, errors.New("entity no ID field")
		}
		_, exist = db.ModelTypes[entityName].FieldByName("VerAndMeta")
		if !exist {
			return nil, errors.New("entity no VerAndMeta field")
		}
		db.RWMutex.RUnlock()

		if id.Type.Kind() == reflect.String {
			if err = db.DB.Table(ormdb.ToTableName(entityName)).Where("id = ?", entityId).Find(entity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, nil
				}
				return nil, errors.WithMessage(err, "query entity failed")
			}
		} else {
			return nil, errors.New("not supported entity ID field type")
		}

		verAndMeta := reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").String()
		returnVersion, returnMetadata, err := DecodeVersionAndMetadata(verAndMeta)
		if err != nil {
			logger.Errorf("marshal version and meta failed [%v]", err)
			return nil, errors.WithMessage(err, "marshal version and meta failed")
		}

		entityBytes, err := json.Marshal(entity)
		if err != nil {
			logger.Errorf("marshal entity failed [%v]", err)
			return nil, errors.WithMessage(err, "marshal entity failed")
		}
		return &statedb.VersionedValue{Version: returnVersion, Metadata: returnMetadata, Value: entityBytes}, nil
	}

}

// GetVersion implements method in VersionedDB interface
func (v *VersionedDB) GetVersion(namespace string, key string) (*version.Height, error) {
	versionedValue, err := v.GetState(namespace, key)
	if err != nil {
		logger.Errorf("get state failed [%v]", err)
		return nil, errors.WithMessage(err, "get state failed")
	}
	if versionedValue == nil {
		return nil, nil
	}
	return versionedValue.Version, nil
}

// GetStateMultipleKeys implements method in VersionedDB interface
func (v *VersionedDB) GetStateMultipleKeys(namespace string, keys []string) ([]*statedb.VersionedValue, error) {
	vals := make([]*statedb.VersionedValue, len(keys))
	for i, key := range keys {
		val, err := v.GetState(namespace, key)
		if err != nil {
			logger.Errorf("get state failed [%v]", err)
			return nil, errors.WithMessage(err, "get state failed")
		}
		vals[i] = val
	}
	return vals, nil
}

// GetStateRangeScanIterator implements method in VersionedDB interface
func (v *VersionedDB) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (statedb.ResultsIterator, error) {
	return v.GetStateRangeScanIteratorWithMetadata(namespace, startKey, endKey, nil)
}

type kvScanner struct {
	namespace            string
	dbItr                *sql.Rows
	requestedLimit       int32
	totalRecordsReturned int32
	entityType           reflect.Type
	isSysState           bool
}

func newKVScanner(namespace string, dbItr *sql.Rows, requestedLimit int32, entityType reflect.Type, isSysState bool) *kvScanner {
	return &kvScanner{namespace, dbItr, requestedLimit, 0, entityType, isSysState}
}

func (scanner *kvScanner) Next() (statedb.QueryResult, error) {
	if scanner.requestedLimit > 0 && scanner.totalRecordsReturned >= scanner.requestedLimit {
		return nil, nil
	}
	if !scanner.dbItr.Next() {
		return nil, nil
	}

	if scanner.isSysState {
		sysState := &SysState{}
		if err := scanner.dbItr.Scan(sysState); err != nil {
			return nil, err
		}

		ver, meta, err := DecodeVersionAndMetadata(sysState.VerAndMeta)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(sysState)
		if err != nil {
			return nil, err
		}

		scanner.totalRecordsReturned++
		return &statedb.VersionedKV{
			CompositeKey: statedb.CompositeKey{Namespace: scanner.namespace, Key: sysState.ID},
			// TODO remove dereferrencing below by changing the type of the field
			// `VersionedValue` in `statedb.VersionedKV` to a pointer
			VersionedValue: statedb.VersionedValue{Value: value, Metadata: meta, Version: ver}}, nil
	} else {
		id, exist := scanner.entityType.FieldByName("ID")
		if !exist {
			return nil, errors.New("entity no ID field")
		}
		_, exist = scanner.entityType.FieldByName("VerAndMeta")
		if !exist {
			return nil, errors.New("entity no VerAndMeta field")
		}
		entity := reflect.New(scanner.entityType).Interface()
		if id.Type.Kind() == reflect.String {
			if err := scanner.dbItr.Scan(entity); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("not supported entity ID field type")
		}

		key := reflect.ValueOf(entity).Elem().FieldByName("ID").String()
		verAndMeta := reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").String()
		returnVersion, returnMetadata, err := DecodeVersionAndMetadata(verAndMeta)
		if err != nil {
			logger.Errorf("marshal version and meta failed [%v]", err)
			return nil, errors.WithMessage(err, "marshal version and meta failed")
		}

		entityBytes, err := json.Marshal(entity)
		if err != nil {
			logger.Errorf("marshal entity failed [%v]", err)
			return nil, errors.WithMessage(err, "marshal entity failed")
		}
		scanner.totalRecordsReturned++
		return &statedb.VersionedKV{
			CompositeKey: statedb.CompositeKey{Namespace: scanner.namespace, Key: key},
			// TODO remove dereferrencing below by changing the type of the field
			// `VersionedValue` in `statedb.VersionedKV` to a pointer
			VersionedValue: statedb.VersionedValue{Value: entityBytes, Metadata: returnMetadata, Version: returnVersion}}, nil
	}
}

func (scanner *kvScanner) Close() {
	err := scanner.dbItr.Close()
	if err != nil {
		logger.Errorf("close scanner failed [%v]", err)
	}
}

func (scanner *kvScanner) GetBookmarkAndClose() string {
	retval := ""

	if !scanner.dbItr.Next() {
		return retval
	}

	if scanner.isSysState {
		sysState := &SysState{}
		if err := scanner.dbItr.Scan(sysState); err != nil {
			logger.Errorf("scan sys state failed [%v]", err)
			return retval
		}

		retval = sysState.ID
	} else {
		id, exist := scanner.entityType.FieldByName("ID")
		if !exist {
			logger.Error("entity no ID field")
			return retval
		}
		_, exist = scanner.entityType.FieldByName("VerAndMeta")
		if !exist {
			logger.Error("entity no VerAndMeta field")
			return retval
		}
		entity := reflect.New(scanner.entityType).Interface()
		if id.Type.Kind() == reflect.String {
			if err := scanner.dbItr.Scan(entity); err != nil {
				logger.Errorf("scan entity failed [%v]", err)
				return retval
			}
		} else {
			logger.Error("not supported entity ID field type")
			return retval
		}

		key := reflect.ValueOf(entity).Elem().FieldByName("ID").String()
		retval = key
	}
	scanner.Close()
	return retval
}

// GetStateRangeScanIteratorWithMetadata implements method in VersionedDB interface
func (v *VersionedDB) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	db, err := v.getNamespaceDBHandle(namespace)
	if err != nil {
		logger.Errorf("get namespaced database failed [%v]", err)
		return nil, errors.WithMessage(err, "get namespaced database failed")
	}

	startKeys, err, isORMKey := parseKey(startKey)
	if err != nil {
		logger.Errorf("parse ormdb start key failed [%v]", err)
		return nil, errors.WithMessage(err, "parse ormdb start key failed")
	}

	endKeys, err, _ := parseKey(endKey)
	if err != nil {
		logger.Errorf("parse ormdb end key failed [%v]", err)
		return nil, errors.WithMessage(err, "parse ormdb end key failed")
	}

	if !isORMKey {
		rows, err := db.DB.Table(ormdb.ToTableName("SysState")).Where("id >= ? AND id < ?", startKey, endKey).Rows()
		if err != nil {
			return nil, errors.WithMessage(err, "rang query sys state failed")
		}

		return &kvScanner{isSysState: true, dbItr: rows, requestedLimit: 0, totalRecordsReturned: 0, entityType: reflect.TypeOf(&SysState{}), namespace: namespace}, nil
	} else {
		entityName := startKeys[0]
		entityStartId := startKeys[1]
		entityEndId := endKeys[1]

		db.RWMutex.RLock()
		id, exist := db.ModelTypes[entityName].FieldByName("ID")
		if !exist {
			return nil, errors.New("entity no ID field")
		}
		_, exist = db.ModelTypes[entityName].FieldByName("VerAndMeta")
		if !exist {
			return nil, errors.New("entity no VerAndMeta field")
		}
		db.RWMutex.RUnlock()

		if id.Type.Kind() == reflect.String {
			rows, err := db.DB.Table(ormdb.ToTableName(entityName)).Where("id >= ? AND id < ?", entityStartId, entityEndId).Rows()
			if err != nil {
				return nil, errors.WithMessage(err, "rang query sys state failed")
			}

			return &kvScanner{isSysState: false, dbItr: rows, requestedLimit: 0, totalRecordsReturned: 0, entityType: db.ModelTypes[entityName].StructType(), namespace: namespace}, nil
		} else {
			return nil, errors.New("not supported entity ID field type")
		}
	}

}

// ExecuteQuery implements method in VersionedDB interface
func (v *VersionedDB) ExecuteQuery(namespace, query string) (statedb.ResultsIterator, error) {
	return nil, errors.New("ExecuteQuery not supported for ormdb")
}

// ExecuteQueryWithMetadata implements method in VersionedDB interface
func (v *VersionedDB) ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	return nil, errors.New("ExecuteQueryWithMetadata not supported for ormdb")
}

// ApplyUpdates implements method in VersionedDB interface
func (v *VersionedDB) ApplyUpdates(updates *statedb.UpdateBatch, height *version.Height) error {
	if height != nil && updates.ContainsPostOrderWrites {
		// height is passed nil when committing missing private data for previously committed blocks
		r := &redoRecord{
			UpdateBatch: updates,
			Version:     height,
		}
		if err := v.redoLog.persist(r); err != nil {
			return err
		}
	}
	return v.applyUpdates(updates, height)
}

// applyUpdates apply write set to state db
func (v *VersionedDB) applyUpdates(updates *statedb.UpdateBatch, height *version.Height) error {
	// stage 1 - buildCommitters builds committers per namespace (per DB). Each committer transforms the
	// given batch in the form of underlying db and keep it in memory.
	committers, err := v.buildCommitters(updates)
	if err != nil {
		return err
	}

	// stage 2 -- executeCommitter executes each committer to push the changes to the DB
	if err = v.executeCommitter(committers); err != nil {
		return err
	}

	// Stgae 3 - postCommitProcessing - flush and record savepoint.
	namespaces := updates.GetUpdatedNamespaces()
	if err := v.postCommitProcessing(committers, namespaces, height); err != nil {
		return err
	}

	return nil
}

func (v *VersionedDB) postCommitProcessing(committers []*committer, namespaces []string, height *version.Height) error {
	var wg sync.WaitGroup

	wg.Add(1)
	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		defer wg.Done()

		cacheUpdates := make(statedb.CacheUpdates)
		for _, c := range committers {
			if !c.cacheEnabled {
				continue
			}
			cacheUpdates.Add(c.namespace, c.cacheKVs)
		}

		if len(cacheUpdates) == 0 {
			return
		}

		// update the cache
		if err := v.cache.UpdateStates(v.chainName, cacheUpdates); err != nil {
			v.cache.Reset()
			errChan <- err
		}

	}()

	// Record a savepoint at a given height
	if err := v.recordSavepoint(height); err != nil {
		logger.Errorf("Error during recordSavepoint: %s", err.Error())
		return err
	}

	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (v *VersionedDB) recordSavepoint(height *version.Height) error {
	if height == nil {
		return nil
	}

	savePoint := &SysState{ID: savePointKey, Value: base64.StdEncoding.EncodeToString(height.ToBytes())}
	err := v.metadataDB.DB.Save(savePoint).Error
	if err != nil {
		logger.Errorf("Failed to save the savepoint to DB %s", err.Error())
		return err
	}

	return nil
}

func (v *VersionedDB) GetLatestSavePoint() (*version.Height, error) {
	sp := &SysState{}
	err := v.metadataDB.DB.Where(&SysState{ID: savePointKey}).Find(sp).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		} else {
			logger.Errorf("Failed to read savepoint data %s", err.Error())
			return nil, err
		}
	}

	versionBytes, err := base64.StdEncoding.DecodeString(sp.Value)
	if err != nil {
		logger.Errorf("Failed to decode savepoint data %s", err.Error())
		return nil, err
	}
	height, _, err := version.NewHeightFromBytes(versionBytes)
	return height, err
}

func (v *VersionedDB) ValidateKeyValue(key string, value []byte) error {
	return nil
}

func (v *VersionedDB) BytesKeySupported() bool {
	return false
}

func (v *VersionedDB) Open() error {
	return nil
}

func (v *VersionedDB) Close() {

}

func (v *VersionedDB) ExecuteConditionQuery(namespace string, search entitydefinition.Search) (interface{}, error) {
	db, err := v.getNamespaceDBHandle(namespace)
	if err != nil {
		logger.Errorf("get namespaced database failed [%v]", err)
		return nil, errors.WithMessage(err, "get namespaced database failed")
	}

	entityName := search.Entity
	gormdb := db.DB
	for _, cond := range search.WhereConditions {
		query := string(cond["query"][0])
		argsBytes := cond["args"]
		args, err := entitydefinition.DecodeSearchValues(argsBytes)
		if err != nil {
			return nil, errors.WithMessage(err, "decode search where args failed")
		}
		gormdb = gormdb.Where(query, args...)
	}

	for _, cond := range search.OrConditions {
		query := string(cond["query"][0])
		argsBytes := cond["args"]
		args, err := entitydefinition.DecodeSearchValues(argsBytes)
		if err != nil {
			return nil, errors.WithMessage(err, "decode search or args failed")
		}
		gormdb = gormdb.Or(query, args...)
	}

	for _, cond := range search.NotConditions {
		query := string(cond["query"][0])
		argsBytes := cond["args"]
		args, err := entitydefinition.DecodeSearchValues(argsBytes)
		if err != nil {
			return nil, errors.WithMessage(err, "decode search not args failed")
		}
		gormdb = gormdb.Not(query, args...)
	}

	for _, order := range search.OrderConditions {
		gormdb = gormdb.Order(order)
	}

	//TODO:make it configurable
	maxLimit := 30
	var limit int
	if search.LimitCondition > maxLimit {
		limit = maxLimit
	} else {
		limit = search.LimitCondition
	}
	gormdb = gormdb.Offset(search.OffsetCondition).Limit(limit)

	db.RWMutex.RLock()
	models := reflect.New(reflect.SliceOf(db.ModelTypes[entityName].StructType())).Interface()
	db.RWMutex.RUnlock()
	err = gormdb.Table(ormdb.ToTableName(entityName)).Find(models).Error
	if err != nil {
		return nil, errors.WithMessage(err, "condition query failed")
	}

	return models, nil
}

// getNamespaceDBHandle gets the handle to a named chaincode database
func (v *VersionedDB) getNamespaceDBHandle(namespace string) (*ormdb.ORMDatabase, error) {
	v.mux.RLock()
	db := v.namespaceDBs[namespace]
	v.mux.RUnlock()
	if db != nil {
		return db, nil
	}
	namespaceDBName := fmt.Sprintf("%s_%s", v.chainName, namespace)
	v.mux.Lock()
	defer v.mux.Unlock()
	db = v.namespaceDBs[namespace]
	if db == nil {
		var err error
		db, err = ormdb.CreateORMDatabase(v.ormDBInstance, namespaceDBName)
		if err != nil {
			logger.Errorf("create orm database failed [%v]", err)
			return nil, errors.WithMessage(err, "create orm database failed")
		}
		v.namespaceDBs[namespace] = db
		sysState := &SysState{}
		if !db.DB.HasTable(sysState) {
			err = db.DB.CreateTable(sysState).Error
			if err != nil {
				logger.Errorf("create sys state table failed [%v]", err)
				return nil, errors.WithMessage(err, "create sys state table failed")
			}
		}
		edf := &entitydefinition.EntityFieldDefinition{}
		if !db.DB.HasTable(edf) {
			err = db.DB.CreateTable(edf).Error
			if err != nil {
				logger.Errorf("create entity definition table failed [%v]", err)
				return nil, errors.WithMessage(err, "create entity definition table failed")
			}
		}

		var entityFieldDefs []entitydefinition.EntityFieldDefinition
		err = db.DB.Order("seq").Find(&entityFieldDefs).Error
		if err != nil {
			logger.Errorf("query entity definition table failed [%v]", err)
			return nil, errors.WithMessage(err, "query entity definition table failed")
		}

		var currentEntity string
		currentEfds := make([]entitydefinition.EntityFieldDefinition, 0)
		for i, efd := range entityFieldDefs {
			if i == 0 {
				currentEntity = efd.Owner
				currentEfds = append(currentEfds, efd)
			} else {
				if currentEntity != efd.Owner {
					db.RWMutex.Lock()
					ds := entitydefinition.NewBuilder().AddEntityFieldDefinition(currentEfds, db.ModelTypes).Build()
					db.ModelTypes[currentEntity] = ds
					db.RWMutex.Unlock()
					entity := ds.Interface()
					tableName := ormdb.ToTableName(currentEntity)
					if !db.DB.HasTable(tableName) {
						db.DB.Table(tableName).CreateTable(entity)
					}
					currentEntity = efd.Owner
					currentEfds = make([]entitydefinition.EntityFieldDefinition, 0)
					currentEfds = append(currentEfds, efd)
				} else {
					currentEfds = append(currentEfds, efd)
				}
			}

			if i == len(entityFieldDefs)-1 {
				db.RWMutex.Lock()
				ds := entitydefinition.NewBuilder().AddEntityFieldDefinition(currentEfds, db.ModelTypes).Build()
				db.ModelTypes[currentEntity] = ds
				db.RWMutex.Unlock()
			}
		}

	}
	return db, nil
}

func constructVersionedValue(cv *statedb.CacheValue) (*statedb.VersionedValue, error) {
	height, _, err := version.NewHeightFromBytes(cv.VersionBytes)
	if err != nil {
		return nil, err
	}

	return &statedb.VersionedValue{
		Value:    cv.Value,
		Version:  height,
		Metadata: cv.Metadata,
	}, nil
}

func constructCacheValue(v *statedb.VersionedValue) *statedb.CacheValue {
	return &statedb.CacheValue{
		VersionBytes: v.Version.ToBytes(),
		Value:        v.Value,
		Metadata:     v.Metadata,
	}
}

func parseKey(key string) ([]string, error, bool) {
	if !strings.Contains(key, entitydefinition.ORMDB_SEPERATOR) {
		return nil, nil, false
	}
	keys := strings.Split(key, entitydefinition.ORMDB_SEPERATOR)
	if len(keys) < 2 {
		logger.Errorf("wrong ormdb key [%s]", key)
		return nil, errors.New("wrong ormdb key"), true
	}
	return keys, nil, true
}
