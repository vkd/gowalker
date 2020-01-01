package gowalker

import (
	"errors"
	"reflect"
)

func walkIface(value interface{}, walker Walker) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return false, errors.New("unsupported type for value: allowed only ptr")
	}
	return walkPrt(v, emptyField, walker)
}

func walk(value reflect.Value, field reflect.StructField, walker Walker) (bool, error) {
	if !value.CanSet() {
		return false, nil
	}

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
