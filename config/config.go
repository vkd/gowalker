package config

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/vkd/gowalker"
	"github.com/vkd/gowalker/setter"
)

type Walker interface {
	Walk(v interface{}) error
}

func Fill(c interface{}, ws ...Walker) error {
	for _, w := range ws {
		err := w.Walk(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func FlagWalker(t gowalker.Tag, n gowalker.Namer) Walker {
	return flagWalker{tag: t, namer: n}
}

type flagWalker struct {
	tag   gowalker.Tag
	namer gowalker.Namer
}

func (w flagWalker) Walk(c interface{}) error {
	var flags = make(map[string]*string)
	var isNotSetString = "~~EMPTY~~"

	var flagWalker gowalker.WalkerFunc = func(value reflect.Value, field reflect.StructField, name gowalker.Name) (bool, error) {
		key, ok := field.Tag.Lookup(string(w.tag))
		if !ok {
			key = name.Get(w.namer)
		}
		s := flag.String(key, isNotSetString, field.Name)
		flags[field.Name] = s
		return false, nil
	}
	err := gowalker.Walk(c, flagWalker)
	if err != nil {
		return fmt.Errorf("error on flag walker: %w", err)
	}

	flag.Parse()

	var flagSetter gowalker.WalkerFunc = func(value reflect.Value, field reflect.StructField, name gowalker.Name) (bool, error) {
		key, ok := field.Tag.Lookup(string(w.tag))
		if !ok {
			key = name.Get(w.namer)
		}
		s, ok := flags[key]
		if !ok {
			return false, nil
		}
		if s == nil {
			return false, nil
		}
		if *s == isNotSetString {
			return false, nil
		}
		return true, setter.SetString(value, field, *s)
	}
	err = gowalker.Walk(c, flagSetter)
	if err != nil {
		return fmt.Errorf("error on flag setter: %w", err)
	}
	return nil
}
