package gowalker

import (
	"strings"
)

type Namer interface {
	Key(fs Fields) string
}

type appendNamer struct {
	separator string
	convFn    func(string) string
}

func (a appendNamer) Key(fs Fields) string {
	var out string
	var sep string
	convFn := a.convFn
	if convFn == nil {
		convFn = func(s string) string { return s }
	}
	for _, f := range fs {
		out += sep
		out += convFn(f.Name)

		sep = a.separator
	}
	return out
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
