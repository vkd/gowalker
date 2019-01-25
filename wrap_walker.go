package gowalker

import (
	"reflect"
	"strings"
)

type wrapFieldNameWalker struct {
	w      Walker
	convFn func([]string) string

	fields []string
}

// NewWrapFieldNameWalker - wrapper with collecting tree of calls
// FieldName_FieldName...
func NewWrapFieldNameWalker(w Walker) Walker {
	return NewWrapFieldNameWalkerConv(w, wrapJoinUnderscoreFunc)
}

// NewWrapFieldNameWalkerConv - wrapper with custom fieldname join function
func NewWrapFieldNameWalkerConv(w Walker, fn func([]string) string) Walker {
	return wrapFieldNameWalker{w: w, convFn: fn}
}

func (w wrapFieldNameWalker) Wrap(field reflect.StructField) Walker {
	w.fields = append(w.fields, field.Name)
	return w
}

func (w wrapFieldNameWalker) Step(value reflect.Value, field reflect.StructField) (bool, error) {
	field.Name = w.convFn(append(w.fields, field.Name))
	return w.w.Step(value, field)
}

func wrapJoinUnderscoreFunc(fields []string) string {
	return strings.Join(fields, "_")
}
