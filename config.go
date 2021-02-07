package gowalker

import (
	"fmt"
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
func Config(cfg interface{}, setters ...Setter) error {
	for _, s := range setters {
		if i, ok := s.(Initer); ok {
			err := i.Init(cfg)
			if err != nil {
				return fmt.Errorf("init %T setter: %w", i, err)
			}
		}
	}

	fs := make(Fields, 0, 6)

	err := Walk(cfg, fs, MultiSetterOR(setters...))
	if err != nil {
		return fmt.Errorf("walk through config structure: %w", err)
	}

	return nil
}

type Initer interface {
	Init(cfg interface{}) error
}
