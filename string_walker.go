package gowalker

import (
	"os"
	"reflect"

	"github.com/vkd/gowalker/setters"
)

// NewStringWalker - simple string walker
func NewStringWalker(tag string, source StringSource) Walker {
	if ss, ok := source.(SliceStringsSource); ok {
		return WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
			return SliceStringsWalkerStep(tag, ss, value, field)
		})
	}
	return WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
		return StringWalkerStep(tag, source, value, field)
	})
}

// StringSource - source of values by string
type StringSource interface {
	Get(key string) (value string, ok bool, err error)
}

// StringWalkerStep - step of walker by string value
func StringWalkerStep(tag string, source StringSource, value reflect.Value, field reflect.StructField) (bool, error) {
	str, ok, err := StringGetValue(tag, source, field)
	if err != nil || !ok {
		return ok, err
	}
	return true, setters.SetString(value, field, str)
}

// StringGetValue - get string value from field
func StringGetValue(tag string, source StringSource, field reflect.StructField) (string, bool, error) {
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

// MapStringSource - map[string]string implement StringSource
type MapStringSource map[string]string

// Get value from source
func (s MapStringSource) Get(key string) (string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}

// StringSourceFunc - function implement string source
type StringSourceFunc func(key string) (value string, ok bool, err error)

// Get value from source
func (s StringSourceFunc) Get(key string) (value string, ok bool, err error) {
	return s(key)
}

// EnvStringSource - source from os env
var EnvStringSource = StringSourceFunc(func(key string) (string, bool, error) {
	v, ok := os.LookupEnv(key)
	return v, ok, nil
})

// NewEnvWalker - walker from env
func NewEnvWalker(tag string) Walker {
	return NewStringWalker(tag, EnvStringSource)
}
