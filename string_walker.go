package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// NewStringWalker - simple string walker
func NewStringWalker(tag string, source Sourcer) Walker {
	if ss, ok := source.(SliceSourcer); ok {
		return WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
			return SliceStringsWalkerStep(tag, ss, value, field)
		})
	}
	return WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
		return StringWalkerStep(tag, source, value, field)
	})
}

// StringWalkerStep - step of walker by string value
func StringWalkerStep(tag string, source Sourcer, value reflect.Value, field reflect.StructField) (bool, error) {
	str, ok, err := StringGetValue(tag, source, field)
	if err != nil || !ok {
		return ok, err
	}
	return true, setter.SetString(value, field, str)
}

// StringGetValue - get string value from field
func StringGetValue(tag string, source Sourcer, field reflect.StructField) (string, bool, error) {
	t := TagStringParse(field, tag)
	str, ok, err := source.Get(t.Value)
	if err != nil {
		return "", false, err
	}
	if !ok && !t.IsDefaultValue {
		return "", false, nil
	}
	if !ok {
		str = t.DefaultValue
	}
	return str, true, nil
}
