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
$ go get github.com/vkd/gowalker
```

## Example of the config parsing

```go
import (
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
	var c Config
	err := config.Fill(&c)
}
```
