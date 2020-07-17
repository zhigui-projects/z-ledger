/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"path/filepath"

	coreconfig "github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/archive"
	config "github.com/hyperledger/fabric/core/ledger/dfs/common"
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

	enableArchive := false
	if viper.IsSet("ledger.archive.enable") {
		enableArchive = viper.GetBool("ledger.archive.enable")
	}

	rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
	conf := &ledger.Config{
		RootFSPath: rootFSPath,
		ArchiveConfig: &archive.Config{
			Enabled: enableArchive,
			Trigger: &archive.Trigger{
				CheckInterval: viper.GetDuration("ledger.archive.trigger.checkInterval"),
				BeginTime:     viper.GetTime("ledger.archive.trigger.beginTime"),
				EndTime:       viper.GetTime("ledger.archive.trigger.endTime"),
			},
			Type: viper.GetString("ledger.archive.type"),
			HdfsConf: &config.HdfsConfig{
				User:                viper.GetString("ledger.archive.hdfsConfig.user"),
				NameNodes:           viper.GetStringSlice("ledger.archive.hdfsConfig.nameNodes"),
				UseDatanodeHostname: viper.GetBool("ledger.archive.hdfsConfig.useDatanodeHostname"),
			},
			IpfsConf: &config.IpfsConfig{
				Url: viper.GetString("ledger.archive.ipfsConfig.url"),
			},
		},
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
		config := &ormdbconfig.ORMDBConfig{}
		config.RedoLogPath = filepath.Join(rootFSPath, "ormdbRedoLogs")
		config.DbType = viper.GetString("ledger.state.ormDBConfig.dbtype")
		config.Host = viper.GetString("ledger.state.ormDBConfig.host")
		config.Port = viper.GetInt("ledger.state.ormDBConfig.port")
		config.UserCacheSizeMBs = viper.GetInt("ledger.state.ormDBConfig.userCacheSizeMBs")
		config.Username = viper.GetString("ledger.state.ormDBConfig.username")
		config.Password = viper.GetString("ledger.state.ormDBConfig.password")
		if config.DbType == "sqlite3" {
			sqlite3Config := &ormdbconfig.Sqlite3Config{}
			sqlite3Config.Path = filepath.Join(rootFSPath, "ormdb")
			config.Sqlite3Config = sqlite3Config
		}
		if config.DbType == "mysql" {
			mysqlConfig := &ormdbconfig.MysqlConfig{}
			mysqlConfig.Charset = viper.GetString("ledger.state.ormDBConfig.mysqlConfig.charset")
			mysqlConfig.Collate = viper.GetString("ledger.state.ormDBConfig.mysqlConfig.collate")
			config.MysqlConfig = mysqlConfig
		}
		conf.StateDBConfig.ORMDB = config
	}
	return conf
}
