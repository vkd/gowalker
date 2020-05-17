package gowalker_test

import (
	"fmt"
	"reflect"

	"github.com/vkd/gowalker"
)

func ExampleWalk_ConfigFromMap() {
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

	w := gowalker.NewStringWalker("config", gowalker.MapStringSource(env))
	err := gowalker.WalkFullname(&cfg, w, gowalker.UpperNamer)

	fmt.Printf("cfg: %#v, %v", cfg, err)
	// Output: cfg: struct { Name string; DB struct { Type string; Port int } }{Name:"service", DB:struct { Type string; Port int }{Type:"postgres", Port:9000}}, <nil>
}

func ExampleWalk_SliceBindingWithDefault() {
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

	w := gowalker.NewStringWalker("uri", gowalker.MapStringsSourcer(uri))
	err := gowalker.Walk(&q, w)
	fmt.Printf("uri: %#v, %v", q, err)
	// Output: uri: struct { Name string "uri:\"name\""; Age int "uri:\"age\""; Friends []string "uri:\"friends\""; Coins []int "uri:\"coins\""; Keys []int }{Name:"mike", Age:0, Friends:[]string{"igor", "alisa"}, Coins:[]int(nil), Keys:[]int(nil)}, <nil>
}

func ExampleWalk_WalkWithMapSource() {
	var cfg struct {
		Name string
		DB   struct {
			Type       string
			PortNumber int `config:"PORT"`
			Username   string
		}
	}

	m := map[string]string{
		"NAME":    "service",
		"DB_TYPE": "postgres",
		"PORT":    "9000",
	}

	w := gowalker.NewStringWalker("config", gowalker.MapStringSource(m))
	err := gowalker.WalkFullname(&cfg, w, gowalker.UpperNamer)
	fmt.Printf("cfg: %v, %v", cfg, err)
	// Output: cfg: {service {postgres 9000 }}, <nil>
}

// var osLookupEnv = os.LookupEnv
var osLookupEnv = func(key string) (string, bool) {
	v, ok := map[string]string{
		"NAME":    "Env",
		"DB_URL":  "postgres",
		"DB_PORT": "5432",
	}[key]
	return v, ok
}

func ExampleWalk_ServiceEnvLoader() {
	type config struct {
		ServiceName string `env:"NAME"`
		Port        int    `env:"PORT"`
		DB          struct {
			Address string `env:"DB_URL"`
			Port    int    // DB_PORT
		}
	}

	var c config
	w := gowalker.NewStringWalker(
		"env",
		gowalker.LookupFuncSource(osLookupEnv),
	)
	err := gowalker.WalkFullname(&c, w, gowalker.UpperNamer)
	fmt.Printf("env: %v, %v", c, err)
	// Output: env: {Env 0 {postgres 5432}}, <nil>
}

type visitedFields []string

func (f *visitedFields) Step(value reflect.Value, field reflect.StructField) (set bool, err error) {
	*f = append(*f, field.Name)
	return false, nil
}

func ExampleWalk_CollectAllPublicFields() {
	var config struct {
		Name string
		Port int
		DB   struct {
			URL  string
			Port int
		}
	}

	var fs visitedFields
	err := gowalker.WalkFullname(&config, &fs, gowalker.DashToLoverNamer)
	fmt.Printf("fields: %v, %v", fs, err)
	// Output: fields: [name port db db-url db-port], <nil>
}
