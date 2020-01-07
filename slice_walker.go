package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// SliceStringGetValue - get value from field for slice of strings
func SliceStringGetValue(tag string, source SliceSourcer, field reflect.StructField) ([]string, bool, error) {
	t := TagStringParse(field, tag)
	ss, ok, err := source.GetStrings(t.Value)
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
func SliceStringsWalkerStep(tag string, source SliceSourcer, value reflect.Value, field reflect.StructField) (bool, error) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		ss, ok, err := SliceStringGetValue(tag, source, field)
		if err != nil || !ok {
			return ok, err
		}
		return true, setter.SetSliceStrings(value, field, ss)
	}
	return StringWalkerStep(tag, source, value, field)
}
