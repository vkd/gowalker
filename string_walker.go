package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// NewStringWalker - simple string walker
func NewStringWalker(tag string, source Sourcer, namer Namer) Walker {
	if ss, ok := source.(SliceSourcer); ok {
		return WalkerFunc(func(value reflect.Value, field reflect.StructField, name Name) (bool, error) {
			return SliceStringsWalkerStep(tag, ss, value, field, name, namer)
		})
	}
	return WalkerFunc(func(value reflect.Value, field reflect.StructField, name Name) (bool, error) {
		return StringWalkerStep(tag, source, value, field, name, namer)
	})
}

// StringWalkerStep - step of walker by string value
func StringWalkerStep(tag string, source Sourcer, value reflect.Value, field reflect.StructField, name Name, namer Namer) (bool, error) {
	str, ok, err := StringGetValue(tag, source, field, name, namer)
	if err != nil || !ok {
		return ok, err
	}
	return true, setter.SetString(value, field, str)
}

// StringGetValue - get string value from field
func StringGetValue(tag string, source Sourcer, field reflect.StructField, name Name, namer Namer) (string, bool, error) {
	t, ok := field.Tag.Lookup(tag)
	if !ok {
		t = name.Get(namer)
	}
	return source.Get(t)
}
