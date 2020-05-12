/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"sync"
)

type batchableEntity struct {
	Entity                 interface{}
	EntityName             string
	EntityFieldDefinitions []entitydefinition.EntityFieldDefinition
	Deleted                bool
}

type committer struct {
	db             *ormdb.ORMDatabase
	batchUpdateMap map[string]*batchableEntity
	namespace      string
	cacheKVs       statedb.CacheKVs
	cacheEnabled   bool
}

func (v *VersionedDB) buildCommitters(updates *statedb.UpdateBatch) ([]*committer, error) {
	namespaces := updates.GetUpdatedNamespaces()

	var wg sync.WaitGroup
	nsCommittersChan := make(chan []*committer, len(namespaces))
	defer close(nsCommittersChan)
	errsChan := make(chan error, len(namespaces))
	defer close(errsChan)

	for _, ns := range namespaces {
		nsUpdates := updates.GetUpdates(ns)
		wg.Add(1)
		go func(ns string) {
			defer wg.Done()
			committers, err := v.buildCommittersForNs(ns, nsUpdates)
			if err != nil {
				errsChan <- err
				return
			}
			nsCommittersChan <- committers
		}(ns)
	}
	wg.Wait()

	var allCommitters []*committer
	select {
	case err := <-errsChan:
		return nil, err
	default:
		for i := 0; i < len(namespaces); i++ {
			allCommitters = append(allCommitters, <-nsCommittersChan...)
		}
	}

	return allCommitters, nil
}

func (v *VersionedDB) buildCommittersForNs(ns string, nsUpdates map[string]*statedb.VersionedValue) ([]*committer, error) {
	db, err := v.getNamespaceDBHandle(ns)
	if err != nil {
		return nil, err
	}
	// for each namespace, build mutiple committers based on the maxBatchSize
	maxBatchSize := db.ORMDBInstance.Config.MaxBatchUpdateSize
	numCommitters := 1
	if maxBatchSize > 0 {
		numCommitters = int(math.Ceil(float64(len(nsUpdates)) / float64(maxBatchSize)))
	}
	committers := make([]*committer, numCommitters)
	cacheEnabled := v.cache.Enabled(ns)

	for i := 0; i < numCommitters; i++ {
		committers[i] = &committer{
			db:             db,
			batchUpdateMap: make(map[string]*batchableEntity),
			namespace:      ns,
			cacheKVs:       make(statedb.CacheKVs),
			cacheEnabled:   cacheEnabled,
		}
	}

	i := 0
	for key, vv := range nsUpdates {
		keys, err := parseKey(key)
		if err != nil {
			return nil, err
		}

		entityName := keys[0]
		entityKey := keys[1]
		verAndMeta, err := encodeVersionAndMetadata(vv.Version, vv.Metadata)
		if err != nil {
			return nil, err
		}
		if entityName == "EntityFieldDefinition" {
			if vv.Value == nil {
				return nil, errors.New("create table value cannot nil")
			} else {
				efds := make([]entitydefinition.EntityFieldDefinition, 0)
				err = json.Unmarshal(vv.Value, efds)
				if err != nil {
					return nil, err
				}

				for _, efd := range efds {
					efd.ID = util.GenerateUUID()
					efd.VerAndMeta = verAndMeta
				}

				committers[i].batchUpdateMap[key] = &batchableEntity{EntityName: entityKey, EntityFieldDefinitions: efds, Deleted: false}
			}
		} else {
			entity := db.ModelTypes[entityName].Interface()
			if vv.Value != nil {
				err = json.Unmarshal(vv.Value, entity)
				if err != nil {
					return nil, err
				}
			} else {
				id, exist := db.ModelTypes[entityName].FieldByName("ID")
				if !exist {
					return nil, errors.New("entity no ID field")
				}
				if id.Type.Kind() == reflect.String {
					reflect.ValueOf(entity).Elem().FieldByName("ID").SetString(entityKey)
				} else {
					return nil, errors.New("not supported entity ID field type")
				}
			}

			reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").SetString(verAndMeta)
			committers[i].batchUpdateMap[key] = &batchableEntity{Entity: entity, Deleted: vv.Value == nil}
		}

		committers[i].addToCacheUpdate(key, vv)
	}
	return committers, nil
}

func (c *committer) addToCacheUpdate(key string, vv *statedb.VersionedValue) {
	if !c.cacheEnabled {
		return
	}

	if vv.Value == nil {
		// nil value denotes a delete operation
		c.cacheKVs[key] = nil
		return
	}

	c.cacheKVs[key] = &statedb.CacheValue{
		VersionBytes: vv.Version.ToBytes(),
		Value:        vv.Value,
		Metadata:     vv.Metadata,
	}
}

func (v *VersionedDB) executeCommitter(committers []*committer) error {
	errsChan := make(chan error, len(committers))
	defer close(errsChan)
	var wg sync.WaitGroup
	wg.Add(len(committers))

	for _, c := range committers {
		go func(c *committer) {
			defer wg.Done()
			if err := c.commitUpdates(); err != nil {
				errsChan <- err
			}
		}(c)
	}
	wg.Wait()

	select {
	case err := <-errsChan:
		return err
	default:
		return nil
	}
}

// commitUpdates commits the given updates to ormdb
func (c *committer) commitUpdates() error {
	var err error
	for _, update := range c.batchUpdateMap {
		if update.Deleted {
			if err = c.db.DB.Delete(update.Entity).Error; err != nil {
				return err
			}
		} else {
			if update.Entity == nil {
				for _, efd := range update.EntityFieldDefinitions {
					if err = c.db.DB.Save(efd).Error; err != nil {
						return err
					}
				}
				c.db.RWMutex.Lock()
				ds := entitydefinition.NewBuilder().AddEntityFieldDefinition(update.EntityFieldDefinitions, c.db.ModelTypes).Build()
				c.db.ModelTypes[update.EntityName] = ds
				c.db.RWMutex.Unlock()
				entity := ds.Interface()
				if !c.db.DB.HasTable(entity) {
					if err = c.db.DB.CreateTable(entity).Error; err != nil {
						return err
					}
				}
			} else {
				if err = c.db.DB.Save(update.Entity).Error; err != nil {
					return err
				}
			}

		}
	}
	return err
}
