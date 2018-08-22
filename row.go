package pg

import (
	"github.com/jackc/pgx"
)

type row struct {
	row *pgx.Row
}

func newRow(r *pgx.Row) Row {
	return &row{row: r}
}

func (r *row) Scan(args ...interface{}) Error {
	return newError(r.row.Scan(args...))
}
