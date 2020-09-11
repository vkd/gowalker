package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// Tag of a struct field.
type Tag string

var _ Walker = Tag("")

// Step of walker implementation.
func (t Tag) Step(value reflect.Value, field reflect.StructField, _ Name) (bool, error) {
	v, ok := field.Tag.Lookup(string(t))
	if !ok {
		return false, nil
	}
	return true, setter.SetString(value, field, v)
}

func (t Tag) GetFieldKey(field reflect.StructField, _ Name) (string, bool) {
	return field.Tag.Lookup(string(t))
}
