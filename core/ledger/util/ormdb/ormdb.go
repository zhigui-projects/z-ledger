/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ormdb

import (
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/sqllite3"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"sync"
)

var logger = flogging.MustGetLogger("ormdb")

//ORMDBInstance represents a ORMDB instance
type ORMDBInstance struct {
	Config *config.ORMDBConfig
}

//ORMDatabase represents a database within a ORMDB instance
type ORMDatabase struct {
	ORMDBInstance *ORMDBInstance
	DBName        string
	DB            *gorm.DB
	Type          string
	RWMutex       sync.RWMutex
	ModelTypes    map[string]entitydefinition.DynamicStruct
}

// NewORMDBInstance create a ORMDB instance through ORMDBConfig
func NewORMDBInstance(config *config.ORMDBConfig, metricsProvider metrics.Provider) (*ORMDBInstance, error) {
	ormDBInstance := &ORMDBInstance{Config: config}
	return ormDBInstance, nil
}

// CreateORMDatabase creates a ORM database object, as well as the underlying database if it does not exist
func CreateORMDatabase(ormDBInstance *ORMDBInstance, dbName string) (*ORMDatabase, error) {
	var db *gorm.DB
	var err error
	switch ormDBInstance.Config.DBType {
	case "sqlite3":
		db, err = sqllite3.CreateIfNotExistAndOpen(ormDBInstance.Config, dbName)
		if err != nil {
			logger.Errorf("create sqlite3 database with dbname [%s] failed", dbName)
			return nil, errors.WithMessage(err, "create sqlite3 database failed")
		}
	default:
		return nil, errors.New("not supported database type")
	}
	ormDBDatabase := &ORMDatabase{
		ORMDBInstance: ormDBInstance,
		DBName:        dbName,
		DB:            db,
		Type:          ormDBInstance.Config.DBType,
		ModelTypes:    make(map[string]entitydefinition.DynamicStruct),
	}

	return ormDBDatabase, nil
}

// DeleteORMDatabase deletes a ORM database object, as well as the underlying database if it exists
func DeleteORMDatabase(ormDatabase *ORMDatabase) error {
	var err error
	switch ormDatabase.Type {
	case "sqlite3":
		err = sqllite3.DropAndDelete(ormDatabase.DB, ormDatabase.ORMDBInstance.Config, ormDatabase.DBName)
		if err != nil {
			logger.Errorf("delete sqlite3 database with dbname [%s] failed", ormDatabase.DBName)
			return errors.WithMessage(err, "delete sqlite3 database failed")
		}
	default:
		return errors.New("not supported database type")
	}
	return nil
}
