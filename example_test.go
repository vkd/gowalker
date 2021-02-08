package gowalker_test

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vkd/gowalker"
)

type titleName string

func (c *titleName) SetString(s string) error {
	*c = titleName(strings.Title(s))
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
			"DB_USERNAME": "postgres",
			"METRICS_URL": "localhost:5678",
		}[key]
		return v, ok
	}

	// osArgs := os.Args
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=gowalker"}

	err := gowalker.Config(&cfg,
		gowalker.Flags("flag", gowalker.Fullname("-", strings.ToLower), osArgs),
		gowalker.Envs("env", gowalker.Fullname("_", strings.ToUpper), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker 5s {postgres example} {localhost:5678}}, <nil>
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
		gowalker.DefaultNamer,
		gowalker.MapStringsSource(uri),
	)
	err := gowalker.Walk(&q, gowalker.MakeFields(1), w)
	fmt.Printf("%+v, %v", q, err)
	// Output: {Name:mike Age:0 Friends:[igor alisa] Coins:[] Keys:[]}, <nil>
}

type visitedFields []string

func (f *visitedFields) TrySet(value reflect.Value, field reflect.StructField, fs gowalker.Fields) (set bool, err error) {
	key := gowalker.FieldKey("", gowalker.DashToLoverNamer, fs)
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
