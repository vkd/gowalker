package gowalker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/vkd/gowalker/setter"
)

// Tag of a struct field.
type Tag string

var _ Setter = Tag("")

// Set of walker implementation.
func (t Tag) TrySet(value reflect.Value, field reflect.StructField, _ Fields) (bool, error) {
	v, ok := field.Tag.Lookup(string(t))
	if !ok {
		return false, nil
	}
	return true, setter.SetString(value, field, v)
}

func Required(tag Tag, updatedFields UpdatedFields) Setter {
	return required{tag: tag, updatedFields: updatedFields}
}

var ErrRequiredField = errors.New("field is required")

type required struct {
	tag           Tag
	updatedFields UpdatedFields
}

var _ Setter = required{}

func (r required) TrySet(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	t, ok := field.Tag.Lookup(string(r.tag))
	if !ok {
		return false, nil
	}

	switch t {
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	}

	if r.updatedFields != nil {
		key := FieldKey("", StructFieldNamer, fs)
		_, ok = r.updatedFields[key]
		if ok {
			return false, nil
		}
	}

	return false, fmt.Errorf("%s: %w", field.Name, ErrRequiredField)
}
