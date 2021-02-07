package gowalker

import (
	"flag"
	"fmt"
	"reflect"
	"sync"

	"github.com/vkd/gowalker/setter"
)

func Flags(tag Tag, namer AppendNamer, osArgs []string) Setter {
	return &flagSetter{tag: tag, namer: namer, osArgs: osArgs}
}

type flagSetter struct {
	mx            sync.Mutex
	initialized   bool
	updatedFields UpdatedFields

	tag    Tag
	namer  AppendNamer
	osArgs []string
}

func (f *flagSetter) Init(ptr interface{}) error {
	f.mx.Lock()
	defer f.mx.Unlock()
	if f.initialized {
		return nil
	}

	updated, err := FlagWalk(ptr, make(Fields, 6), f.tag, f.namer, f.osArgs)
	if err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	f.updatedFields = updated
	f.initialized = true
	return nil
}

func (f *flagSetter) TrySet(v reflect.Value, sf reflect.StructField, fs Fields) (ok bool, _ error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	if !f.initialized {
		return false, fmt.Errorf("setter is not initialized")
	}

	_, ok = f.updatedFields[FieldKey("", StructFieldNamer, fs)]
	return ok, nil
}

func FlagWalk(ptr interface{}, fs Fields, tag Tag, namer AppendNamer, osArgs []string) (UpdatedFields, error) {
	var name string
	if len(osArgs) > 0 {
		name = osArgs[0]
		osArgs = osArgs[1:]
	}

	updatedFields := make(UpdatedFields)

	fset := flag.NewFlagSet(name, flag.ContinueOnError)

	err := Walk(ptr, fs, flagWalker{
		tag:           tag,
		namer:         namer,
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
