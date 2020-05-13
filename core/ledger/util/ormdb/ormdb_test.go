/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ormdb

import (
	"bytes"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/mitchellh/mapstructure"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type ABC struct {
	ID   string `gorm:"primary_key"`
	Test string
}

func TestNewORMDBInstance(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	config := &ormdbconfig.ORMDBConfig{Sqlite3Config: &ormdbconfig.Sqlite3Config{}}
	_ = mapstructure.Decode(viper.Get("ledger.state.ormDBConfig"), config)
	ormDBInstance, err := NewORMDBInstance(config, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ormDBInstance)
	assert.Equal(t, "test", ormDBInstance.Config.Username)
	assert.Equal(t, "/tmp/ormdb", ormDBInstance.Config.Sqlite3Config.Path)
}

func TestCreateORMDatabase(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	config := &ormdbconfig.ORMDBConfig{Sqlite3Config: &ormdbconfig.Sqlite3Config{}}
	_ = mapstructure.Decode(viper.Get("ledger.state.ormDBConfig"), config)
	ormDBInstance, err := NewORMDBInstance(config, nil)
	assert.NoError(t, err)

	db, err := CreateORMDatabase(ormDBInstance, "test")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.DB.CreateTable(&ABC{})
	testABC := &ABC{ID: "6F1351B1-F6D1-4B05-AE8F-99CBE6ED106B", Test: "123"}
	db.DB.Save(testABC)
	testABC1 := &ABC{ID: "6F1351B1-F6D1-4B05-AE8F-99CBE6ED106B"}
	db.DB.Find(testABC1)
	assert.Equal(t, "123", testABC1.Test)
	testABC1.Test = "789"
	db.DB.Save(testABC1)
	testABC2 := &ABC{}
	db.DB.Where(&ABC{ID: "6F1351B1-F6D1-4B05-AE8F-99CBE6ED106B"}).Find(testABC2)
	assert.Equal(t, "789", testABC2.Test)
	err = DeleteORMDatabase(db)
	assert.NoError(t, err)

}
