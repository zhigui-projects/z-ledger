package entitydefinition

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"
)

const ORMDB_SEPERATOR = "$#$"

const (
	SampleString  = ""
	SampleBool    = true
	SampleInt     = int(1)
	SampleInt8    = int8(1)
	SampleInt16   = int16(1)
	SampleInt32   = int32(1)
	SampleInt64   = int64(1)
	SampleUint    = uint(1)
	SampleUint8   = uint8(1)
	SampleUint16  = uint16(1)
	SampleUint32  = uint32(1)
	SampleUint64  = uint64(1)
	SampleFloat32 = float32(1.1)
	SampleFloat64 = float64(1.1)
)

var datatype map[string]reflect.Type

func init() {
	datatype = make(map[string]reflect.Type)
	timeTimeType := reflect.TypeOf(time.Time{})
	datatype[timeTimeType.PkgPath()+ORMDB_SEPERATOR+timeTimeType.Name()] = timeTimeType
	sqlNullBoolType := reflect.TypeOf(sql.NullBool{})
	datatype[sqlNullBoolType.PkgPath()+ORMDB_SEPERATOR+sqlNullBoolType.Name()] = sqlNullBoolType
	sqlNullFloat64 := reflect.TypeOf(sql.NullFloat64{})
	datatype[sqlNullFloat64.PkgPath()+ORMDB_SEPERATOR+sqlNullFloat64.Name()] = sqlNullFloat64
	sqlNullInt64 := reflect.TypeOf(sql.NullInt64{})
	datatype[sqlNullInt64.PkgPath()+ORMDB_SEPERATOR+sqlNullInt64.Name()] = sqlNullInt64
	sqlNullString := reflect.TypeOf(sql.NullString{})
	datatype[sqlNullString.PkgPath()+ORMDB_SEPERATOR+sqlNullString.Name()] = sqlNullString
}

type EntityFieldDefinition struct {
	Name         string            `json:"name"`
	Kind         reflect.Kind      `json:"kind"`
	ElemKind     reflect.Kind      `json:"value_kind"`
	IsPtr        bool              `json:"is_ptr"`
	IsElemPtr    bool              `json:"is_elem_ptr"`
	Tag          reflect.StructTag `json:"tag"`
	IsEntity     bool              `json:"is_entity"`
	EntityName   string            `json:"entity_name"`
	IsDataType   bool              `json:"is_data_type"`
	DatatypeName string            `json:"datatype_name"`
}

// Builder is the interface that builds a dynamic and runtime struct.
type Builder interface {
	AddEntityFieldDefinition(definition []*EntityFieldDefinition, registry map[string]DynamicStruct) Builder
	Remove(name string) Builder
	Exists(name string) bool
	NumField() int
	Build() DynamicStruct
	BuildNonPtr() DynamicStruct
}

// BuilderImpl is the default Builder implementation.
type BuilderImpl struct {
	fields map[string]reflect.Type
	tags   map[string]reflect.StructTag
}

// NewBuilder returns a concrete Builder
func NewBuilder() Builder {
	return &BuilderImpl{fields: map[string]reflect.Type{}, tags: map[string]reflect.StructTag{}}
}

