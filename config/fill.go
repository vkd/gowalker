package config

import (
	"os"

	"github.com/vkd/gowalker"
)

// Config fills config structure.
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
