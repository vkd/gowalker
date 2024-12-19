package config

import (
	"fmt"
	"time"
)

func Example() {
	var cfg struct {
		Name    string
		Timeout time.Duration `default:"3s"`

		DB struct {
			Username string `required:""`
			Password string `required:""`
		}

		Metrics struct {
			Addr  string `fkey:"URL"`
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

	err := defaultConfig(&cfg, osArgs, osLookupEnv)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker 5s {postgres example} {localhost:5678 gowalker}}, <nil>
}

func ExamplePrintHelp() {
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
	osArgs := []string{"gowalker", "--timeout=5s", "--db-password", "example", "--name=Gowalker", "--help"}

	err := defaultConfig(&cfg, osArgs, osLookupEnv)
	fmt.Printf("%v", err)
	// Output:
	// flag          | ENV           | default  | required
	// ------------- | ------------- | -------- | --------
	// name          | NAME          |          |
	// timeout       | TIMEOUT       | 3s       |
	// db-username   | DB_USERNAME   |          | *
	// db-password   | DB_PASSWORD   |          | *
	// metrics-addr  | METRICS_URL   |          |
	// metrics-label | METRICS_LABEL | gowalker |
	// parse flags: print help
}

func ExampleMaybe() {
	var cfg struct {
		Name    maybe[string]
		Timeout maybe[time.Duration] `default:"3s"`

		Metrics struct {
			Addr  maybe[string] `fkey:"URL"`
			Label maybe[string] `default:"gowalker"`
			Tag   maybe[string]
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
	osArgs := []string{"gowalker", "--timeout=5s", "--name=Gowalker"}

	err := defaultConfig(&cfg, osArgs, osLookupEnv)
	fmt.Printf("%v, %v, %v, %v, %v", cfg.Metrics.Addr.IsSet, cfg.Metrics.Label.IsSet, cfg.Metrics.Tag.IsSet, cfg.Timeout, err)
	// Output: true, true, false, {true 5s}, <nil>
}

type maybe[T any] struct {
	IsSet bool
	Value T
}

func (m maybe[T]) GoWalkerMaybe() {}

func Example_embed() {
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
			Addr  string `fkey:"URL"`
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

	err := defaultConfig(&cfg, osArgs, osLookupEnv)
	fmt.Printf("%v, %v", cfg, err)
	// Output: {Gowalker 5s {{postgres example}} {localhost:5678 gowalker}}, <nil>
}
