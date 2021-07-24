package gowalker_test

import (
	"reflect"
	"testing"

	"github.com/vkd/gowalker"
	"github.com/vkd/gowalker/setter"
)

func BenchmarkWalk_NewWalker(b *testing.B) {
	var cfg struct {
		Name string
		DB   struct {
			Type string
			Port int
		}
	}

	env := map[string]string{
		"NAME":    "service",
		"DB_TYPE": "postgres",
		"DB_PORT": "9000",
	}

	fs := gowalker.MakeFields(2)

	w := StringSetter(
		gowalker.FieldKey("config", gowalker.UpperNamer),
		MapStringSource(env),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, fs, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

func BenchmarkWalk_MapSource(b *testing.B) {
	var cfg struct {
		Name string
		DB   struct {
			Type string
			Port int
		}
	}

	env := map[string]string{
		"NAME":    "service",
		"DB_TYPE": "postgres",
		"DB_PORT": "9000",
	}

	fs := gowalker.MakeFields(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, fs,
			StringSetter(
				gowalker.FieldKey("config", gowalker.UpperNamer),
				MapStringSource(env),
			),
		)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

func BenchmarkWalk_NewWalker_ConcatNamer(b *testing.B) {
	var cfg struct {
		Name string
		DB   struct {
			Type string
			Port int
		}
	}

	env := map[string]string{
		"Name":    "service",
		"DB.Type": "postgres",
		"DB.Port": "9000",
	}

	w := StringSetter(
		gowalker.FieldKey("config", gowalker.StructFieldNamer),
		MapStringSource(env),
	)

	fs := gowalker.MakeFields(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, fs, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

type MapStringSource map[string]string

func StringSetter(fk gowalker.FieldKeyer, m MapStringSource) gowalker.Walker {
	return gowalker.WalkerFunc(func(value reflect.Value, field reflect.StructField, fs gowalker.Fields) (stop bool, _ error) {
		key, ok := fk.FieldKey(field, fs)
		if !ok {
			return false, nil
		}
		v, ok := m[key]
		if !ok {
			return false, nil
		}
		return true, setter.SetString(value, field, v)
	})
}
