/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"fmt"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/mitchellh/mapstructure"
	"path/filepath"

	coreconfig "github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	"github.com/spf13/viper"
)

func ledgerConfig() *ledger.Config {
	// set defaults
	warmAfterNBlocks := 1
	if viper.IsSet("ledger.state.couchDBConfig.warmIndexesAfterNBlocks") {
		warmAfterNBlocks = viper.GetInt("ledger.state.couchDBConfig.warmIndexesAfterNBlocks")
	}
	internalQueryLimit := 1000
	if viper.IsSet("ledger.state.couchDBConfig.internalQueryLimit") {
		internalQueryLimit = viper.GetInt("ledger.state.couchDBConfig.internalQueryLimit")
	}
	maxBatchUpdateSize := 500
	if viper.IsSet("ledger.state.couchDBConfig.maxBatchUpdateSize") {
		maxBatchUpdateSize = viper.GetInt("ledger.state.couchDBConfig.maxBatchUpdateSize")
	}
	collElgProcMaxDbBatchSize := 5000
	if viper.IsSet("ledger.pvtdataStore.collElgProcMaxDbBatchSize") {
		collElgProcMaxDbBatchSize = viper.GetInt("ledger.pvtdataStore.collElgProcMaxDbBatchSize")
	}
	collElgProcDbBatchesInterval := 1000
	if viper.IsSet("ledger.pvtdataStore.collElgProcDbBatchesInterval") {
		collElgProcDbBatchesInterval = viper.GetInt("ledger.pvtdataStore.collElgProcDbBatchesInterval")
	}
	purgeInterval := 100
	if viper.IsSet("ledger.pvtdataStore.purgeInterval") {
		purgeInterval = viper.GetInt("ledger.pvtdataStore.purgeInterval")
	}

	rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
	conf := &ledger.Config{
		RootFSPath: rootFSPath,
		StateDBConfig: &ledger.StateDBConfig{
			StateDatabase: viper.GetString("ledger.state.stateDatabase"),
			CouchDB:       &couchdb.Config{},
		},
		PrivateDataConfig: &ledger.PrivateDataConfig{
			MaxBatchSize:    collElgProcMaxDbBatchSize,
			BatchesInterval: collElgProcDbBatchesInterval,
			PurgeInterval:   purgeInterval,
		},
		HistoryDBConfig: &ledger.HistoryDBConfig{
			Enabled: viper.GetBool("ledger.history.enableHistoryDatabase"),
		},
	}

	if conf.StateDBConfig.StateDatabase == "CouchDB" {
		conf.StateDBConfig.CouchDB = &couchdb.Config{
			Address:                 viper.GetString("ledger.state.couchDBConfig.couchDBAddress"),
			Username:                viper.GetString("ledger.state.couchDBConfig.username"),
			Password:                viper.GetString("ledger.state.couchDBConfig.password"),
			MaxRetries:              viper.GetInt("ledger.state.couchDBConfig.maxRetries"),
			MaxRetriesOnStartup:     viper.GetInt("ledger.state.couchDBConfig.maxRetriesOnStartup"),
			RequestTimeout:          viper.GetDuration("ledger.state.couchDBConfig.requestTimeout"),
			InternalQueryLimit:      internalQueryLimit,
			MaxBatchUpdateSize:      maxBatchUpdateSize,
			WarmIndexesAfterNBlocks: warmAfterNBlocks,
			CreateGlobalChangesDB:   viper.GetBool("ledger.state.couchDBConfig.createGlobalChangesDB"),
			RedoLogPath:             filepath.Join(rootFSPath, "couchdbRedoLogs"),
			UserCacheSizeMBs:        viper.GetInt("ledger.state.couchDBConfig.cacheSize"),
		}
	}

	if conf.StateDBConfig.StateDatabase == "ORMDB" {
		config := &ormdbconfig.ORMDBConfig{Sqlite3Config: &ormdbconfig.Sqlite3Config{}}
		_ = mapstructure.Decode(viper.Get("ledger.state.ormDBConfig"), config)
		fmt.Println("232313123123" + config.DbType + "232313123123")
		fmt.Println("232313123123" + config.Username + "232313123123")
		config.RedoLogPath = filepath.Join(rootFSPath, "ormdbRedoLogs")
		config.DbType = viper.GetString("ledger.state.ormDBConfig.dbtype")
		fmt.Println("232313123123" + config.DbType + "232313123123")
		config.UserCacheSizeMBs = viper.GetInt("ledger.state.ormDBConfig.userCacheSizeMBs")
		config.Username = viper.GetString("ledger.state.ormDBConfig.username")
		fmt.Println("232313123123" + config.Username + "232313123123")
		if config.DbType == "sqlite3" {
			config.Sqlite3Config.Path = filepath.Join(rootFSPath, "ormdb")
		}
		conf.StateDBConfig.ORMDB = config
	}
	return conf
}
