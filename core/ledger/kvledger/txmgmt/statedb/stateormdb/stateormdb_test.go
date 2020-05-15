package stateormdb

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	ormdbconfig "github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

type TestSubModel struct {
	ID        string `gorm:"primary_key"`
	Model     string
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

type TestModel1 struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Age       sql.NullInt64 `ormdb:"datatype"`
	CreatedAt time.Time     `ormdb:"datatype"`
	UpdatedAt time.Time     `ormdb:"datatype"`
}

type TestModel2 struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Age       sql.NullInt64 `ormdb:"datatype"`
	CreatedAt time.Time     `ormdb:"datatype"`
	UpdatedAt time.Time     `ormdb:"datatype"`
}

type TestModel3 struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Age       sql.NullInt64 `ormdb:"datatype"`
	CreatedAt time.Time     `ormdb:"datatype"`
	UpdatedAt time.Time     `ormdb:"datatype"`
}

func TestNewVersionedDBProvider(t *testing.T) {
	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      redoLogPath: /tmp/ormdbredolog\n" +
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
	for _, db := range vdb.namespaceDBs {
		ormdb.DeleteORMDatabase(db)
	}
	provider.Close()
	os.RemoveAll(config.RedoLogPath)
}

func TestVersionedDB_ApplyUpdates(t *testing.T) {
	batchUpdate := statedb.NewUpdateBatch()
	key, testSubEfds, err := entitydefinition.RegisterEntity(&TestSubModel{}, 1)
	assert.NoError(t, err)
	testSubEfdsBytes, err := json.Marshal(testSubEfds)
	assert.NoError(t, err)

	key1, testEfds, err := entitydefinition.RegisterEntity(&TestModel{}, 2)
	assert.NoError(t, err)
	testEfdsBytes, err := json.Marshal(testEfds)
	assert.NoError(t, err)

	testsubmodelDel := TestSubModel{ID: "testsubmodelid3", Name: "testsubmodel3", Model: "testmodelid1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	testsubmodelDelBytes, err := json.Marshal(&testsubmodelDel)
	assert.NoError(t, err)
	batchUpdate.Put("mycc", "EntityFieldDefinition"+entitydefinition.ORMDB_SEPERATOR+key, testSubEfdsBytes, &version.Height{BlockNum: 1, TxNum: 1})
	batchUpdate.Put("mycc", "EntityFieldDefinition"+entitydefinition.ORMDB_SEPERATOR+key1, testEfdsBytes, &version.Height{BlockNum: 1, TxNum: 2})
	batchUpdate.Put("mycc", key+entitydefinition.ORMDB_SEPERATOR+"testsubmodelid3", testsubmodelDelBytes, &version.Height{BlockNum: 1, TxNum: 3})

	yaml := "---\n" +
		"ledger:\n" +
		"  state:\n" +
		"    ormDBConfig:\n" +
		"      username: test\n" +
		"      dbtype: sqlite3\n" +
		"      redoLogPath: /tmp/ormdbredolog\n" +
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

	err = vdb.ApplyUpdates(batchUpdate, &version.Height{BlockNum: 2, TxNum: 5})
	assert.NoError(t, err)

	savepoint := &SavePoint{}
	vdb.metadataDB.DB.Where(&SavePoint{Key: savePointKey}).Find(savepoint)
	savepointHeightBytes, err := base64.StdEncoding.DecodeString(savepoint.Height)
	assert.NoError(t, err)
	savepointHeight, _, err := version.NewHeightFromBytes(savepointHeightBytes)
	assert.NotNil(t, savepointHeight)
	assert.Equal(t, uint64(2), savepointHeight.BlockNum)
	assert.Equal(t, uint64(5), savepointHeight.TxNum)

	myccdb := vdb.namespaceDBs["mycc"]
	assert.True(t, myccdb.DB.HasTable(ormdb.ToTableName(key)))
	assert.True(t, myccdb.DB.HasTable(ormdb.ToTableName(key1)))

	var subModelEfds []entitydefinition.EntityFieldDefinition
	myccdb.DB.Where("owner=?", key).Find(&subModelEfds)
	assert.Equal(t, 6, len(subModelEfds))

	var modelEfds []entitydefinition.EntityFieldDefinition
	myccdb.DB.Where("owner=?", key1).Find(&modelEfds)
	assert.Equal(t, 14, len(modelEfds))

	var total []entitydefinition.EntityFieldDefinition
	myccdb.DB.Find(&total)
	assert.Equal(t, 20, len(total))

	vdb.metadataDB.DB.Close()
	for _, db := range vdb.namespaceDBs {
		db.DB.Close()
	}
	provider.Close()

	// Test restart provider
	batchUpdate1 := statedb.NewUpdateBatch()
	key2, testEfds1, err := entitydefinition.RegisterEntity(&TestModel1{}, 3)
	assert.NoError(t, err)
	testEfds1Bytes, err := json.Marshal(testEfds1)
	assert.NoError(t, err)

	key3, testEfds2, err := entitydefinition.RegisterEntity(&TestModel2{}, 4)
	assert.NoError(t, err)
	testEfds2Bytes, err := json.Marshal(testEfds2)
	assert.NoError(t, err)

	key4, testEfds3, err := entitydefinition.RegisterEntity(&TestModel3{}, 5)
	assert.NoError(t, err)
	testEfds3Bytes, err := json.Marshal(testEfds3)
	assert.NoError(t, err)

	testsubmodelRec := TestSubModel{ID: "testsubmodelid1", Name: "testsubmodel1", Model: "testmodelid1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	testsubmodelRecBytes, err := json.Marshal(&testsubmodelRec)
	assert.NoError(t, err)
	testsubmodelRec1 := TestSubModel{ID: "testsubmodelid2", Name: "testsubmodel2", Model: "testmodelid1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	testsubmodelRec1Bytes, err := json.Marshal(&testsubmodelRec1)
	assert.NoError(t, err)
	birthday := time.Now()
	MemberNumber := "testnum"
	testmodelRec := &TestModel{ID: "testmodelid1", Name: "testmodel1", Age: sql.NullInt64{Int64: 11}, Birthday: &birthday, Email: "test@test.com", Role: "admin", MemberNumber: &MemberNumber, Num: 12, Address: "address", IgnoreMe: 22, TestSubModels: []TestSubModel{testsubmodelRec, testsubmodelRec1}, Attachment: []byte("test"), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	testmodelRecBytes, err := json.Marshal(testmodelRec)
	assert.NoError(t, err)

	batchUpdate1.Put("mycc", "EntityFieldDefinition"+entitydefinition.ORMDB_SEPERATOR+key2, testEfds1Bytes, &version.Height{BlockNum: 2, TxNum: 1})
	batchUpdate1.Put("mycc", "EntityFieldDefinition"+entitydefinition.ORMDB_SEPERATOR+key3, testEfds2Bytes, &version.Height{BlockNum: 2, TxNum: 2})
	batchUpdate1.Put("mycc", "EntityFieldDefinition"+entitydefinition.ORMDB_SEPERATOR+key4, testEfds3Bytes, &version.Height{BlockNum: 2, TxNum: 3})
	batchUpdate1.Put("mycc", key+entitydefinition.ORMDB_SEPERATOR+"testsubmodelid1", testsubmodelRecBytes, &version.Height{BlockNum: 2, TxNum: 4})
	batchUpdate1.Put("mycc", key+entitydefinition.ORMDB_SEPERATOR+"testsubmodelid2", testsubmodelRec1Bytes, &version.Height{BlockNum: 2, TxNum: 5})
	batchUpdate1.Delete("mycc", key+entitydefinition.ORMDB_SEPERATOR+"testsubmodelid3", &version.Height{BlockNum: 2, TxNum: 6})
	batchUpdate1.Put("mycc", key1+entitydefinition.ORMDB_SEPERATOR+"testmodelid1", testmodelRecBytes, &version.Height{BlockNum: 2, TxNum: 7})

	provider1, err := NewVersionedDBProvider(config, nil, statedb.NewCache(config.UserCacheSizeMBs, []string{"lscc"}))
	assert.NoError(t, err)
	assert.Equal(t, 64, provider1.ormDBInstance.Config.UserCacheSizeMBs)

	_, err = provider1.GetDBHandle("mychannel")
	assert.NoError(t, err)
	vdb1 := provider1.databases["mychannel"]

	err = vdb1.ApplyUpdates(batchUpdate1, &version.Height{BlockNum: 3, TxNum: 5})
	assert.NoError(t, err)

	savepoint1 := &SavePoint{}
	vdb1.metadataDB.DB.Where(&SavePoint{Key: savePointKey}).Find(savepoint1)
	savepoint1HeightBytes, err := base64.StdEncoding.DecodeString(savepoint1.Height)
	assert.NoError(t, err)
	savepoint1Height, _, err := version.NewHeightFromBytes(savepoint1HeightBytes)
	assert.NotNil(t, savepoint1Height)
	assert.Equal(t, uint64(3), savepoint1Height.BlockNum)
	assert.Equal(t, uint64(5), savepoint1Height.TxNum)

	myccdb1 := vdb1.namespaceDBs["mycc"]
	assert.True(t, myccdb1.DB.HasTable(ormdb.ToTableName(key)))
	assert.True(t, myccdb1.DB.HasTable(ormdb.ToTableName(key1)))
	assert.True(t, myccdb1.DB.HasTable(ormdb.ToTableName(key2)))

	var subModelEfds1 []entitydefinition.EntityFieldDefinition
	myccdb1.DB.Where("owner=?", key).Find(&subModelEfds1)
	assert.Equal(t, 6, len(subModelEfds1))

	var modelEfds1 []entitydefinition.EntityFieldDefinition
	myccdb1.DB.Where("owner=?", key1).Find(&modelEfds1)
	assert.Equal(t, 14, len(modelEfds1))

	var model1Efds []entitydefinition.EntityFieldDefinition
	myccdb1.DB.Where("owner=?", key2).Find(&model1Efds)
	assert.Equal(t, 5, len(model1Efds))

	var total1 []entitydefinition.EntityFieldDefinition
	myccdb1.DB.Find(&total1)
	assert.Equal(t, 35, len(total1))

	submodels := reflect.New(reflect.SliceOf(myccdb1.ModelTypes[key].StructType())).Interface()
	myccdb1.DB.Table(ormdb.ToTableName(key)).Find(submodels)
	submodelsbytes, err := json.Marshal(&submodels)
	assert.NoError(t, err)

	submodels1 := make([]TestSubModel, 0)
	err = json.Unmarshal(submodelsbytes, &submodels1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(submodels1))

	models := reflect.New(reflect.SliceOf(myccdb1.ModelTypes[key1].StructType())).Interface()
	myccdb1.DB.Table(ormdb.ToTableName(key1)).Find(models)
	modelsbytes, err := json.Marshal(&models)
	assert.NoError(t, err)

	models1 := make([]TestModel, 0)
	err = json.Unmarshal(modelsbytes, &models1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(models1))

	ormdb.DeleteORMDatabase(vdb1.metadataDB)
	for _, db := range vdb1.namespaceDBs {
		ormdb.DeleteORMDatabase(db)
	}
	provider1.Close()
	os.RemoveAll(config.RedoLogPath)
}
