/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/spf13/viper"
)

var logger = flogging.MustGetLogger("stateormdb")

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBinstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	//openCounts    uint64
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	ormDBInstance  *ormdb.ORMDBInstance
	metadataDB     *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName      string                        // The name of the chain/channel.
	namespaceDBs   map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	mux            sync.RWMutex
	lsccStateCache *statedb.LsccStateCache
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider(metricsProvider metrics.Provider) (*VersionedDBProvider, error) {
	logger.Debugf("constructing ORMDB VersionedDBProvider")
	instance, err := ormdb.NewORMDBInstance(metricsProvider)
	if err != nil {
		return nil, err
	}
	return &VersionedDBProvider{ormDBinstance: instance, databases: make(map[string]*VersionedDB)}, nil
}

// GetDBHandle gets the handle to a named database
func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()
	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(provider.ormDBinstance, dbName)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(ormDBInstance *ormdb.ORMDBInstance, dbName string) (*VersionedDB, error) {
	chainName := dbName
	dbName = dbName + "_"

	dbType := viper.GetString("ledger.state.stateDatabase")
	metadataDB, err := ormdb.CreateORMDatabase(ormDBInstance, dbName, dbType)
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*ormdb.ORMDatabase)
	return &VersionedDB{
		ormDBInstance: ormDBInstance,
		metadataDB:    metadataDB,
		chainName:     chainName,
		namespaceDBs:  namespaceDBMap,
		lsccStateCache: &statedb.LsccStateCache{
			Cache: make(map[string]*statedb.VersionedValue),
		},
	}, nil
}

func (v *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	//
	//logger.Debugf("GetState(). ns=%s, key=%s", namespace, key)
	//if namespace == "lscc" {
	//	if value := v.lsccStateCache.GetState(key); value != nil {
	//		return value, nil
	//	}
	//}
	//
	//db, err := v.getNamespaceDBHandle(namespace)
	//if err != nil {
	//	return nil, err
	//}
	//
	//keys := strings.Split(key, "#@#")
	//typeKey := keys[0]
	//entId := keys[1]
	//entType := db.ModelTypes[typeKey]
	//entPtr := reflect.New(entType)
	//db.DB.Find(entPtr, entId)
	//
	//if entPtr.IsNil() {
	//	return nil, nil
	//}
	//
	//kv, err := couchDocToKeyValue(entPtr.Elem())
	//if err != nil {
	//	return nil, err
	//}
	//
	//if namespace == "lscc" {
	//	vdb.lsccStateCache.SetState(key, kv.VersionedValue)
	//}
	//
	//return kv.VersionedValue, nil
	return nil, nil
}

func (v *VersionedDB) GetVersion(namespace string, key string) (*version.Height, error) {
	panic("implement me")
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
		var err error
		dbType := viper.GetString("ledger.state.stateDatabase")
		db, err = ormdb.CreateORMDatabase(v.ormDBInstance, namespaceDBName, dbType)
		if err != nil {
			return nil, err
		}
		v.namespaceDBs[namespace] = db
	}
	return db, nil
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
