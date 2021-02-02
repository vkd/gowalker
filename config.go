package gowalker

import (
	"fmt"
	"reflect"
)

// Config fills config structure.
// Useful if use `type Config struct{}` as a configuration for
// an application.
//
// Usage:
// var cfg struct {
// 	Port int `default:"8080"`
// 	DB   struct {
// 		Port     int    `default:"5432"`
// 		Username string `env:"DB_USERNAME" flag:"db-username"`
// 		Timeout time.Duration `default:"1s" flag:"postgres-timeout"`
// 	}
// }
// err := gowalker.Config(&cfg, os.LookupEnv, os.Args)
//
// DB_USERNAME=postgres ./app --port=8000 --db-port=5432 --postgres-timeout 3s.
func Config(cfg interface{}, osLookupEnv LookupFuncSource, osArgs []string) error {
	fs := make(Fields, 0, 6)

	updatedFields, err := FlagWalk(cfg, fs, osArgs)
	if err != nil {
		return fmt.Errorf("flag walk: %w", err)
	}

	setters := []Setter{
		StringSetter("env", EnvNamer, osLookupEnv),
		Tag("default"),
		Required("required", updatedFields),
	}

	err = Walk(cfg, fs, SetterFunc(func(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
		_, ok := updatedFields[FieldKey("", StructFieldNamer, fs)]
		if ok {
			return false, nil
		}
		for _, s := range setters {
			ok, err := s.TrySet(value, field, fs)
			if err != nil {
				return ok, fmt.Errorf("setter %T: %w", s, err)
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}))
	if err != nil {
		return fmt.Errorf("walk through config structure: %w", err)
	}

	return nil
}
