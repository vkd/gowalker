package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

func StringSetter(fname FieldKeyer, source StringSource) Walker {
	return WalkerFunc(func(value reflect.Value, field reflect.StructField, name Name) (bool, error) {
		key, ok := fname.GetFieldKey(field, name)
		if !ok {
			return false, nil
		}
		s, ok, err := source.Get(key)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}

		return true, setter.SetString(value, field, s)
	})
}

type FieldKeyer interface {
	GetFieldKey(field reflect.StructField, name Name) (string, bool)
}

func StringsSetter(fname FieldKeyer, source StringsSource) Walker {
	return WalkerFunc(func(value reflect.Value, field reflect.StructField, name Name) (bool, error) {
		key, ok := fname.GetFieldKey(field, name)
		if !ok {
			return false, nil
		}
		ss, ok, err := source.GetStrings(key)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}

		return true, setter.SetSliceStrings(value, field, ss)
	})
}

func FieldKeys(keyers ...FieldKeyer) FieldKeyer {
	return fieldKeyers(keyers)
}

type fieldKeyers []FieldKeyer

var _ FieldKeyer = (fieldKeyers)(nil)

func (f fieldKeyers) GetFieldKey(field reflect.StructField, name Name) (string, bool) {
	for _, fkey := range f {
		key, ok := fkey.GetFieldKey(field, name)
		if ok {
			return key, true
		}
	}
	return "", false
}
