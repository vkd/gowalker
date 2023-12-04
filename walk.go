package gowalker

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrUnsupportedValue is raised if value is passed not as a pointer.
var ErrUnsupportedValue = errors.New("unsupported type for value: allowed only ptr")

type Walker interface {
	Step(reflect.Value, reflect.StructField, Fields) (stop bool, _ error)
}

type WalkerFunc func(reflect.Value, reflect.StructField, Fields) (stop bool, _ error)

var _ Walker = WalkerFunc(nil)

func (f WalkerFunc) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	return f(value, field, fs)
}

// Walk - walk struct by all public fields.
//
// Should be passed as pointer:
// type myStruct struct
// var s myStruct
// gowalker.Walk(&s, ...)
func Walk(value interface{}, fs Fields, w Walker) error {
	_, err := walkIface(value, w, fs)
	return err
}

func walkIface(value interface{}, w Walker, fs Fields) (bool, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return false, ErrUnsupportedValue
	}
	if w == nil {
		return false, nil
	}
	return walkPrt(v, emptyField, w, fs)
}

func walk(value reflect.Value, field reflect.StructField, w Walker, fs Fields) (set bool, _ error) {
	if !value.CanSet() {
		return false, nil
	}

	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, w, fs)
	}

	if !isEmptyField(field) && w != nil {
		stop, err := w.Step(value, field, fs)
		if err != nil {
			return false, err
		}
		if stop {
			return true, nil
		}
	}

	if kind == reflect.Struct {
		return walkStruct(value, w, fs)
	}

	return false, nil
}

func walkPrt(value reflect.Value, field reflect.StructField, w Walker, fs Fields) (set bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	set, err = walk(vPtr.Elem(), field, w, fs)
	if err != nil {
		return false, err
	}
	if isCreateNew && set {
		value.Set(vPtr)
	}
	return set, nil
}

func walkStruct(value reflect.Value, w Walker, fs Fields) (set bool, err error) {
	tp := value.Type()

	var isStructSet bool
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanSet() {
			continue
		}
		tField := tp.Field(i)

		var nextFs Fields
		switch tField.Tag.Get("walker") {
		case "embed":
			nextFs = fs
		default:
			if fs != nil {
				nextFs = append(fs, tField)
			}
		}
		set, err := walk(value.Field(i), tField, w, nextFs)
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

type WalkersOR []Walker

func (w WalkersOR) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	for _, s := range w {
		stop, err := s.Step(value, field, fs)
		if err != nil {
			return stop, fmt.Errorf("walker %T: %w", s, err)
		}
		if stop {
			return true, nil
		}
	}
	return false, nil
}
