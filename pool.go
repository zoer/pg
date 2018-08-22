package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
)

type pool struct {
	pool *pgx.ConnPool
}

var _ = ConnPool(&pool{})

// NewPool creates a new DB pool
func NewPool(c PoolConfig) (ConnPool, error) {
	p, err := pgx.NewConnPool(c.config())
	if err != nil {
		return nil, err
	}

	return &pool{pool: p}, nil
}

func (p *pool) QueryRow(
	ctx context.Context, query string, args ...interface{}) Row {

	row := p.pool.QueryRowEx(ctx, query, nil, args...)
	return newRow(row)
}

func (p *pool) Query(
	ctx context.Context, query string, args ...interface{}) (Rows, Error) {

	defer recover()
	rows, err := p.pool.QueryEx(ctx, query, nil, args...)
	if err != nil {
		return nil, newError(err)
	}

	return newRows(rows), nil
}

func (p *pool) Exec(
	ctx context.Context, query string, args ...interface{}) (CommandTag, Error) {

	tag, err := p.pool.ExecEx(ctx, query, nil, args...)
	if err != nil {
		return nil, newError(err)
	}

	return newCommandTag(tag), nil
}

func (p *pool) CopyFrom(
	identifier []string, columns []string, data interface{}) error {

	iden := pgx.Identifier(identifier)

	var rows [][]interface{}

	switch t := data.(type) {
	case [][]interface{}:
		rows = t
	case [][]string:
		for _, r := range t {
			var rr []interface{}
			for _, v := range r {
				rr = append(rr, v)
			}
			rows = append(rows, rr)
		}
	default:
		return fmt.Errorf("CopyFrom doesn't work with the %T type", t)
	}

	_, err := p.pool.CopyFrom(iden, columns, pgx.CopyFromRows(rows))

	return err
}

func (p *pool) Close() error {
	p.pool.Close()
	return nil
}
