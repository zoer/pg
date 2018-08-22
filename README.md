# pg
[![Build Status](https://travis-ci.org/zoer/pg.svg)](https://travis-ci.org/zoer/pg)
[![Go Report
Card](https://goreportcard.com/badge/github.com/zoer/pg)](https://goreportcard.com/report/github.com/zoer/pg)
[![GoDoc](https://godoc.org/github.com/zoer/pg?status.svg)](https://godoc.org/github.com/zoer/pg)

https://github.com/jackc/pgx wrapper for internal usage

## Install

```
$ go get github.com/zoer/pg
```
## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/zoer/pg"
)

func main() {
	// creating a new DB pool
	pool, err := pg.NewPool(pg.PoolConfig{
		ConnConfig: ConnConfig{
			Host:     "127.0.0.1",
			Database: "test",
		},
		MaxConnections: 3,
	})
	if err != nil {
		log.Fatalf("unable connect to database: %v", err)
	}
	defer pool.Close()

	var i int
	if err = pool.QueryRow(context.TODO(), "SELECT $1", 123).Scan(&i); err != nil {
		log.Fatalf("unable to query: %v", err)
	}
}
```
