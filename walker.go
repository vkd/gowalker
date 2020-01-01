package gowalker

import (
	"reflect"
)

// Walk - walk struct by all public fields
func Walk(value interface{}, walker Walker) error {
	_, err := walkIface(value, walker)
	return err
}

// Walker - interface to walk through struct fields
type Walker interface {
	Step(value reflect.Value, field reflect.StructField) (bool, error)
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