func RegisterEntity(model interface{}) (string, []*EntityFieldDefinition, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		return "", nil, errors.New("model must be a pointer")
	}
	modelElem := modelType.Elem()
	if modelElem.Kind() != reflect.Struct {
		return "", nil, errors.New("model value must be a struct")
	}
	key := modelElem.Name()
	var entityFieldDefinitions []*EntityFieldDefinition
	for i := 0; i < modelElem.NumField(); i++ {
		field := modelElem.Field(i)
		entityFieldDefinition := &EntityFieldDefinition{}
		entityFieldDefinition.Name = field.Name
		entityFieldDefinition.Tag = field.Tag
		entityTag := entityFieldDefinition.Tag.Get("ormdb")
		switch field.Type.Kind() {
		case reflect.Ptr:
			entityFieldDefinition.IsPtr = true
			name := field.Type.Elem().Name()
			if strings.ToUpper(entityTag) == "DATATYPE" {
				pkgPath := field.Type.Elem().PkgPath()
				fullName := pkgPath + ORMDB_SEPERATOR + name
				entityFieldDefinition.IsDataType = true
				entityFieldDefinition.DatatypeName = fullName
			}
			if strings.ToUpper(entityTag) == "ENTITY" {
				entityFieldDefinition.IsEntity = true
				entityFieldDefinition.EntityName = name
			}
			if field.Type.Elem().Kind() == reflect.Interface || field.Type.Elem().Kind() == reflect.Array || field.Type.Elem().Kind() == reflect.Chan || field.Type.Elem().Kind() == reflect.Func || field.Type.Elem().Kind() == reflect.Map || field.Type.Elem().Kind() == reflect.Uintptr || field.Type.Elem().Kind() == reflect.UnsafePointer || field.Type.Elem().Kind() == reflect.Complex64 || field.Type.Elem().Kind() == reflect.Complex128 {
				return "", nil, errors.New("not supported field kind for orm entity")
			}
			entityFieldDefinition.ElemKind = field.Type.Elem().Kind()
			entityFieldDefinition.Kind = reflect.Ptr
		case reflect.Struct:
			name := field.Type.Name()
			if strings.ToUpper(entityTag) == "DATATYPE" {
				pkgPath := field.Type.PkgPath()
				fullName := pkgPath + ORMDB_SEPERATOR + name
				entityFieldDefinition.IsDataType = true
				entityFieldDefinition.DatatypeName = fullName
			}
			if strings.ToUpper(entityTag) == "ENTITY" {
				entityFieldDefinition.IsEntity = true
				entityFieldDefinition.EntityName = name
			}

			entityFieldDefinition.Kind = reflect.Struct
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Ptr {
				name := field.Type.Elem().Elem().Name()
				if strings.ToUpper(entityTag) == "DATATYPE" {
					pkgPath := field.Type.Elem().PkgPath()
					fullName := pkgPath + ORMDB_SEPERATOR + name
					entityFieldDefinition.IsDataType = true
					entityFieldDefinition.DatatypeName = fullName
				}
				if strings.ToUpper(entityTag) == "ENTITY" {
					entityFieldDefinition.IsEntity = true
					entityFieldDefinition.EntityName = name
				}
				if field.Type.Elem().Elem().Kind() == reflect.Interface || field.Type.Elem().Elem().Kind() == reflect.Array || field.Type.Elem().Elem().Kind() == reflect.Chan || field.Type.Elem().Elem().Kind() == reflect.Func || field.Type.Elem().Elem().Kind() == reflect.Map || field.Type.Elem().Elem().Kind() == reflect.Uintptr || field.Type.Elem().Elem().Kind() == reflect.UnsafePointer || field.Type.Elem().Elem().Kind() == reflect.Complex64 || field.Type.Elem().Elem().Kind() == reflect.Complex128 {
					return "", nil, errors.New("not supported field kind for orm entity")
				}
				entityFieldDefinition.IsElemPtr = true
				entityFieldDefinition.ElemKind = field.Type.Elem().Elem().Kind()
				entityFieldDefinition.Kind = reflect.Slice
			} else if field.Type.Elem().Kind() == reflect.Struct {
				name := field.Type.Elem().Name()
				if strings.ToUpper(entityTag) == "DATATYPE" {
					pkgPath := field.Type.Elem().PkgPath()
					fullName := pkgPath + ORMDB_SEPERATOR + name
					entityFieldDefinition.IsDataType = true
					entityFieldDefinition.DatatypeName = fullName
				}
				if strings.ToUpper(entityTag) == "ENTITY" {
					entityFieldDefinition.IsEntity = true
					entityFieldDefinition.EntityName = name
				}

				entityFieldDefinition.ElemKind = field.Type.Elem().Kind()
				entityFieldDefinition.Kind = reflect.Slice
			} else if field.Type.Elem().Kind() == reflect.Interface || field.Type.Elem().Kind() == reflect.Array || field.Type.Elem().Kind() == reflect.Chan || field.Type.Elem().Kind() == reflect.Func || field.Type.Elem().Kind() == reflect.Map || field.Type.Elem().Kind() == reflect.Uintptr || field.Type.Elem().Kind() == reflect.UnsafePointer || field.Type.Elem().Kind() == reflect.Complex64 || field.Type.Elem().Kind() == reflect.Complex128 {
				return "", nil, errors.New("not supported field kind for orm entity")
			} else {
				entityFieldDefinition.ElemKind = field.Type.Elem().Kind()
				entityFieldDefinition.Kind = reflect.Slice
			}
		case reflect.Interface, reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Uintptr, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
			return "", nil, errors.New("not supported field kind for orm entity")
		default:
			entityFieldDefinition.Kind = field.Type.Kind()
		}

		entityFieldDefinitions = append(entityFieldDefinitions, entityFieldDefinition)
	}

	return key, entityFieldDefinitions, nil
}

