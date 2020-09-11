package gowalker

import (
	"flag"
	"reflect"

	"github.com/vkd/gowalker/setter"
)

func DefaultFlagWalker() Walker {
	return Flag(FieldKeys(Tag("flag"), FlagNamer))
}

func Flag(k FieldKeyer) Walker {
	return flagWalker{k: k}
}

type flagWalker struct {
	k FieldKeyer
}

func (f flagWalker) Step(value reflect.Value, field reflect.StructField, name Name) (isSet bool, _ error) {
	if value.Kind() == reflect.Struct {
		return false, nil
	}

	n, ok := f.k.GetFieldKey(field, name)
	if !ok {
		return false, nil
	}

	flag.Var(fieldValue{Value: value, StructField: field}, n, name.Get(ConcatNamer))

	return true, nil
}

type fieldValue struct {
	reflect.Value
	reflect.StructField
}

var _ flag.Value = fieldValue{}

func (f fieldValue) String() string {
	return "StructField"
}

func (f fieldValue) Set(s string) error {
	return setter.SetString(f.Value, f.StructField, s)
}
