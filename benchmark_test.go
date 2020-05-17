package gowalker_test

import (
	"reflect"
	"testing"

	"github.com/vkd/gowalker"
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

	w := gowalker.NewStringWalker("config", gowalker.MapStringSource(env), gowalker.UpperNamer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, w)
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, mapSourceConfigWalker(env))
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
		"NAME":    "service",
		"DB_TYPE": "postgres",
		"DB_PORT": "9000",
	}

	w := gowalker.NewStringWalker("config", gowalker.MapStringSource(env), gowalker.ConcatNamer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

type mapSourceConfigWalker gowalker.MapStringSource

var _ gowalker.Walker = (mapSourceConfigWalker)(nil)

func (m mapSourceConfigWalker) Step(value reflect.Value, field reflect.StructField, name gowalker.Name) (bool, error) {
	return gowalker.StringWalkerStep("config", gowalker.MapStringSource(m), value, field, name, gowalker.UpperNamer)
}
