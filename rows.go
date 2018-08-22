package pg

import (
	"github.com/jackc/pgx"
)

type rows struct {
	rows *pgx.Rows
}

var _ = Rows(&rows{})

func newRows(pRows *pgx.Rows) Rows {
	return &rows{rows: pRows}
}

func (r *rows) Scan(args ...interface{}) Error {
	return newError(r.rows.Scan(args...))
}

func (r *rows) Next() bool {
	return r.rows.Next()
}

func (r *rows) Err() Error {
	return newError(r.rows.Err())
}

func (r *rows) Close() {
	r.rows.Close()
}
