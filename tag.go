package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// Tag of a struct field.
type Tag string

// Step of walker implementation.
func (t Tag) Step(value reflect.Value, field reflect.StructField) (bool, error) {
	v, ok := field.Tag.Lookup(string(t))
	if !ok {
		return false, nil
	}
	return true, setter.SetString(value, field, v)
}

// Walk - implementation of the config.Walker interface.
func (t Tag) Walk(v interface{}) error {
	return Walk(v, t)
}
