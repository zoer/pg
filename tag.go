package pg

import (
	"github.com/jackc/pgx"
)

type commandTag struct {
	tag pgx.CommandTag
}

var _ = CommandTag(&commandTag{})

func newCommandTag(tag pgx.CommandTag) CommandTag {
	return &commandTag{tag}
}

func (t *commandTag) RowsAffected() int64 {
	return t.tag.RowsAffected()
}
