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
	"reflect"
	"sort"
	"strings"
	"sync"
)

type batchableEntity struct {
	Entity  interface{}
	Deleted bool
}

type batchableEntityFieldDefinition struct {
	EntityName             string
	Seq                    int
	EntityFieldDefinitions []entitydefinition.EntityFieldDefinition
}

type committer struct {
	db             *ormdb.ORMDatabase
	batchUpdateMap []*batchableEntity
	efdMap         []*batchableEntityFieldDefinition
	namespace      string
	cacheKVs       statedb.CacheKVs
	cacheEnabled   bool
}

func (v *VersionedDB) buildCommitters(updates *statedb.UpdateBatch) ([]*committer, error) {
	namespaces := updates.GetUpdatedNamespaces()

	var wg sync.WaitGroup
	nsCommittersChan := make(chan *committer, len(namespaces))
	defer close(nsCommittersChan)
	errsChan := make(chan error, len(namespaces))
	defer close(errsChan)

	for _, ns := range namespaces {
		nsUpdates := updates.GetUpdates(ns)
		wg.Add(1)
		go func(ns string) {
			defer wg.Done()
			committer, err := v.buildCommittersForNs(ns, nsUpdates)
			if err != nil {
				errsChan <- err
				return
			}
			nsCommittersChan <- committer
		}(ns)
	}
	wg.Wait()

	var allCommitters []*committer
	select {
	case err := <-errsChan:
		return nil, err
	default:
		for i := 0; i < len(namespaces); i++ {
			allCommitters = append(allCommitters, <-nsCommittersChan)
		}
	}

	return allCommitters, nil
}

func (v *VersionedDB) buildCommittersForNs(ns string, nsUpdates map[string]*statedb.VersionedValue) (*committer, error) {
	db, err := v.getNamespaceDBHandle(ns)
	if err != nil {
		return nil, err
	}
	cacheEnabled := v.cache.Enabled(ns)

	committer := &committer{
		db:             db,
		batchUpdateMap: make([]*batchableEntity, 0),
		efdMap:         make([]*batchableEntityFieldDefinition, 0),
		namespace:      ns,
		cacheKVs:       make(statedb.CacheKVs),
		cacheEnabled:   cacheEnabled,
	}

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
				err = json.Unmarshal(vv.Value, &efds)
				if err != nil {
					return nil, err
				}

				var seq int
				var aefds []entitydefinition.EntityFieldDefinition
				for i, efd := range efds {
					efdr := &efd
					if i == 0 {
						seq = efd.Seq
					}
					efdr.ID = strings.ReplaceAll(util.GenerateUUID(), "-", "")
					efdr.VerAndMeta = verAndMeta
					aefds = append(aefds, *efdr)
				}

				committer.efdMap = append(committer.efdMap, &batchableEntityFieldDefinition{EntityName: entityKey, EntityFieldDefinitions: aefds, Seq: seq})
			}
		} else {
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

			if vv.Value != nil {
				err = json.Unmarshal(vv.Value, entity)
				if err != nil {
					return nil, err
				}
			} else {
				if id.Type.Kind() == reflect.String {
					reflect.ValueOf(entity).Elem().FieldByName("ID").SetString(entityKey)
				} else {
					return nil, errors.New("not supported entity ID field type")
				}
			}

			reflect.ValueOf(entity).Elem().FieldByName("VerAndMeta").SetString(verAndMeta)
			committer.batchUpdateMap = append(committer.batchUpdateMap, &batchableEntity{Entity: entity, Deleted: vv.Value == nil})
		}

		committer.addToCacheUpdate(key, vv)
	}
	return committer, nil
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
	sort.SliceStable(c.efdMap, func(i, j int) bool {
		return c.efdMap[i].Seq < c.efdMap[j].Seq
	})

	for _, update := range c.efdMap {
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
	}

	for _, update := range c.batchUpdateMap {
		if update.Deleted {
			if err = c.db.DB.Delete(update.Entity).Error; err != nil {
				return err
			}
		} else {
			if err = c.db.DB.Save(update.Entity).Error; err != nil {
				return err
			}
		}
	}
	return err
}
