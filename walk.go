package gowalker

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/vkd/gowalker/setter"
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
func Walk(value interface{}, fs Fields, w Walker, cs ...Option) error {
	var cfg config
	for _, c := range cs {
		c.apply(&cfg)
	}

	return walkIface(value, w, fs, cfg)
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) { f(cfg) }

type config struct {
	StepOnStructFields bool
}

func StepOnStructFields() Option {
	return optionFunc(func(c *config) {
		c.StepOnStructFields = true
	})
}

func walkIface(value interface{}, w Walker, fs Fields, cfg config) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return ErrUnsupportedValue
	}
	if w == nil {
		return nil
	}
	_, _, err := walkPrt(v, emptyField, w, fs, cfg)
	return err
}

func walk(value reflect.Value, field reflect.StructField, w Walker, fs Fields, cfg config) (stepped bool, stop bool, _ error) {
	if !value.CanSet() {
		return false, false, nil
	}

	kind := value.Kind()
	if kind == reflect.Ptr {
		return walkPrt(value, field, w, fs, cfg)
	}

	if kind == reflect.Struct {
		var err error
		stepped, stop, err = walkStruct(value, w, fs, cfg)
		if err != nil {
			return false, false, fmt.Errorf("struct %q: %w", field.Name, err)
		}
		if stop {
			return stepped, true, nil
		}
	}

	if !isEmptyField(field) && w != nil {
		var err error
		stepped, stop, err = step(value, field, w, fs, cfg, stepped)
		if err != nil {
			return false, false, err
		}
		return stepped, stop, nil
	}

	return stepped, false, nil
}

var setStringerType = reflect.TypeOf((*setter.SetStringer)(nil)).Elem()

func step(value reflect.Value, field reflect.StructField, w Walker, fs Fields, cfg config, stepped bool) (steppedOut bool, stop bool, _ error) {
	if field.Tag.Get("walker") == "embed" {
		return stepped, false, nil
	}

	switch field.Type.Kind() {
	case reflect.Struct:
		iface := value.Interface()
		_, isTime := iface.(time.Time)

		switch {
		case isTime:
		case value.Type().Implements(setStringerType):
		case value.CanAddr() && value.Addr().Type().Implements(setStringerType):
		default:
			if stepped && !cfg.StepOnStructFields {
				return stepped, false, nil
			}
		}
	default:
	}

	stop, err := w.Step(value, field, fs)
	if err != nil {
		return false, false, err
	}

	return true, stop, nil
}

func walkPrt(value reflect.Value, field reflect.StructField, w Walker, fs Fields, cfg config) (stepped bool, stop bool, err error) {
	isCreateNew := value.IsNil()

	vPtr := value
	if isCreateNew {
		vPtr = reflect.New(value.Type().Elem())
	}
	fPtr := field
	if isCreateNew {
		fPtr.Type = field.Type.Elem()
	}
	stepped, stop, err = walk(vPtr.Elem(), fPtr, w, fs, cfg)
	if err != nil {
		return false, false, err
	}
	if isCreateNew && stop {
		value.Set(vPtr)
	}
	return stepped, stop, nil
}

func walkStruct(value reflect.Value, w Walker, fs Fields, cfg config) (stepped bool, set bool, _ error) {
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
		var err error
		stepped, set, err = walk(value.Field(i), tField, w, nextFs, cfg)
		if err != nil {
			return false, false, fmt.Errorf("field %q: %w", tField.Name, err)
		}
		isStructSet = isStructSet || set
	}
	return stepped, isStructSet, nil
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
