package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vkd/gowalker"
)

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
