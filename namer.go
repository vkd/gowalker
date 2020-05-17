package gowalker

import (
	"strings"
)

// Namer - interface to get a full field name.
type Namer interface {
	FieldName(parent, name string) string
}

type namerFunc func(parent, name string) string

var _ Namer = namerFunc(func(_, _ string) string { return "" })

func (f namerFunc) FieldName(parent, name string) string { return f(parent, name) }

type appendNamer struct {
	separator string
	convFn    func(string) string
}

var _ Namer = appendNamer{}

func (a appendNamer) FieldName(parent, name string) string {
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
func Fullname(sep string, convFn func(string) string) Namer {
	return appendNamer{sep, convFn}
}

// UpperNamer - concat a uppercase parent's name with a uppercase child's one with underscore.
var UpperNamer Namer = Fullname("_", strings.ToUpper)

// EnvNamer - STRUCT_FIELD naming.
var EnvNamer = UpperNamer

// ConcatNamer - concat a parent's name with a child's one with underscore.
var ConcatNamer Namer = Fullname("_", nil)

// DashToLoverNamer - concat a lowercase parent's name with a lowercase child's one with dash.
var DashToLoverNamer Namer = Fullname("-", strings.ToLower)

// FlagNamer - struct-field naming.
var FlagNamer = DashToLoverNamer
