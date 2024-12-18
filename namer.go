package gowalker

import (
	"iter"
	"strings"
)

type Namer interface {
	Key(fs iter.Seq[string]) string
}

type NamerFunc func(fs iter.Seq[string]) string

func (n NamerFunc) Key(fs iter.Seq[string]) string { return n(fs) }

// Fullname namer.
func Fullname(separator string, convFn func(string) string) Namer {
	return NamerFunc(func(fs iter.Seq[string]) string {
		var out string
		var sep string
		if convFn == nil {
			convFn = func(s string) string { return s }
		}
		for f := range fs {
			out += sep
			out += convFn(f)

			sep = separator
		}
		return out
	})
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
