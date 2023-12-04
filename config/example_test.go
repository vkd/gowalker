package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vkd/gowalker"
)

func ExampleConfig() {
	var cfg struct {
		Name    string
		Timeout time.Duration `default:"3s"`

		DB struct {
			Username string `required:""`
			Password string `required:""`
		}

		Metrics struct {
			Addr  string `env:"METRICS_URL"`
			Label string `default:"gowalker"`
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
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=Gowalker"}

	err := Walk(&cfg, nil,
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.Fullname("-", strings.ToLower)), osArgs),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.Fullname("_", strings.ToUpper)), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker 5s {postgres example} {localhost:5678 gowalker}}, <nil>
}

func ExamplePrintHelp() {
	var cfg struct {
		Name    string
		Timeout time.Duration `default:"3s"`

		DB struct {
			Username string `required:""`
			Password string `required:""`
		}

		Metrics struct {
			Addr  string `env:"METRICS_URL"`
			Label string `default:"gowalker"`
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
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=Gowalker", "--help"}

	err := Walk(&cfg, log.New(os.Stdout, "", 0),
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.Fullname("-", strings.ToLower)), osArgs),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.Fullname("_", strings.ToUpper)), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	fmt.Printf("%v", err)
	// Output:
	// flag          | ENV           | default  | required
	// ------------- | ------------- | -------- | --------
	// name          | NAME          |          |
	// timeout       | TIMEOUT       | 3s       |
	// db            | DB            |          |
	// db-username   | DB_USERNAME   |          | *
	// db-password   | DB_PASSWORD   |          | *
	// metrics       | METRICS       |          |
	// metrics-addr  | METRICS_URL   |          |
	// metrics-label | METRICS_LABEL | gowalker |
	// parse flags: print help
}

func ExampleConfigInline() {
	var cfg struct {
		Name    string
		Timeout time.Duration `default:"3s"`

		Wrap struct {
			DB struct {
				Username string `required:""`
				Password string `required:""`
			}
		} `walker:"embed"`

		Metrics struct {
			Addr  string `env:"METRICS_URL"`
			Label string `default:"gowalker"`
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
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=Gowalker"}

	err := Walk(&cfg, nil,
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.Fullname("-", strings.ToLower)), osArgs),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.Fullname("_", strings.ToUpper)), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker 5s {{postgres example}} {localhost:5678 gowalker}}, <nil>
}
