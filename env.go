package gowalker

import (
	"reflect"

	"github.com/vkd/gowalker/setter"
)

type LookupFunc func(key string) (string, bool)

func Envs(fk FieldKeyer, osLookupEnv LookupFunc) *Env {
	return &Env{FieldKeyer: fk, LookupFunc: osLookupEnv}
}

type Env struct {
	FieldKeyer
	LookupFunc
}

func (e *Env) Name() string {
	return "ENV"
}

func (e *Env) Doc(field reflect.StructField, fs Fields) string {
	return e.FieldKeyer.FieldKey(field, fs)
}

func (e *Env) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	key := e.FieldKeyer.FieldKey(field, fs)

	v, ok := e.LookupFunc(key)
	if !ok {
		return false, nil
	}

	return true, setter.SetString(value, field, v)
}
