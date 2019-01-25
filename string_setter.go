package gowalker

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)

// ErrUnknownKind - unknown kind of value
var ErrUnknownKind = errors.New("unknown kind")

// ErrNotSupportKind - kind not support to set by string
var ErrNotSupportKind = errors.New("not support kind")

// ErrNotImplemented - not implemented yet functionality
var ErrNotImplemented = errors.New("not implemented")

// SetValueBySliceOfString - set value by slice of strings
//
// Support multi values only: Slice, Array
// other kindes setted by string (first value from slice)
func SetValueBySliceOfString(value reflect.Value, field reflect.StructField, strs []string) error {
	switch value.Kind() {
	case reflect.Slice:
		return sliceStringSetter(value, field, strs)
	case reflect.Array:
		return arrayStringSetter(value, field, strs)
	default:
		var str string
		if len(strs) > 0 {
			str = strs[0]
		}
		return SetValueByString(value, field, str)
	}
}

// SetValueByString - set value by string
//
// Not implemented kinds: Complex, Chan
func SetValueByString(value reflect.Value, field reflect.StructField, str string) error { // nolint: gocyclo
	switch value.Kind() {
	case reflect.Bool:
		return boolStringSetter(value, field, str)
	// Int
	case reflect.Int:
		return int0StrSetter(value, field, str)
	case reflect.Int8:
		return int8StrSetter(value, field, str)
	case reflect.Int16:
		return int16StrSetter(value, field, str)
	case reflect.Int32:
		return int32StrSetter(value, field, str)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return timeDurationStringSetter(value, field, str)
		}
		return int64StrSetter(value, field, str)

	// Uint
	case reflect.Uint:
		return uint0StrSetter(value, field, str)
	case reflect.Uint8:
		return uint8StrSetter(value, field, str)
	case reflect.Uint16:
		return uint16StrSetter(value, field, str)
	case reflect.Uint32:
		return uint32StrSetter(value, field, str)
	case reflect.Uint64:
		return uint64StrSetter(value, field, str)

	// Float
	case reflect.Float32:
		return float32StrSetter(value, field, str)
	case reflect.Float64:
		return float64StrSetter(value, field, str)

	// Complex
	case reflect.Complex64:
		return ErrNotImplemented
	case reflect.Complex128:
		return ErrNotImplemented

	case reflect.Array:
		return arrayStringSetter(value, field, []string{str})

	case reflect.Chan:
		return ErrNotImplemented

	case reflect.Map:
		return json.Unmarshal([]byte(str), value.Addr().Interface())

	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()).Elem())
		}
		return SetValueByString(value.Elem(), field, str)

	case reflect.Slice:
		return sliceStringSetter(value, field, []string{str})

	case reflect.String:
		value.SetString(str)
		return nil

	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return timeStringSetter(value, field, str)
		}
		return json.Unmarshal([]byte(str), value.Addr().Interface())

	case
		reflect.Uintptr,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer:
		return ErrNotSupportKind

	default:
		return ErrUnknownKind
	}
}

type stringSetterFunc func(value reflect.Value, field reflect.StructField, str string) error

func boolStringSetter(value reflect.Value, field reflect.StructField, s string) error {
	if s == "" {
		s = "false"
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	value.SetBool(b)
	return nil
}

func intStringSetterFunc(bitSize int) stringSetterFunc {
	return func(value reflect.Value, field reflect.StructField, s string) error {
		if s == "" {
			s = "0"
		}
		x, err := strconv.ParseInt(s, 10, bitSize)
		if err != nil {
			return err
		}
		value.SetInt(x)
		return nil
	}
}

var (
	int0StrSetter  = intStringSetterFunc(0)
	int8StrSetter  = intStringSetterFunc(8)
	int16StrSetter = intStringSetterFunc(16)
	int32StrSetter = intStringSetterFunc(32)
	int64StrSetter = intStringSetterFunc(64)
)

func timeDurationStringSetter(value reflect.Value, field reflect.StructField, s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func uintStringSetterFunc(bitSize int) stringSetterFunc {
	return func(value reflect.Value, field reflect.StructField, s string) error {
		if s == "" {
			s = "0"
		}
		x, err := strconv.ParseUint(s, 10, bitSize)
		if err != nil {
			return err
		}
		value.SetUint(x)
		return nil
	}
}

var (
	uint0StrSetter  = uintStringSetterFunc(0)
	uint8StrSetter  = uintStringSetterFunc(8)
	uint16StrSetter = uintStringSetterFunc(16)
	uint32StrSetter = uintStringSetterFunc(32)
	uint64StrSetter = uintStringSetterFunc(64)
)

func floatStringSetterFunc(bitSize int) stringSetterFunc {
	return func(value reflect.Value, field reflect.StructField, s string) error {
		if s == "" {
			s = "0.0"
		}
		x, err := strconv.ParseFloat(s, bitSize)
		if err != nil {
			return err
		}
		value.SetFloat(x)
		return nil
	}
}

var (
	float32StrSetter = floatStringSetterFunc(32)
	float64StrSetter = floatStringSetterFunc(64)
)

func sliceStringSetter(value reflect.Value, field reflect.StructField, strs []string) error {
	slice := reflect.MakeSlice(value.Type(), len(strs), len(strs))
	err := arrayStringSetter(slice, field, strs)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func arrayStringSetter(value reflect.Value, field reflect.StructField, strs []string) error {
	for i, s := range strs {
		err := SetValueByString(value.Index(i), field, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func timeStringSetter(value reflect.Value, field reflect.StructField, s string) error {
	layout := field.Tag.Get("time_layout")
	if layout == "" {
		layout = time.RFC3339
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}
