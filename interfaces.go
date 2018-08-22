package pg

import (
	"context"
	"io"
)

//go:generate mockgen -destination=mock/pg.go -package=mock -imports=.=github.com/zoer/pg -source=interfaces.go ConnPool,Row,Error

// ConnPool represents database pool interface
type ConnPool interface {
	io.Closer
	QueryRow(context.Context, string, ...interface{}) Row
	Query(context.Context, string, ...interface{}) (Rows, Error)
	Exec(context.Context, string, ...interface{}) (CommandTag, Error)
	// Acquire() (Conn, Error)
	// Release(Conn)

	CopyFrom([]string, []string, interface{}) error
}

// CommandTag is a command's result interface
type CommandTag interface {
	RowsAffected() int64
}

// Row is a database row interface
type Row interface {
	Scan(...interface{}) Error
}

// Rows is a database rows interface
type Rows interface {
	Map(func(Row) error) Error
	Next() bool
	Scan(...interface{}) Error
	Close()
	Err() Error
}

// Error is a database error interface
type Error interface {
	Error() string
	Message() string
	Hint() string
	IsCode(string) bool
	NoRows() bool
}
