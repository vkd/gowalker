package gowalker_test

import (
	"strings"
	"testing"

	"github.com/vkd/gowalker"
)

func BenchmarkWalk(b *testing.B) {
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

	w := gowalker.NewStringWalker("config", gowalker.StringSourceMapString(env))

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
		w := gowalker.NewStringWalker("config", gowalker.StringSourceMapString(env))

		err := gowalker.Walk(&cfg, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

func BenchmarkWalk_Wrap(b *testing.B) {
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

	w := gowalker.NewWrapFieldNameWalkerConv(
		gowalker.NewStringWalker("config", gowalker.StringSourceMapString(env)),
		func(fields []string) string {
			return strings.ToUpper(strings.Join(fields, "_"))
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gowalker.Walk(&cfg, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}

func BenchmarkWalk_MapSource_Wrap(b *testing.B) {
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
		w := gowalker.NewWrapFieldNameWalkerConv(
			gowalker.NewStringWalker("config", gowalker.StringSourceMapString(env)),
			func(fields []string) string {
				return strings.ToUpper(strings.Join(fields, "_"))
			},
		)

		err := gowalker.Walk(&cfg, w)
		if err != nil {
			b.Fatalf("Error on walk: %v", err)
		}
	}
}
