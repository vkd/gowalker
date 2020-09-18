package gowalker

import (
	"reflect"
	"strings"
)

func FieldKey(tag Tag, namer AppendNamer, fs Fields) string {
	var key string
	for _, field := range fs {
		if tag != "" {
			k, ok := field.Tag.Lookup(string(tag))
			if ok {
				key = k
				continue
			}
		}
		key = namer.FieldKey(key, field)
	}
	return key
}

type AppendNamer struct {
	separator string
	convFn    func(string) string
}

func (a AppendNamer) FieldKey(parent string, field reflect.StructField) string {
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

// Fullname namer.
func Fullname(sep string, convFn func(string) string) AppendNamer {
	return AppendNamer{sep, convFn}
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
