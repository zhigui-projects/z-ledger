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
	err = DeleteORMDatabase(db, "test")
	assert.NoError(t, err)
}
