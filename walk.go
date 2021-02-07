package gowalker

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrUnsupportedValue is raised if value is passed not as a pointer.
var ErrUnsupportedValue = errors.New("unsupported type for value: allowed only ptr")

// Setter - interface to walk through struct fields.
type Setter interface {
	TrySet(reflect.Value, reflect.StructField, Fields) (ok bool, _ error)
}

// SetterFunc - func implemented Walk interface.
type SetterFunc func(reflect.Value, reflect.StructField, Fields) (bool, error)

var _ Setter = SetterFunc(nil)

// Set - one step of walker.
func (f SetterFunc) TrySet(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	return f(value, field, fs)
}

// MultiSetterOR - set value only one (first returns true).
func MultiSetterOR(ss ...Setter) Setter {
	if len(ss) == 1 {
		return ss[0]
	}
	return SetterFunc(func(v reflect.Value, sf reflect.StructField, f Fields) (bool, error) {
		for _, s := range ss {
			ok, err := s.TrySet(v, sf, f)
			if err != nil {
				return ok, fmt.Errorf("setter %T: %w", s, err)
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	})
}

// Walk - walk struct by all public fields.
//
// Should be passed as pointer:
// type myStruct struct
// var s myStruct
// gowalker.Walk(&s, ...)
func Walk(value interface{}, fs Fields, ss ...Setter) error {
	for _, s := range ss {
		_, err := walkIface(value, s, fs)
		if err != nil {
			if len(ss) == 1 {
				return err
			}
			return fmt.Errorf("setter %T: %w", s, err)
		}
	}
	return nil
}

func WalkFast(value interface{}, ws ...Setter) error {
	return Walk(value, nil, ws...)
}

func walkIface(value interface{}, s Setter, fs Fields) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return false, ErrUnsupportedValue
	}
	if s == nil {
		return false, nil
	}
	return walkPrt(v, emptyField, s, fs)
}

func walk(value reflect.Value, field reflect.StructField, s Setter, fs Fields) (set bool, _ error) {
	if !value.CanSet() {
		return false, nil
	}

	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, s, fs)
	}

	if !isEmptyField(field) && s != nil {
		ok, err := s.TrySet(value, field, fs)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if kind == reflect.Struct {
		return walkStruct(value, s, fs)
	}

	return false, nil
}

func walkPrt(value reflect.Value, field reflect.StructField, s Setter, fs Fields) (set bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	set, err = walk(vPtr.Elem(), field, s, fs)
	if err != nil {
		return false, err
	}
	if isCreateNew && set {
		value.Set(vPtr)
	}
	return set, nil
}

func walkStruct(value reflect.Value, s Setter, fs Fields) (set bool, err error) {
	tp := value.Type()

	var isStructSet bool
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		tField := tp.Field(i)
		var nextFs Fields
		if fs != nil {
			nextFs = append(fs, tField)
		}
		set, err := walk(value.Field(i), tField, s, nextFs)
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

type Fields []reflect.StructField

func MakeFields(cap int) Fields {
	return make(Fields, 0, cap)
}
