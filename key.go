package gowalker

import "reflect"

type FieldKeyer interface {
	FieldKey(reflect.StructField, Fields) string
}

type FieldKeyerFunc func(reflect.StructField, Fields) string

func (f FieldKeyerFunc) FieldKey(field reflect.StructField, fs Fields) string { return f(field, fs) }

func NestedFieldKey(t, fkey Tag, namer Namer) FieldKeyer {
	return FieldKeyerFunc(func(field reflect.StructField, fs Fields) string {
		key, ok := t.get(field)
		if ok {
			return key
		}
		return namer.Key(fkey.Names(fs))

	})
}

func FieldKey(t Tag, namer Namer) FieldKeyer {
	return FieldKeyerFunc(func(field reflect.StructField, fs Fields) string {
		key, ok := t.get(field)
		if ok {
			return key
		}
		return namer.Key(fs.Names())

	})
}

func Prefix(p string, fk FieldKeyer) FieldKeyer {
	return prefix{p, fk}
}

type prefix struct {
	p  string
	fk FieldKeyer
}

func (p prefix) FieldKey(field reflect.StructField, fs Fields) string {
	return p.p + p.fk.FieldKey(field, fs)
}
