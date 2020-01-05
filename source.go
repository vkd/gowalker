package gowalker

import "os"

// Sourcer - source of values by key.
type Sourcer interface {
	Get(key string) (value string, ok bool, err error)
}

// MapStringSource - map[string]string implement Sourcer.
type MapStringSource map[string]string

var _ Sourcer = MapStringSource(nil)

// Get value from source.
func (s MapStringSource) Get(key string) (string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}

// LookupFuncSource - func(key string) (string, bool) sourcer.
type LookupFuncSource func(key string) (string, bool)

var _ Sourcer = LookupFuncSource(nil)

// Get value from source.
func (f LookupFuncSource) Get(key string) (string, bool, error) {
	v, ok := f(key)
	return v, ok, nil
}

type EnvFuncSource = LookupFuncSource

var _ Sourcer = EnvFuncSource(nil)
var _ Sourcer = EnvFuncSource(os.LookupEnv)

// SliceSourcer - source of values by slice of strings.
type SliceSourcer interface {
	Sourcer
	GetStrings(key string) (value []string, ok bool, err error)
}

// MapStringsSourcer - map[string][]string implement SliceSourcer.
type MapStringsSourcer map[string][]string

var _ SliceSourcer = MapStringsSourcer(nil)

// Get value from source.
func (s MapStringsSourcer) GetStrings(key string) ([]string, bool, error) {
	v, ok := s[key]
	return v, ok, nil
}

// Get value from source.
func (s MapStringsSourcer) Get(key string) (string, bool, error) {
	return sliceStringsToGetString(s, key)
}
