package gowalker

import (
	"reflect"
)

func StringSetter(tag Tag, namer AppendNamer, source StringSource) Setter {
	return SetterFunc(func(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
		key := FieldKey(tag, namer, fs)
		return SetStringSource(value, field, source, key)
	})
}

func StringsSetter(tag Tag, namer AppendNamer, source StringsSource) Setter {
	return SetterFunc(func(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
		key := FieldKey(tag, namer, fs)
		return SetStringsSource(value, field, source, key)
	})
}

func Envs(tag Tag, namer AppendNamer, osLookupEnv LookupFuncSource) Setter {
	return StringSetter(tag, namer, osLookupEnv)
}