func (b *BuilderImpl) AddEntityFieldDefinition(efds []*EntityFieldDefinition, registry map[string]DynamicStruct) Builder {
	for _, efd := range efds {
		var fieldType reflect.Type
		switch efd.Kind {
		case reflect.Slice:
			var elemType reflect.Type
			if efd.IsDataType {
				elemType = datatype[efd.DatatypeName]
			}
			if efd.IsEntity {
				elemType = registry[efd.EntityName].StructType()
			}
			if !efd.IsDataType && !efd.IsEntity {
				switch efd.ElemKind {
				case reflect.Uint:
					elemType = reflect.TypeOf(SampleUint)
				case reflect.Uint8:
					elemType = reflect.TypeOf(SampleUint8)
				case reflect.Uint16:
					elemType = reflect.TypeOf(SampleUint16)
				case reflect.Uint32:
					elemType = reflect.TypeOf(SampleUint32)
				case reflect.Uint64:
					elemType = reflect.TypeOf(SampleUint64)
				case reflect.Int:
					elemType = reflect.TypeOf(SampleInt)
				case reflect.Int8:
					elemType = reflect.TypeOf(SampleInt8)
				case reflect.Int16:
					elemType = reflect.TypeOf(SampleInt16)
				case reflect.Int32:
					elemType = reflect.TypeOf(SampleInt32)
				case reflect.Int64:
					elemType = reflect.TypeOf(SampleInt64)
				case reflect.String:
					elemType = reflect.TypeOf(SampleString)
				case reflect.Bool:
					elemType = reflect.TypeOf(SampleBool)
				case reflect.Float32:
					elemType = reflect.TypeOf(SampleFloat32)
				case reflect.Float64:
					elemType = reflect.TypeOf(SampleFloat64)
				default:
					panic(errors.New("not supported data type"))
				}
			}
			if efd.IsElemPtr {
				elemType = reflect.PtrTo(elemType)
			}
			fieldType = reflect.SliceOf(elemType)
		case reflect.Ptr:
			var elemType reflect.Type
			if efd.IsDataType {
				elemType = datatype[efd.DatatypeName]
			}
			if efd.IsEntity {
				elemType = registry[efd.EntityName].StructType()
			}

			if !efd.IsDataType && !efd.IsEntity {
				switch efd.ElemKind {
				case reflect.Uint:
					elemType = reflect.TypeOf(SampleUint)
				case reflect.Uint8:
					elemType = reflect.TypeOf(SampleUint8)
				case reflect.Uint16:
					elemType = reflect.TypeOf(SampleUint16)
				case reflect.Uint32:
					elemType = reflect.TypeOf(SampleUint32)
				case reflect.Uint64:
					elemType = reflect.TypeOf(SampleUint64)
				case reflect.Int:
					elemType = reflect.TypeOf(SampleInt)
				case reflect.Int8:
					elemType = reflect.TypeOf(SampleInt8)
				case reflect.Int16:
					elemType = reflect.TypeOf(SampleInt16)
				case reflect.Int32:
					elemType = reflect.TypeOf(SampleInt32)
				case reflect.Int64:
					elemType = reflect.TypeOf(SampleInt64)
				case reflect.String:
					elemType = reflect.TypeOf(SampleString)
				case reflect.Bool:
					elemType = reflect.TypeOf(SampleBool)
				case reflect.Float32:
					elemType = reflect.TypeOf(SampleFloat32)
				case reflect.Float64:
					elemType = reflect.TypeOf(SampleFloat64)
				default:
					panic(errors.New("not supported data type"))
				}
			}
			if efd.IsPtr {
				fieldType = reflect.PtrTo(elemType)
			}
		case reflect.Struct:
			if efd.IsDataType {
				fieldType = datatype[efd.DatatypeName]
			}
			if efd.IsEntity {
				fieldType = registry[efd.EntityName].StructType()
			}
		case reflect.Uint:
			fieldType = reflect.TypeOf(SampleUint)
		case reflect.Uint8:
			fieldType = reflect.TypeOf(SampleUint8)
		case reflect.Uint16:
			fieldType = reflect.TypeOf(SampleUint16)
		case reflect.Uint32:
			fieldType = reflect.TypeOf(SampleUint32)
		case reflect.Uint64:
			fieldType = reflect.TypeOf(SampleUint64)
		case reflect.Int:
			fieldType = reflect.TypeOf(SampleInt)
		case reflect.Int8:
			fieldType = reflect.TypeOf(SampleInt8)
		case reflect.Int16:
			fieldType = reflect.TypeOf(SampleInt16)
		case reflect.Int32:
			fieldType = reflect.TypeOf(SampleInt32)
		case reflect.Int64:
			fieldType = reflect.TypeOf(SampleInt64)
		case reflect.String:
			fieldType = reflect.TypeOf(SampleString)
		case reflect.Bool:
			fieldType = reflect.TypeOf(SampleBool)
		case reflect.Float32:
			fieldType = reflect.TypeOf(SampleFloat32)
		case reflect.Float64:
			fieldType = reflect.TypeOf(SampleFloat64)
		default:
			panic(errors.New("not supported data type"))
		}

		b.fields[efd.Name] = fieldType
		b.tags[efd.Name] = efd.Tag
	}

	return b
}

