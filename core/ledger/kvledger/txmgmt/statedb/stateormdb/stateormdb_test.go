package stateormdb

import (
	"bytes"
	"database/sql"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type TestSubModel struct {
	ID        string     `gorm:"primary_key"`
	CreatedAt time.Time  `ormdb:"datatype"`
	UpdatedAt time.Time  `ormdb:"datatype"`
	DeletedAt *time.Time `sql:"index" ormdb:"datatype"`
	Name      string
}

type TestModel struct {
	ID            string `gorm:"primary_key"`
	Name          string
	Age           sql.NullInt64  `ormdb:"datatype"`
	Birthday      *time.Time     `ormdb:"datatype"`
	Email         string         `gorm:"type:varchar(100);unique_index"`
	Role          string         `gorm:"size:255"`        // set field size to 255
	MemberNumber  *string        `gorm:"unique;not null"` // set member number to unique and not null
	Num           int            `gorm:"AUTO_INCREMENT"`  // set num to auto incrementable
	Address       string         `gorm:"index:addr"`      // create index with name `addr` for address
	IgnoreMe      int            `gorm:"-"`               // ignore this field
	TestSubModels []TestSubModel `ormdb:"entity"`
	Attachment    []byte
	CreatedAt     time.Time `ormdb:"datatype"`
	UpdatedAt     time.Time `ormdb:"datatype"`
}

func TestNewVersionedDBProvider(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      redoLogPath: /tmp/ormdbredolog\n" +
		"      maxBatchUpdateSize: 50\n" +
		"      userCacheSizeMBs: 64\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	config := &ormdbconfig.ORMDBConfig{Sqlite3Config: &ormdbconfig.Sqlite3Config{}}
	_ = mapstructure.Decode(viper.Get("ledger.state.ormDBConfig"), config)

	provider, err := NewVersionedDBProvider(config, nil, statedb.NewCache(config.UserCacheSizeMBs, []string{"lscc"}))
	assert.NoError(t, err)
	assert.Equal(t, 64, provider.ormDBInstance.Config.UserCacheSizeMBs)
	provider.Close()
	os.RemoveAll(config.RedoLogPath)
}

func TestVersionedDBProvider_GetDBHandle(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      redoLogPath: /tmp/ormdbredolog\n" +
		"      maxBatchUpdateSize: 50\n" +
		"      userCacheSizeMBs: 64\n" +
		"      sqlite3Config:\n" +
		"        path: /tmp/ormdb\n"

	defer viper.Reset()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader([]byte(yaml))); err != nil {
		t.Fatalf("Error reading config: %s", err)
	}

	config := &ormdbconfig.ORMDBConfig{Sqlite3Config: &ormdbconfig.Sqlite3Config{}}
	_ = mapstructure.Decode(viper.Get("ledger.state.ormDBConfig"), config)

	provider, err := NewVersionedDBProvider(config, nil, statedb.NewCache(config.UserCacheSizeMBs, []string{"lscc"}))
	assert.NoError(t, err)
	assert.Equal(t, 64, provider.ormDBInstance.Config.UserCacheSizeMBs)

	_, err = provider.GetDBHandle("mychannel")
	assert.NoError(t, err)
	vdb := provider.databases["mychannel"]
	ormdb.DeleteORMDatabase(vdb.metadataDB)
	provider.Close()
	os.RemoveAll(config.RedoLogPath)
}

func TestVersionedDB_ApplyUpdates(t *testing.T) {
	//nsUpdates := make(map[string]*statedb.VersionedValue)
	//statedb.VersionedValue{}
	//updates := &statedb.UpdateBatch{}
}
