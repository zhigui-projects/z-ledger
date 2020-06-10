/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

// ORMDBConfig
type ORMDBConfig struct {
	DbType           string         `mapstructure:"dbtype" json:"dbtype" yaml:"dbtype"`
	Username         string         `mapstructure:"username" json:"username" yaml:"username"`
	Password         string         `mapstructure:"password" json:"password" yaml:"password"`
	Host             string         `mapstructure:"host" json:"host" yaml:"host"`
	Port             int            `mapstructure:"port" json:"port" yaml:"port"`
	RedoLogPath      string         `mapstructure:"redoLogPath" json:"redoLogPath" yaml:"redoLogPath"`
	UserCacheSizeMBs int            `mapstructure:"userCacheSizeMBs" json:"userCacheSizeMBs" yaml:"userCacheSizeMBs"`
	Sqlite3Config    *Sqlite3Config `mapstructure:"sqlite3Config" json:"sqlite3Config" yaml:"sqlite3Config"`
	MysqlConfig      *MysqlConfig   `mapstructure:"mysqlConfig" json:"mysqlConfig" yaml:"mysqlConfig"`
}

// Sqlite3Config
type Sqlite3Config struct {
	Path string `mapstructure:"path" json:"path" yaml:"path"`
}

type MysqlConfig struct {
	Charset string `mapstructure:"charset" json:"charset" yaml:"charset"`
	Collate string `mapstructure:"collate" json:"collate" yaml:"collate"`
}
