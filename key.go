package gowalker

import "reflect"

type FieldKeyer interface {
	FieldKey(reflect.StructField, Fields) (string, bool)
}

func FieldKey(t Tag, namer Namer) FieldKeyer {
	return &fieldKey{Tag: t, Namer: namer}
}

type fieldKey struct {
	Tag
	Namer
}

func (f *fieldKey) FieldKey(field reflect.StructField, fs Fields) (string, bool) {
	key, ok := f.Tag.get(field)
	if ok {
		return key, true
	}
	return f.Namer.Key(fs), true
}
