package gowalker_test

import (
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

	fs := gowalker.MakeFields(2)

	w := gowalker.StringSetter(
		"config",
		gowalker.UpperNamer,
		gowalker.MapStringSource(env),
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
			gowalker.StringSetter(
				"config",
				gowalker.UpperNamer,
				gowalker.MapStringSource(env),
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

	w := gowalker.StringSetter(
		"config",
		gowalker.StructFieldNamer,
		gowalker.MapStringSource(env),
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
