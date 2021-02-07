package config

import (
	"os"

	"github.com/vkd/gowalker"
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
func Fill(ptr interface{}) error {
	return gowalker.Config(ptr, DefaultSetters()...)
}

func DefaultSetters() []gowalker.Setter {
	return []gowalker.Setter{
		gowalker.Flags("flag", gowalker.FlagNamer, os.Args),
		gowalker.Envs("env", gowalker.EnvNamer, os.LookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	}
}
