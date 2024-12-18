package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/vkd/gowalker"
)

func Default(cfg interface{}) error {
	return defaultConfig(cfg, os.Args, os.LookupEnv)
}

func defaultConfig(cfg interface{}, osArgs []string, osLookupEnv func(string) (string, bool)) error {
	return Walk(cfg, log.New(os.Stdout, "", 0),
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.FlagNamer), osArgs),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.EnvNamer), osLookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
}

type Configer interface {
	Name() string
	Doc(reflect.StructField, gowalker.Fields) string
	gowalker.Walker
}

// Walk - fills config structure.
func Walk(cfg interface{}, l Logger, cs ...Configer) error {
	for _, s := range cs {
		if i, ok := s.(Initer); ok {
			err := i.Init(cfg)
			if err != nil {
				if errors.Is(err, gowalker.ErrPrintHelp) {
					PrintHelp(cfg, l, cs...)
					return err
				}
				return fmt.Errorf("init %T setter: %w", i, err)
			}
		}
	}

	var svs []gowalker.Walker
	for _, c := range cs {
		svs = append(svs, c)
	}

	err := gowalker.Walk(
		cfg,
		gowalker.MakeFields(4),
		gowalker.WalkersOR(svs),
	)
	if err != nil {
		return fmt.Errorf("walk through config structure: %w", err)
	}

	return nil
}

type Initer interface {
	Init(cfg interface{}) error
}
