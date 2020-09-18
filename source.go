package gowalker

import (
	"fmt"
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// StringSource - source with strings as values.
type StringSource interface {
	Get(key string) (value string, ok bool, err error)
}

func SetStringSource(value reflect.Value, field reflect.StructField, source StringSource, key string) (bool, error) {
	if source == nil {
		return false, nil
	}
	s, ok, err := source.Get(key)
	if err != nil {
		return false, fmt.Errorf("source get value: %w", err)
	}
	if !ok {
		return false, nil
	}

	return true, setter.SetString(value, field, s)
}

// MapStringSource - map[string]string implement Sourcer.
type MapStringSource map[string]string

var _ StringSource = MapStringSource(nil)

// Get value from source.
func (s MapStringSource) Get(key string) (string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}

// LookupFuncSource - func(key string) (string, bool) sourcer.
type LookupFuncSource func(key string) (string, bool)

var _ StringSource = LookupFuncSource(nil)

// Get value from source.
func (f LookupFuncSource) Get(key string) (string, bool, error) {
	v, ok := f(key)
	return v, ok, nil
}

// StringsSource - source of values by slice of strings.
type StringsSource interface {
	GetStrings(key string) (value []string, ok bool, err error)
}

func SetStringsSource(value reflect.Value, field reflect.StructField, source StringsSource, key string) (bool, error) {
	if source == nil {
		return false, nil
	}
	ss, ok, err := source.GetStrings(key)
	if err != nil {
		return false, fmt.Errorf("source get values: %w", err)
	}
	if !ok {
		return false, nil
	}

	return true, setter.SetSliceStrings(value, field, ss)
}

// MapStringsSource - map[string][]string implement StringsSource.
type MapStringsSource map[string][]string

var _ StringsSource = MapStringsSource(nil)

// Get value from source.
func (s MapStringsSource) GetStrings(key string) ([]string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}
