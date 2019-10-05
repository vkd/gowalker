package gowalker

import (
	"errors"
	"reflect"
)

// Walk - walk struct by all public fields
func Walk(value interface{}, walker Walker) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return errors.New("unsupported type for value: allowed only ptr")
	}
	_, err := walkPrt(v, emptyField, walker)
	return err
}

// Walker - interface to walk through struct fields
type Walker interface {
	Step(value reflect.Value, field reflect.StructField) (bool, error)
}

// Wrapper - interface to allow build tree of struct fields
type Wrapper interface {
	Walker
	Wrap(field reflect.StructField) Walker
}

// WalkerFunc - func implemented Walk interface
type WalkerFunc func(value reflect.Value, field reflect.StructField) (bool, error)

// Step - one step of walker
func (f WalkerFunc) Step(value reflect.Value, field reflect.StructField) (bool, error) {
	return f(value, field)
}

func walk(value reflect.Value, field reflect.StructField, walker Walker) (bool, error) {
	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, walker)
	}

	if !isEmptyField(field) {
		ok, err := walker.Step(value, field)
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

func walkPrt(value reflect.Value, field reflect.StructField, walker Walker) (setted bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	setted, err = walk(vPtr.Elem(), field, walker)
	if err != nil {
		return false, err
	}
	if isCreateNew && setted {
		value.Set(vPtr)
	}
	return setted, nil
}

func walkStruct(value reflect.Value, field reflect.StructField, walker Walker) (setted bool, err error) {
	tp := value.Type()

	var isStructSetted bool
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		var nextW = walker
		if ww, ok := walker.(Wrapper); ok && !isEmptyField(field) {
			nextW = ww.Wrap(field)
		}
		tField := tp.Field(i)
		setted, err := walk(value.Field(i), tField, nextW)
		if err != nil {
			return false, err
		}
		isStructSetted = isStructSetted || setted
	}
	return isStructSetted, nil
}

var emptyField = reflect.StructField{}

func isEmptyField(field reflect.StructField) bool {
	return field.Name == emptyField.Name
}
