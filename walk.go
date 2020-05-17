package gowalker

import (
	"errors"
	"reflect"
)

// Walker - interface to walk through struct fields
type Walker interface {
	Step(value reflect.Value, field reflect.StructField, name Name) (bool, error)
}

// Name of the struct field.
type Name interface {
	Get(n Namer) string
}

// WalkerFunc - func implemented Walk interface
type WalkerFunc func(value reflect.Value, field reflect.StructField, name Name) (bool, error)

// Step - one step of walker
func (f WalkerFunc) Step(value reflect.Value, field reflect.StructField, name Name) (bool, error) {
	return f(value, field, name)
}

// Walk - walk struct by all public fields
func Walk(value interface{}, walker Walker) error {
	var name = make(fieldsStack, 0, 4)
	_, err := walkIface(value, &name, walker)
	return err
}

func walkIface(value interface{}, name nameBuilder, walker Walker) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return false, errors.New("unsupported type for value: allowed only ptr")
	}
	return walkPrt(v, emptyField, name, walker)
}

func walk(value reflect.Value, field reflect.StructField, name nameBuilder, walker Walker) (bool, error) {
	if !value.CanSet() {
		return false, nil
	}

	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, name, walker)
	}

	if !isEmptyField(field) {
		ok, err := walker.Step(value, field, name)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if kind == reflect.Struct {
		return walkStruct(value, field, name, walker)
	}

	return false, nil
}

func walkPrt(value reflect.Value, field reflect.StructField, name nameBuilder, walker Walker) (set bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	set, err = walk(vPtr.Elem(), field, name, walker)
	if err != nil {
		return false, err
	}
	if isCreateNew && set {
		value.Set(vPtr)
	}
	return set, nil
}

func walkStruct(value reflect.Value, _ reflect.StructField, name nameBuilder, walker Walker) (set bool, err error) {
	tp := value.Type()

	var isStructSet bool
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		tField := tp.Field(i)
		name.Next(tField)
		set, err := walk(value.Field(i), tField, name, walker)
		name.Pop()
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

type nameBuilder interface {
	Name

	Next(reflect.StructField)
	Pop()
}

type fieldsStack []string

func (fs *fieldsStack) Get(n Namer) string {
	if n == nil {
		return (*fs)[len(*fs)-1]
	}
	var out string
	for _, f := range *fs {
		out = n.FieldName(out, f)
	}
	return out
}

func (fs *fieldsStack) Next(f reflect.StructField) {
	*fs = append(*fs, f.Name)
}

func (fs *fieldsStack) Pop() {
	*fs = (*fs)[:len(*fs)-1]
}
