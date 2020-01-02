package gowalker

import (
	"reflect"
	"strings"
)

// Struct - walk by struct fields.
func Struct(value interface{}, tag string, source StringSource) error {
	return StructFullname(value, tag, source, nil)
}

// StructFullname - walk by struct fields with custom field name generator.
func StructFullname(value interface{}, tag string, source StringSource, namer Namer) error {
	w := NewStringWalker(tag, source)
	_, err := walkIface(value, w, namer)
	return err
}

// Walk - walk struct by all public fields
func Walk(value interface{}, walker Walker) error {
	return WalkFullname(value, walker, nil)
}

// WalkFullname - walk struct by all public fields with custom field name generator.
func WalkFullname(value interface{}, walker Walker, namer Namer) error {
	_, err := walkIface(value, walker, namer)
	return err
}

// Wrapper - interface to allow build tree of struct fields
type Wrapper interface {
	Walker
	Wrap(field reflect.StructField) Walker
}

// WalkerFunc - func implemented Walk interface
type WalkerFunc func(value reflect.Value, field reflect.StructField) (bool, error)

// Step - one step of walker
func (f WalkerFunc) Step(value reflect.Value, field reflect.StructField) (bool, error) {
	return f(value, field)
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

// UpperNamer - concat a uppercase parent's name with a uppercase child's one with underscore.
var UpperNamer Namer = appendNamer{"_", strings.ToUpper}

// ConcatNamer - concat a parent's name with a child's one with underscore.
var ConcatNamer Namer = appendNamer{separator: "_"}

// DashToLoverNamer - concat a lowercase parent's name with a lowercase child's one with dash.
var DashToLoverNamer Namer = appendNamer{"-", strings.ToLower}
