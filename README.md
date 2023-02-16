# Dotenv

[![go test](https://github.com/tangelo-labs/go-dotenv/actions/workflows/go-test.yml/badge.svg)](https://github.com/tangelo-labs/go-dotenv/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/tangelo-labs/go-dotenv/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/tangelo-labs/go-dotenv/actions/workflows/golangci-lint.yml)

This package provides a simple mechanism for loading environment variables values into typed struct fields.

## Installation

```bash
go get github.com/tangelo-labs/go-dotenv
```

## Usage

The following example assumes that either variables are set or that a `.env` file exists in the current working
directory, or are already defined in the environment.

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tangelo-labs/go-dotenv"
)

type config struct {
	Foo  string     `env:"ENV_FOO,required" default:"fooValue"`
	Bar  int        `env:"ENV_BAR,notEmpty"`
	IPs  []string	`env:"ENV_IPS" delimiter:";"`
	When time.Time	`env:"ENV_WHEN" default:"2021-12-24T17:04:05Z07:00" timeLayout:"2006-01-02T15:04:05Z07:00"`
}

func main() {
	var cfg config

	if err := dotenv.LoadAndParse(&cfg); err != nil {
		panic(err)
	}

	fmt.Printf("Foo: %s\n", cfg.Foo)
	fmt.Printf("Bar: %d\n", cfg.Bar)
	fmt.Printf("IPs: %+v\n", cfg.IPs)
	fmt.Printf("When: %s\n", cfg.When)
}
```

See the `dotenv.Parse` function for further details.
