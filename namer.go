package gowalker

import (
	"reflect"
	"strings"
)

type Namer interface {
	Key(fs Fields) string
}

type appendNamer struct {
	separator string
	convFn    func(string) string
}

func (a appendNamer) FieldKey(parent string, field reflect.StructField) string {
	name := field.Name
	if parent != "" {
		if a.convFn != nil {
			return parent + a.separator + a.convFn(name)
		}
		return parent + a.separator + name
	}
	if a.convFn != nil {
		return a.convFn(name)
	}
	return name
}

func (a appendNamer) Key(fs Fields) string {
	var key string
	for _, f := range fs {
		key = a.FieldKey(key, f)
	}
	return key
}

// Fullname namer.
func Fullname(sep string, convFn func(string) string) Namer {
	return appendNamer{sep, convFn}
}

// UpperNamer - concat a uppercase parent's name with a uppercase child's one with underscore.
var UpperNamer = Fullname("_", strings.ToUpper)

// EnvNamer - STRUCT_FIELD naming.
var EnvNamer = UpperNamer

// StructFieldNamer - concat a parent's name with a child's one with dot.
var StructFieldNamer = Fullname(".", nil)

// DashToLoverNamer - concat a lowercase parent's name with a lowercase child's one with dash.
var DashToLoverNamer = Fullname("-", strings.ToLower)

// FlagNamer - struct-field naming.
var FlagNamer = DashToLoverNamer

var DefaultNamer = Fullname("", nil)
