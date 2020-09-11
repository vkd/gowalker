package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// SliceStringGetValue - get value from field for slice of strings
func SliceStringGetValue(tag string, source StringsSource, field reflect.StructField, name Name, namer Namer) ([]string, bool, error) {
	t, ok := field.Tag.Lookup(tag)
	if !ok {
		t = name.Get(namer)
	}
	ss, ok, err := source.GetStrings(t)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return ss, true, nil
}

// SliceStringsWalkerStep - step of walker by slice of strings
func SliceStringsWalkerStep(tag string, source SliceSourcer, value reflect.Value, field reflect.StructField, name Name, namer Namer) (bool, error) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		ss, ok, err := SliceStringGetValue(tag, source, field, name, namer)
		if err != nil || !ok {
			return ok, err
		}
		return true, setter.SetSliceStrings(value, field, ss)
	}
	return StringWalkerStep(tag, source, value, field, name, namer)
}
