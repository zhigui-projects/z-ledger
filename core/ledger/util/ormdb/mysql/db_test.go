package mysql

import (
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMySQLDB(t *testing.T) {
	conf := &config.ORMDBConfig{Username: "root", Password: "test", Host: "localhost", Port: 3306, MysqlConfig: &config.MysqlConfig{Charset: "utf8mb4", Collate: "utf8mb4_general_ci"}}
	db, err := CreateIfNotExistAndOpen(conf, "test")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = DropAndDelete(db)
	assert.NoError(t, err)
}
