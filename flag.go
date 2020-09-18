package gowalker

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/vkd/gowalker/setter"
)

func FlagWalk(ptr interface{}, fs Fields, osArgs []string) (UpdatedFields, error) {
	var name string
	if len(osArgs) > 0 {
		name = osArgs[0]
		osArgs = osArgs[1:]
	}

	updatedFields := make(UpdatedFields)

	fset := flag.NewFlagSet(name, flag.ContinueOnError)

	err := Walk(ptr, fs, flagWalker{
		tag:           Tag("flag"),
		namer:         FlagNamer,
		fset:          fset,
		updatedFields: updatedFields,
	})
	if err != nil {
		return updatedFields, fmt.Errorf("setup flags values: %w", err)
	}

	err = fset.Parse(osArgs)
	if err != nil {
		return updatedFields, fmt.Errorf("parse flags: %w", err)
	}

	return updatedFields, nil
}

type flagWalker struct {
	tag           Tag
	namer         AppendNamer
	fset          *flag.FlagSet
	updatedFields UpdatedFields
}

func (f flagWalker) TrySet(value reflect.Value, field reflect.StructField, fs Fields) (isSet bool, _ error) {
	if f.fset == nil {
		return false, nil
	}
	if value.Kind() == reflect.Struct {
		return false, nil
	}

	key := FieldKey(f.tag, f.namer, fs)
	fieldPath := FieldKey("", StructFieldNamer, fs)

	f.fset.Var(fieldValue{Value: value, StructField: field, fieldPath: fieldPath, updatedFields: f.updatedFields}, key, fieldPath)

	return true, nil
}

type fieldValue struct {
	reflect.Value
	reflect.StructField

	fieldPath     string
	updatedFields UpdatedFields
}

var _ flag.Value = fieldValue{}

func (f fieldValue) String() string {
	return "StructField"
}

func (f fieldValue) Set(s string) error {
	f.updatedFields[f.fieldPath] = struct{}{}
	return setter.SetString(f.Value, f.StructField, s)
}

type UpdatedFields map[string]struct{}
