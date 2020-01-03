package gowalker

import (
	"errors"
	"reflect"
)

// Walker - interface to walk through struct fields
type Walker interface {
	Step(value reflect.Value, field reflect.StructField) (bool, error)
}

// Namer - interface to get a full field name.
type Namer interface {
	FieldName(parent, name string) string
}

func walkIface(value interface{}, walker Walker, namer Namer) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return false, errors.New("unsupported type for value: allowed only ptr")
	}
	return walkPrt(v, structField{namer: namer}, walker)
}

func walk(value reflect.Value, field structField, walker Walker) (bool, error) {
	if !value.CanSet() {
		return false, nil
	}

	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, walker)
	}

	if !isEmptyField(field.field) {
		ok, err := walker.Step(value, field.field)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if kind == reflect.Struct {
		return walkStruct(value, field, walker)
	}

	return false, nil
}

func walkPrt(value reflect.Value, field structField, walker Walker) (set bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	set, err = walk(vPtr.Elem(), field, walker)
	if err != nil {
		return false, err
	}
	if isCreateNew && set {
		value.Set(vPtr)
	}
	return set, nil
}

func walkStruct(value reflect.Value, field structField, walker Walker) (set bool, err error) {
	tp := value.Type()

	var isStructSet bool
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		tField := tp.Field(i)
		set, err := walk(value.Field(i), field.Child(tField), walker)
		if err != nil {
			return false, err
		}
		isStructSet = isStructSet || set
	}
	return isStructSet, nil
}

var emptyField = reflect.StructField{}

func isEmptyField(field reflect.StructField) bool {
	return field.Name == emptyField.Name
}

type structField struct {
	field reflect.StructField

	namer Namer
}

func (f structField) Child(c reflect.StructField) structField {
	if f.namer != nil {
		c.Name = f.namer.FieldName(f.field.Name, c.Name)
	}

	f.field = c
	return f
}
