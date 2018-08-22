package pg

import (
	"github.com/jackc/pgx"
)

type pgError struct {
	err error
}

var _ = Error(&pgError{})

func newError(e error) Error {
	if e == nil {
		return nil
	}

	return &pgError{err: e}
}

func (e *pgError) Error() string {
	return e.err.Error()
}

func (e *pgError) IsCode(code string) bool {
	pge, ok := e.err.(pgx.PgError)
	if !ok {
		return false
	}

	return pge.Code == code
}

func (e *pgError) NoRows() bool {
	return e != nil && e.err == pgx.ErrNoRows
}

func (e *pgError) Message() string {
	pge, ok := e.err.(pgx.PgError)
	if !ok {
		return ""
	}

	return pge.Message
}

func (e *pgError) Hint() string {
	pge, ok := e.err.(pgx.PgError)
	if !ok {
		return ""
	}

	return pge.Hint
}
