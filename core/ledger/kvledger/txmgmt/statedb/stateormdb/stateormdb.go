/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/msgs"
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

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBInstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	//openCounts    uint64
	cache *statedb.Cache
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	ormDBInstance *ormdb.ORMDBInstance
	metadataDB    *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName     string                        // The name of the chain/channel.
	namespaceDBs  map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	mux           sync.RWMutex
	cache         *statedb.Cache
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider(config *ormdbconfig.ORMDBConfig, metricsProvider metrics.Provider, cache *statedb.Cache) (*VersionedDBProvider, error) {
	logger.Debugf("constructing ORMDB VersionedDBProvider")
	instance, err := ormdb.NewORMDBInstance(config, metricsProvider)
	if err != nil {
		return nil, err
	}
	return &VersionedDBProvider{ormDBInstance: instance, databases: make(map[string]*VersionedDB), cache: cache}, nil
}

// GetDBHandle gets the handle to a named database
func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()
	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(provider.ormDBInstance, dbName, provider.cache)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(ormDBInstance *ormdb.ORMDBInstance, dbName string, cache *statedb.Cache) (*VersionedDB, error) {
	chainName := dbName
	dbName = dbName + "_"

	metadataDB, err := ormdb.CreateORMDatabase(ormDBInstance, dbName)
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*ormdb.ORMDatabase)
	return &VersionedDB{
		ormDBInstance: ormDBInstance,
		metadataDB:    metadataDB,
		chainName:     chainName,
		namespaceDBs:  namespaceDBMap,
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
	verAndMetaft, exist := db.ModelTypes[entityName].FieldByName("ID")
	if !exist {
		return nil, errors.New("entity no ID field")
	}
	if verAndMetaft.Type.Kind() == reflect.Int {
		entityIdInt, err := strconv.Atoi(entityId)
		if err != nil {
			return nil, errors.WithMessage(err, "entity int ID convert failed")
		}
		db.DB.Find(entity, entityIdInt)
	} else if verAndMetaft.Type.Kind() == reflect.String {
		db.DB.Find(entity, entityId)
	} else {
		return nil, errors.New("not supported entity ID field type")
	}

	if entity == nil {
		return nil, nil
	}

	verAndMetaft, exist = db.ModelTypes[entityName].FieldByName("VerAndMeta")
	if !exist {
		return nil, errors.New("entity no VerAndMeta field")
	}
	verAndMeta := reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").String()
	returnVersion, returnMetadata, err := msgs.DecodeVersionAndMetadata(verAndMeta)

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

func (v *VersionedDB) ApplyUpdates(batch *statedb.UpdateBatch, height *version.Height) error {
	panic("implement me")
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

//func modelToKeyValue() (*statedb.VersionedValue, error) {
//	// initialize the return value
//	var returnValue []byte
//	var err error
//	reflect.Type()
//	// create a generic map unmarshal the json
//	var valueBytes bytes.Buffer
//	encoder:=gob.NewEncoder(&valueBytes)
//	encoder.EncodeValue()
//	decoder := json.NewDecoder(bytes.NewBuffer(doc.JSONValue))
//	decoder.UseNumber()
//	if err = decoder.Decode(&jsonResult); err != nil {
//		return nil, err
//	}
//	// verify the version field exists
//	if _, fieldFound := jsonResult[versionField]; !fieldFound {
//		return nil, errors.Errorf("version field %s was not found", versionField)
//	}
//	key := jsonResult[idField].(string)
//	// create the return version from the version field in the JSON
//
//	returnVersion, returnMetadata, err := decodeVersionAndMetadata(jsonResult[versionField].(string))
//	if err != nil {
//		return nil, err
//	}
//	// remove the _id, _rev and version fields
//	delete(jsonResult, idField)
//	delete(jsonResult, revField)
//	delete(jsonResult, versionField)
//
//	// handle binary or json data
//	if doc.Attachments != nil { // binary attachment
//		// get binary data from attachment
//		for _, attachment := range doc.Attachments {
//			if attachment.Name == binaryWrapper {
//				returnValue = attachment.AttachmentBytes
//			}
//		}
//	} else {
//		// marshal the returned JSON data.
//		if returnValue, err = json.Marshal(jsonResult); err != nil {
//			return nil, err
//		}
//	}
//	return &keyValue{key, &statedb.VersionedValue{
//		Value:    returnValue,
//		Metadata: returnMetadata,
//		Version:  returnVersion},
//	}, nil
//}
