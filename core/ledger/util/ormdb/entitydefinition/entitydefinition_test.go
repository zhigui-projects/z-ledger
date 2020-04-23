package entitydefinition

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

type TestSubModel struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `ormdb:"datatype"`
	UpdatedAt time.Time  `ormdb:"datatype"`
	DeletedAt *time.Time `sql:"index" ormdb:"datatype"`
	Name      string
}

type TestModel struct {
	Name           string
	Age            sql.NullInt64   `ormdb:"datatype"`
	Birthday       *time.Time      `ormdb:"datatype"`
	Email          string          `gorm:"type:varchar(100);unique_index"`
	Role           string          `gorm:"size:255"`        // set field size to 255
	MemberNumber   *string         `gorm:"unique;not null"` // set member number to unique and not null
	Num            int             `gorm:"AUTO_INCREMENT"`  // set num to auto incrementable
	Address        string          `gorm:"index:addr"`      // create index with name `addr` for address
	IgnoreMe       int             `gorm:"-"`               // ignore this field
	TestSubModels  []TestSubModel  `ormdb:"entity"`
	TestSubModels1 []*TestSubModel `ormdb:"entity"`
	IgnoreMe1      *int            `gorm:"-"` // ignore this field

}

func TestRegisterEntity(t *testing.T) {
	gormModelType := reflect.TypeOf(&gorm.Model{})
	fmt.Println(gormModelType.Elem().PkgPath())
	subKey, subEntityFieldDefinitions, err := RegisterEntity(&TestSubModel{})
	key, entityFieldDefinitions, err := RegisterEntity(&TestModel{})
	subEntityFieldDefinitionsBytes, err := json.Marshal(subEntityFieldDefinitions)
	entityFieldDefinitionsBytes, err := json.Marshal(entityFieldDefinitions)
	var subEntityFieldDefinitions1 []*EntityFieldDefinition
	err = json.Unmarshal(subEntityFieldDefinitionsBytes, &subEntityFieldDefinitions1)
	var entityFieldDefinitions1 []*EntityFieldDefinition
	err = json.Unmarshal(entityFieldDefinitionsBytes, &entityFieldDefinitions1)
	fmt.Println(subKey)
	fmt.Println(subEntityFieldDefinitions)
	fmt.Println(key)
	fmt.Println(entityFieldDefinitions)
	fmt.Println(err)

	registry := make(map[string]DynamicStruct)
	subEntityDS := NewBuilder().AddEntityFieldDefinition(subEntityFieldDefinitions1, registry).Build()
	registry[subKey] = subEntityDS
	entityDS := NewBuilder().AddEntityFieldDefinition(entityFieldDefinitions1, registry).Build()
	registry[key] = entityDS
	fmt.Println(subEntityDS)
	fmt.Println(entityDS)
}
