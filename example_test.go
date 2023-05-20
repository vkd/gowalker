package gowalker_test

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vkd/gowalker"
	"github.com/vkd/gowalker/config"
	"github.com/vkd/gowalker/setter"
)

type titleName string

func (c *titleName) SetString(s string) error {
	*c = titleName(strings.Title(s) + "!")
	return nil
}

func ExampleConfig() {
	var cfg struct {
		Name    titleName
		Timeout time.Duration `default:"3s"`

		DB struct {
			Username string `required:""`
			Password string `required:""`
		}

		Metrics struct {
			Addr string `env:"METRICS_URL"`
		}
	}

	// osLookupEnv := os.LookupEnv
	osLookupEnv := func(key string) (string, bool) {
		v, ok := map[string]string{
			"MY_ORG_DB_USERNAME": "postgres",
			"MY_ORG_METRICS_URL": "localhost:5678",
		}[key]
		return v, ok
	}

	// osArgs := os.Args
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=gowalker"}

	err := config.Walk(&cfg, nil,
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.Fullname("-", strings.ToLower)), osArgs),
		gowalker.Envs(gowalker.Prefix("MY_ORG_", gowalker.FieldKey("env", gowalker.Fullname("_", strings.ToUpper))), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker! 5s {postgres example} {localhost:5678}}, <nil>
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

	fk := gowalker.FieldKey("uri", gowalker.DefaultNamer)
	w := gowalker.WalkerFunc(func(value reflect.Value, field reflect.StructField, fs gowalker.Fields) (stop bool, _ error) {
		key, ok := fk.FieldKey(field, fs)
		if !ok {
			return false, nil
		}
		v, ok := uri[key]
		if !ok {
			return false, nil
		}
		return true, setter.SetSliceStrings(value, field, v)
	})
	err := gowalker.Walk(&q, gowalker.MakeFields(1), w)
	fmt.Printf("%+v, %v", q, err)
	// Output: {Name:mike Age:0 Friends:[igor alisa] Coins:[] Keys:[]}, <nil>
}

type visitedFields []string

func (f *visitedFields) Step(value reflect.Value, field reflect.StructField, fs gowalker.Fields) (set bool, err error) {
	key := gowalker.DashToLoverNamer.Key(fs)
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
	err := gowalker.Walk(&config, gowalker.MakeFields(2), &fs)
	fmt.Printf("fields: %v, %v", fs, err)
	// Output: fields: [name port db db-url db-port], <nil>
}
