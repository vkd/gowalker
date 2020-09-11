package gowalker_test

import (
	"fmt"
	"reflect"

	"github.com/vkd/gowalker"
)

func ExampleWalk_envConfig() {
	var cfg struct {
		DB struct {
			Name    string
			Address string `env:"DB_URL"`
			Port    int    `default:"5432" env:"DB_PORT"`
		}
	}

	// osLookupEnv := os.LookupEnv
	osLookupEnv := func(key string) (string, bool) {
		v, ok := map[string]string{
			"DB_NAME": "Env",
			"DB_URL":  "postgres",
		}[key]
		return v, ok
	}

	err := gowalker.Walk(&cfg,
		gowalker.Tag("default"),

		gowalker.StringSetter(
			gowalker.FieldKeys(
				gowalker.Tag("env"),
				gowalker.EnvNamer,
			),
			gowalker.LookupFuncSource(osLookupEnv),
		),
	)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {{Env postgres 5432}}, <nil>
}

func ExampleWalk_fromMapOfStrings() {
	var q struct {
		Name    string   `uri:"name"`
		Age     int      `uri:"age"`
		Friends []string `uri:"friends"`
		Coins   []int    `uri:"coins"`
		Keys    []int
	}

	uri := map[string][]string{
		"name":    {"mike"},
		"friends": {"igor", "alisa"},
	}

	w := gowalker.StringsSetter(
		gowalker.Tag("uri"),
		gowalker.MapStringsSource(uri),
	)
	err := gowalker.Walk(&q, w)
	fmt.Printf("%+v, %v", q, err)
	// Output: {Name:mike Age:0 Friends:[igor alisa] Coins:[] Keys:[]}, <nil>
}

type visitedFields []string

func (f *visitedFields) Step(value reflect.Value, field reflect.StructField, name gowalker.Name) (set bool, err error) {
	key := name.Get(gowalker.DashToLoverNamer)
	*f = append(*f, key)
	return false, nil
}

func ExampleWalk_collectAllPublicFields() {
	var config struct {
		Name string
		Port int
		DB   struct {
			URL  string
			Port int
		}
	}

	var fs visitedFields
	err := gowalker.Walk(&config, &fs)
	fmt.Printf("fields: %v, %v", fs, err)
	// Output: fields: [name port db db-url db-port], <nil>
}
