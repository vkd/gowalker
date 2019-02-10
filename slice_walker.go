package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setters"
)

// SliceStringsSource - source of values by slice of strings
type SliceStringsSource interface {
	Get(key string) (value []string, ok bool, err error)
	StringSource() StringSource
}

// SliceStringGetValue - get value from field for slice of strings
func SliceStringGetValue(tag string, source SliceStringsSource, field reflect.StructField) ([]string, bool, error) {
	t := TagStringParse(field, tag)
	ss, ok, err := source.Get(t.Value)
	if err != nil {
		return nil, false, err
	}
	if !ok && !t.IsDefaultValue {
		return nil, false, nil
	}
	if !ok {
		ss = []string{t.DefaultValue}
	}
	return ss, true, nil
}

// SliceStringsWalkerStep - step of walker by slice of strings
func SliceStringsWalkerStep(tag string, source SliceStringsSource, value reflect.Value, field reflect.StructField) (bool, error) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		ss, ok, err := SliceStringGetValue(tag, source, field)
		if err != nil || !ok {
			return ok, err
		}
		return true, setters.SetValueBySliceOfString(value, field, ss)
	}
	return StringWalkerStep(tag, source.StringSource(), value, field)
}

// SliceStringsSourceMapStrings - map[string][]string implement SliceStringsSource
type SliceStringsSourceMapStrings map[string][]string

// Get value from source
func (s SliceStringsSourceMapStrings) Get(key string) ([]string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}

// StringSource - return string source
func (s SliceStringsSourceMapStrings) StringSource() StringSource {
	return StringSourceMapStringsByFirst(s)
}

// SliceStringsSourceFunc - func implement SliceStringSource
type SliceStringsSourceFunc func(key string) ([]string, bool, error)

var _ SliceStringsSource = SliceStringsSourceFunc(nil)

// Get value from source
func (fn SliceStringsSourceFunc) Get(key string) ([]string, bool, error) {
	return fn(key)
}

// StringSource - return string source
func (fn SliceStringsSourceFunc) StringSource() StringSource {
	return StringSourceFunc(func(key string) (string, bool, error) {
		ss, ok, err := fn(key)
		if err != nil || !ok {
			return "", ok, err
		}
		var s string
		if len(ss) > 0 {
			s = ss[0]
		}
		return s, true, nil
	})
}
