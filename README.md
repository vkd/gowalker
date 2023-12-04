# Golang struct walker

[![Build Status](https://travis-ci.org/vkd/gowalker.svg)](https://travis-ci.org/vkd/gowalker)
[![codecov](https://codecov.io/gh/vkd/gowalker/branch/master/graph/badge.svg)](https://codecov.io/gh/vkd/gowalker)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkd/gowalker)](https://goreportcard.com/report/github.com/vkd/gowalker)
[![GoDoc](https://godoc.org/github.com/vkd/gowalker?status.svg)](https://godoc.org/github.com/vkd/gowalker)
[![Sourcegraph](https://sourcegraph.com/github.com/vkd/gowalker/-/badge.svg)](https://sourcegraph.com/github.com/vkd/gowalker?badge)
[![Release](https://img.shields.io/github/release/vkd/gowalker.svg)](https://github.com/vkd/gowalker/releases)

Walking throught golang struct to fullfil its fields from ENV variables.

## Install

```sh
$ go get github.com/vkd/gowalker
```

## Example of the config parsing

```go
import (
	"log"
	"os"

	"github.com/vkd/gowalker"
	"github.com/vkd/gowalker/config"
)

type Config struct {
	LogLevel string        `flag:"loglevel" env:"LOGLEVEL" required:"true"`
	Timeout  time.Duration `default:"3s"`

	DB  struct {
		Port  int `default:"5432" flag:"db-port" env:"DB_PORT"`
	}
}

func ParseConfig() {
	var cfg Config
	err := config.Walk(&cfg, log.New(os.Stdout, "", 0),
		gowalker.Flags(gowalker.FieldKey("flag", gowalker.Fullname("-", strings.ToLower)), os.Args),
		gowalker.Envs(gowalker.FieldKey("env", gowalker.Fullname("_", strings.ToUpper)), os.LookupEnv),
		gowalker.Tag("default"),
		gowalker.Required("required"),
	)
	if err != nil {
		if errors.Is(err, gowalker.ErrPrintHelp) {
			return nil
		}
		...
	}
}
```
