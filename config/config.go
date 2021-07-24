package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/vkd/gowalker"
)

func Default(cfg interface{}) error {
	return Walk(&cfg, log.New(os.Stdout, "", 0),
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.FlagNamer), os.Args),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.EnvNamer), os.LookupEnv),
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
				if errors.Is(err, gowalker.ErrHelp) {
					PrintHelp(cfg, l, cs...)
					return nil
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

type Logger interface {
	Print(...interface{})
}

func PrintHelp(cfg interface{}, log Logger, cs ...Configer) {
	if log == nil {
		return
	}

	var docs [][]string
	max := make([]int, len(cs))

	header := make([]string, 0, len(cs))
	for i, c := range cs {
		title := c.Name()
		max[i] = len(title)
		header = append(header, title)
	}
	docs = append(docs, header)

	err := gowalker.Walk(cfg, gowalker.MakeFields(4), gowalker.WalkerFunc(func(v reflect.Value, sf reflect.StructField, f gowalker.Fields) (stop bool, _ error) {

		line := make([]string, len(cs))
		for i, c := range cs {
			d := c.Doc(sf, f)

			line[i] = d
			if len(d) > max[i] {
				max[i] = len(d)
			}
		}
		docs = append(docs, line)

		return false, nil
	}))
	if err != nil {
		log.Print(fmt.Sprintf("Error on print help: %v", err))
		return
	}

	for line, doc := range docs {
		for i, d := range doc {
			dx := max[i] - len(d)
			if dx > 0 {
				doc[i] += strings.Repeat(" ", dx)
			}
		}
		log.Print(strings.TrimSpace(strings.Join(doc, " | ")))

		if line == 0 {
			splitter := make([]string, len(doc))
			for i := range doc {
				splitter[i] = strings.Repeat("-", max[i])
			}
			log.Print(strings.TrimSpace(strings.Join(splitter, " | ")))
		}
	}
}