// Remove returns a Builder that was removed a field named by name parameter.
func (b *BuilderImpl) Remove(name string) Builder {
	delete(b.fields, name)
	return b
}

// Exists returns true if the specified name field exists
func (b *BuilderImpl) Exists(name string) bool {
	_, ok := b.fields[name]
	return ok
}

// NumField returns the number of built struct fields.
func (b *BuilderImpl) NumField() int {
	return len(b.fields)
}

// Build returns a concrete struct pointer built by Builder.
func (b *BuilderImpl) Build() DynamicStruct {
	return b.build(true)
}

// BuildNonPtr returns a concrete struct built by Builder.
func (b *BuilderImpl) BuildNonPtr() DynamicStruct {
	return b.build(false)
}

func (b *BuilderImpl) build(isPtr bool) DynamicStruct {
	var i int
	fs := make([]reflect.StructField, len(b.fields))
	for name, typ := range b.fields {
		fs[i] = reflect.StructField{Name: name, Type: typ, Tag: b.tags[name]}
		i++
	}

	return newDs(fs, isPtr)
}

// DynamicStruct is the interface that built dynamic struct by Builder.Build().
type DynamicStruct interface {
	NumField() int
	Field(i int) reflect.StructField
	FieldByName(name string) (reflect.StructField, bool)
	IsPtr() bool
	Interface() interface{}
	StructType() reflect.Type
	//DecodeMap(m map[string]interface{}) (interface{}, error)
}

// Impl is the default DynamicStruct implementation.
type Impl struct {
	structType reflect.Type
	isPtr      bool
	intf       interface{}
}

func newDs(fs []reflect.StructField, isPtr bool) DynamicStruct {
	ds := &Impl{structType: reflect.StructOf(fs), isPtr: isPtr}

	n := reflect.New(ds.structType)
	if isPtr {
		ds.intf = n.Interface()
	} else {
		ds.intf = reflect.Indirect(n).Interface()
	}

	return ds
}

// NumField returns the number of built struct fields.
func (ds *Impl) NumField() int {
	return ds.structType.NumField()
}

// Field returns the i'th field of the built struct.
func (ds *Impl) Field(i int) reflect.StructField {
	return ds.structType.Field(i)
}

// FieldByName returns the struct field with the given name
// and a boolean indicating if the field was found.
func (ds *Impl) FieldByName(name string) (reflect.StructField, bool) {
	return ds.structType.FieldByName(name)
}

// IsPtr reports whether the built struct type is pointer.
func (ds *Impl) IsPtr() bool {
	return ds.isPtr
}

// Interface returns the interface of built struct.
func (ds *Impl) Interface() interface{} {
	return ds.intf
}

func (ds *Impl) StructType() reflect.Type {
	return ds.structType
}