package gowalker

import (
	"flag"
	"fmt"
	"reflect"
	"sync"

	"github.com/vkd/gowalker/setter"
)

func Flags(fk FieldKeyer, osArgs []string) *Flag {
	return &Flag{FieldKeyer: fk, OsArgs: osArgs}
}

type Flag struct {
	mx            sync.Mutex
	updatedFields UpdatedFields

	FieldKeyer
	OsArgs []string
}

func (*Flag) Name() string {
	return "flag"
}

func (f *Flag) Init(ptr interface{}) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	updated, err := FlagWalk(ptr, make(Fields, 0, 6), f.FieldKeyer, f.OsArgs)
	if err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}
	f.updatedFields = updated

	return nil
}

func (f *Flag) Step(v reflect.Value, sf reflect.StructField, fs Fields) (ok bool, _ error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	return f.updatedFields.Has(fs), nil
}

func (f *Flag) Doc(field reflect.StructField, fs Fields) string {
	v := f.FieldKeyer.FieldKey(field, fs)
	return v
}

var ErrPrintHelp = fmt.Errorf("print help")

func FlagWalk(ptr interface{}, fs Fields, kb FieldKeyer, osArgs []string) (UpdatedFields, error) {
	var name string
	if len(osArgs) > 0 {
		name = osArgs[0]
		osArgs = osArgs[1:]
	}

	updatedFields := make(UpdatedFields)

	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	hLong := fset.Bool("help", false, "print help")
	hShort := fset.Bool("h", false, "print help")

	err := Walk(ptr, fs, flagWalker{
		FieldKeyer:    kb,
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

	if *hLong || *hShort {
		return nil, ErrPrintHelp
	}

	return updatedFields, nil
}

type flagWalker struct {
	FieldKeyer
	fset          *flag.FlagSet
	updatedFields UpdatedFields
}

func (f flagWalker) Step(value reflect.Value, field reflect.StructField, fs Fields) (isSet bool, _ error) {
	if f.fset == nil {
		return false, nil
	}
	if value.Kind() == reflect.Struct {
		return false, nil
	}

	key := f.FieldKey(field, fs)
	fieldPath := KeyUpdatedFields(fs)

	f.fset.Var(fieldValue{Value: value, StructField: field, fieldPath: fieldPath, updatedFields: f.updatedFields}, key, fieldPath)

	return true, nil
}

type fieldValue struct {
	reflect.Value
	reflect.StructField

	fieldPath     string
	updatedFields UpdatedFields
}

func (f fieldValue) Set(s string) error {
	f.updatedFields[f.fieldPath] = struct{}{}
	return setter.SetString(f.Value, f.StructField, s)
}

type UpdatedFields map[string]struct{}

func KeyUpdatedFields(fs Fields) string {
	return StructFieldNamer.Key(fs.Names())
}

func (u UpdatedFields) Has(fs Fields) bool {
	_, ok := u[KeyUpdatedFields(fs)]
	return ok
}
