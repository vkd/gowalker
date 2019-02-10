package gowalker_test

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vkd/gowalker"
)

func ExampleWalk_upperConfigEnv() {
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
		gowalker.WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
			return gowalker.StringWalkerStep("config", gowalker.StringSourceMapString(env), value, field)
		}),
		func(fields []string) string {
			return strings.ToUpper(strings.Join(fields, "_"))
		},
	)

	err := gowalker.Walk(&cfg, w)
	fmt.Printf("cfg: %#v, %v", cfg, err)
	// Output: cfg: struct { Name string; DB struct { Type string; Port int } }{Name:"service", DB:struct { Type string; Port int }{Type:"postgres", Port:9000}}, <nil>
}

func ExampleWalk_ginBinding() {
	var q struct {
		Name    string   `uri:"name"`
		Age     int      `uri:"age,default=25"`
		Friends []string `uri:"friends"`
		Coins   []int    `uri:"coins,default=40"`
		Keys    []int
	}

	uri := map[string][]string{
		"name":    {"mike"},
		"friends": {"igor", "alisa"},
	}

	w := gowalker.WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
		return gowalker.SliceStringsWalkerStep("uri", gowalker.SliceStringsSourceMapStrings(uri), value, field)
	})

	err := gowalker.Walk(&q, w)
	fmt.Printf("uri: %#v, %v", q, err)
	// Output: uri: struct { Name string "uri:\"name\""; Age int "uri:\"age,default=25\""; Friends []string "uri:\"friends\""; Coins []int "uri:\"coins,default=40\""; Keys []int }{Name:"mike", Age:25, Friends:[]string{"igor", "alisa"}, Coins:[]int{40}, Keys:[]int(nil)}, <nil>
}
