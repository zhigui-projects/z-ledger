/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
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

type savePoint struct {
	key    string
	height string
}

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBInstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	//openCounts    uint64
	cache              *statedb.Cache
	redoLoggerProvider *statecouchdb.RedoLoggerProvider
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	ormDBInstance *ormdb.ORMDBInstance
	metadataDB    *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName     string                        // The name of the chain/channel.
	namespaceDBs  map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	redoLog       *statecouchdb.RedoLogger
	mux           sync.RWMutex
	cache         *statedb.Cache
}

// ormDBSavepointData data for couchdb
type ormDBSavepointData struct {
	BlockNum uint64 `json:"BlockNum"`
	TxNum    uint64 `json:"TxNum"`
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider(config *ormdbconfig.ORMDBConfig, metricsProvider metrics.Provider, cache *statedb.Cache) (*VersionedDBProvider, error) {
	logger.Debugf("constructing ORMDB VersionedDBProvider")
	instance, err := ormdb.NewORMDBInstance(config, metricsProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "create ormdb instance failed")
	}
	p, err := statecouchdb.NewRedoLoggerProvider(config.RedoLogPath)
	if err != nil {
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
		vdb, err = newVersionedDB(provider.ormDBInstance, dbName, provider.redoLoggerProvider.NewRedoLogger(dbName), provider.cache)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(ormDBInstance *ormdb.ORMDBInstance, dbName string, redoLog *statecouchdb.RedoLogger, cache *statedb.Cache) (*VersionedDB, error) {
	chainName := dbName
	dbName = dbName + "_"

	metadataDB, err := ormdb.CreateORMDatabase(ormDBInstance, dbName)
	if err != nil {
		return nil, err
	}
	err = metadataDB.DB.CreateTable(&savePoint{}).Error
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*ormdb.ORMDatabase)
	return &VersionedDB{
		ormDBInstance: ormDBInstance,
		metadataDB:    metadataDB,
		chainName:     chainName,
		namespaceDBs:  namespaceDBMap,
		redoLog:       redoLog,
		cache:         cache,
	}, nil
}

// GetState implements method in VersionedDB interface
func (v *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	logger.Debugf("GetState(). ns=%s, key=%s", namespace, key)

	// (1) read the KV from the cache if available
	cacheEnabled := v.cache.Enabled(namespace)
	if cacheEnabled {
		cv, err := v.cache.GetState(v.chainName, namespace, key)
		if err != nil {
			return nil, err
		}
		if cv != nil {
			vv, err := constructVersionedValue(cv)
			if err != nil {
				return nil, err
			}
			return vv, nil
		}
	}

	// (2) read from the database if cache miss occurs
	vv, err := v.readFromDB(namespace, key)
	if err != nil {
		return nil, err
	}
	if vv == nil {
		return nil, nil
	}

	// (3) if the value is not nil, store in the cache
	if cacheEnabled {
		cacheValue := constructCacheValue(vv)
		if err := v.cache.PutState(v.chainName, namespace, key, cacheValue); err != nil {
			return nil, err
		}
	}

	return vv, nil
}

func (v *VersionedDB) readFromDB(namespace, key string) (*statedb.VersionedValue, error) {
	db, err := v.getNamespaceDBHandle(namespace)
	if err != nil {
		return nil, err
	}

	keys, err := parseKey(key)
	if err != nil {
		return nil, err
	}

	entityName := keys[0]
	entityId := keys[1]
	entity := db.ModelTypes[entityName].Interface()
	id, exist := db.ModelTypes[entityName].FieldByName("ID")
	if !exist {
		return nil, errors.New("entity no ID field")
	}
	if id.Type.Kind() == reflect.Int {
		entityIdInt, err := strconv.Atoi(entityId)
		if err != nil {
			return nil, errors.WithMessage(err, "entity int ID convert failed")
		}
		db.DB.Find(entity, entityIdInt)
	} else if id.Type.Kind() == reflect.String {
		db.DB.Find(entity, entityId)
	} else {
		return nil, errors.New("not supported entity ID field type")
	}

	if entity == nil {
		return nil, nil
	}

	id, exist = db.ModelTypes[entityName].FieldByName("VerAndMeta")
	if !exist {
		return nil, errors.New("entity no VerAndMeta field")
	}
	verAndMeta := reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").String()
	returnVersion, returnMetadata, err := statecouchdb.DecodeVersionAndMetadata(verAndMeta)

	entityBytes, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	return &statedb.VersionedValue{Version: returnVersion, Metadata: returnMetadata, Value: entityBytes}, nil
}

// GetVersion implements method in VersionedDB interface
func (v *VersionedDB) GetVersion(namespace string, key string) (*version.Height, error) {
	versionedValue, err := v.GetState(namespace, key)
	if err != nil {
		return nil, err
	}
	if versionedValue == nil {
		return nil, nil
	}
	return versionedValue.Version, nil
}

func (v *VersionedDB) GetStateMultipleKeys(namespace string, keys []string) ([]*statedb.VersionedValue, error) {
	panic("implement me")
}

func (v *VersionedDB) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (statedb.ResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) ExecuteQuery(namespace, query string) (statedb.ResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	panic("implement me")
}

// ApplyUpdates implements method in VersionedDB interface
func (v *VersionedDB) ApplyUpdates(updates *statedb.UpdateBatch, height *version.Height) error {
	if height != nil && updates.ContainsPostOrderWrites {
		// height is passed nil when committing missing private data for previously committed blocks
		r := &statecouchdb.RedoRecord{
			UpdateBatch: updates,
			Version:     height,
		}
		if err := v.redoLog.Persist(r); err != nil {
			return err
		}
	}
	return v.applyUpdates(updates, height)
}

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

func (vdb *VersionedDB) recordSavepoint(height *version.Height) error {
	if height == nil {
		return nil
	}

	savePoint := &savePoint{key: savePointKey, height: base64.StdEncoding.EncodeToString(height.ToBytes())}
	err := vdb.metadataDB.DB.Save(savePoint).Error

	if err != nil {
		logger.Errorf("Failed to save the savepoint to DB %s", err.Error())
		return err
	}

	return nil
}

func (v *VersionedDB) GetLatestSavePoint() (*version.Height, error) {
	panic("implement me")
}

func (v *VersionedDB) ValidateKeyValue(key string, value []byte) error {
	panic("implement me")
}

func (v *VersionedDB) BytesKeySupported() bool {
	panic("implement me")
}

func (v *VersionedDB) Open() error {
	panic("implement me")
}

func (v *VersionedDB) Close() {
	panic("implement me")
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
		db, err := ormdb.CreateORMDatabase(v.ormDBInstance, namespaceDBName)
		if err != nil {
			return nil, err
		}
		db.DB.CreateTable(&entitydefinition.EntityFieldDefinition{})
		v.namespaceDBs[namespace] = db
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

func parseKey(key string) ([]string, error) {
	keys := strings.Split(key, entitydefinition.ORMDB_SEPERATOR)
	if len(keys) != 2 {
		return nil, errors.New("wrong ormdb key")
	}
	return keys, nil
}
