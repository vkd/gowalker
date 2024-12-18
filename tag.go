package gowalker

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strings"

	"github.com/vkd/gowalker/setter"
)

type Tag string

func (t Tag) Name() string {
	return string(t)
}

func (t Tag) Doc(field reflect.StructField, fs Fields) string {
	v, _ := t.get(field)
	return v
}

func (t Tag) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	v, ok := t.get(field)
	if !ok {
		return false, nil
	}

	return true, setter.SetString(value, field, v)
}

func (t Tag) get(field reflect.StructField) (string, bool) {
	return field.Tag.Lookup(string(t))
}

func (t Tag) Names(f Fields) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, field := range f {
			v, ok := t.get(field)
			if !ok {
				v = field.Name
			}
			for _, vv := range strings.Split(v, " ") {
				if !yield(vv) {
					return
				}
			}
		}
	}
}

var ErrRequiredField = errors.New("field is required")

type Required Tag

func (r Required) Name() string {
	return Tag(r).Name()
}

func (r Required) Doc(field reflect.StructField, fs Fields) string {
	_, ok := Tag(r).get(field)
	if ok {
		return "*"
	}
	return ""
}

func (r Required) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	t, ok := Tag(r).get(field)
	if !ok {
		return false, nil
	}

	switch t {
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	}

	return false, fmt.Errorf("%s: %w", field.Name, ErrRequiredField)
}
