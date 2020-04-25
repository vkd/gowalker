# Golang struct walker

[![Build Status](https://travis-ci.org/vkd/gowalker.svg)](https://travis-ci.org/vkd/gowalker)
[![codecov](https://codecov.io/gh/vkd/gowalker/branch/master/graph/badge.svg)](https://codecov.io/gh/vkd/gowalker)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkd/gowalker)](https://goreportcard.com/report/github.com/vkd/gowalker)
[![GoDoc](https://godoc.org/github.com/vkd/gowalker?status.svg)](https://godoc.org/github.com/vkd/gowalker)
[![Sourcegraph](https://sourcegraph.com/github.com/vkd/gowalker/-/badge.svg)](https://sourcegraph.com/github.com/vkd/gowalker?badge)
[![Release](https://img.shields.io/github/release/vkd/gowalker.svg)](https://github.com/vkd/gowalker/releases)

Walking throught golang struct

## Install

```sh
$ go get -u github.com/vkd/gowalker
```

## Examples

```go
package main

import (
	"log"
	"reflect"
	"time"

	"github.com/vkd/gowalker"
)

type Config struct {
	Addr    string
	Timeout time.Duration
	Start   time.Time

	Seconds *int

	DB *struct {
		Name string
		Port int `config:"db_port"`
	}
}

type configWalker map[string]string

func (c configWalker) Step(value reflect.Value, field reflect.StructField) (bool, error) {
	return gowalker.StringWalkerStep("config", gowalker.StringSourceMapString(c), value, field)
}

func main() {
	var cfg Config

	env := map[string]string{
		"Addr":    "localhost:9000",
		"Timeout": "5m",

		"Name":    "postgres",
		"db_port": "9001",
	}

	err := gowalker.Walk(&cfg, configWalker(env))
	if err != nil {
		log.Fatal(err)
	}

  ...
}
```

## Example of the config parsing

```go
type Config struct {
	Env int
	DB  struct {
		Port  int `default:"1000" env:"DB_PORT"`
	}
}

func ParseConfig() (Config, error) {
	var c Config
	err := config.Fill(&c,
		gowalker.Tag("default"),
		config.Name(
			gowalker.Fullname("_", strings.ToUpper),
			gowalker.NewStringWalker("env", gowalker.EnvFuncSource(os.LookupEnv)),
		),
		config.FlagWalker(
			gowalker.Tag("flag"),
			gowalker.Fullname("-", strings.ToLower),
		),
	)
	return c, err
}
```