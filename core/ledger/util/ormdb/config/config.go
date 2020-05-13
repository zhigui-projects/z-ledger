/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package config

// ORMDBConfig
type ORMDBConfig struct {
	DBType             string
	Username           string
	Password           string
	Host               string
	Port               int
	RedoLogPath        string
	MaxBatchUpdateSize int
	UserCacheSizeMBs   int
	Sqlite3Config      *Sqlite3Config
}

// Sqlite3Config
type Sqlite3Config struct {
	Path string `mapstructure:"path" json:"path" yaml:"path"`
}
