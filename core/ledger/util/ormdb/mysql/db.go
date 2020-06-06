/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package mysql

import (
	"fmt"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

// CreateIfNotExistAndOpen opens a gorm mysql database instance
func CreateIfNotExistAndOpen(config *config.ORMDBConfig, dbName string) (db *gorm.DB, err error) {
	con := fmt.Sprintf("%s:%s@(%s:%d)/", config.Username, config.Password, config.Host, config.Port)
	fmt.Println(con)
	cdb, err := gorm.Open("mysql", con)
	if err != nil {
		panic(errors.WithMessage(err, "error connect mysql"))
	}

	createDBStatement := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET %s COLLATE %s", dbName, config.MysqlConfig.Charset, config.MysqlConfig.Collate)
	err = cdb.Exec(createDBStatement).Error
	if err != nil {
		panic(errors.WithMessage(err, "error creating mysql db"))
	}
	cdb.Close()
	dbcon := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True", config.Username, config.Password, config.Host, config.Port, dbName, config.MysqlConfig.Charset)

	db, err = gorm.Open("mysql", dbcon)
	if err != nil {
		panic(errors.WithMessage(err, "error connect mysql database"))
	}
	return
}

// DropAndDelete deletes a gorm mysql database instance
func DropAndDelete(db *gorm.DB) error {
	err := db.Close()
	if err != nil {
		return errors.WithMessage(err, "error closing mysql db")
	}
	return nil
}
