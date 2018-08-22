package pg

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnPool_Exec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	a := require.New(t)
	p := newTestPool()
	defer p.Close()

	_, err := p.Exec(ctx, "SELECT $1::int", 1)
	a.NoError(err)

	_, err = p.Exec(ctx, "SELECT $1::foo", 1)
	a.True(err.IsCode("42704"))
}

func TestConnPool_QueryRow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	a := require.New(t)
	p := newTestPool()
	defer p.Close()

	var i int
	err := p.QueryRow(ctx, "SELECT $1::int", 1).Scan(&i)
	a.NoError(err)
	a.Equal(i, 1)

	err = p.QueryRow(ctx, "SELECT $1::foo", 1).Scan(&i)
	a.True(err.IsCode("42704"))
	a.False(err.NoRows())

	q := "SELECT 1 WHERE FALSE"
	err = p.QueryRow(ctx, q).Scan(&i)
	a.True(err.NoRows())
}

func TestConnPool_Query(t *testing.T) {
	t.Parallel()

	p := newTestPool()
	defer p.Close()

	t.Run("Scan", func(t *testing.T) {
		a := require.New(t)

		var i int
		rows, err := p.Query(context.TODO(), "SELECT $1::int", 1)
		a.NoError(err)
		a.True(rows.Next())
		rows.Scan(&i)
		a.Equal(i, 1)
		a.False(rows.Next())
		rows.Close()
	})

	t.Run("Map with valid iteration", func(t *testing.T) {
		a := require.New(t)

		var ints []int
		rows, err := p.Query(context.TODO(), "SELECT unnest('{1,2}'::int[]);")
		a.NoError(err)
		err = rows.Map(func(r Row) error {
			var i int
			a.NoError(r.Scan(&i))
			ints = append(ints, i)
			return nil
		})
		a.NoError(err)
		a.Equal(ints, []int{1, 2})
	})

	t.Run("Map, when there's an error during iteration", func(t *testing.T) {
		a := require.New(t)

		var ints []int
		rows, err := p.Query(context.TODO(), "SELECT unnest('{3,4}'::int[]);")
		a.NoError(err)
		err = rows.Map(func(r Row) error {
			return errors.New("foo")
		})
		a.Equal(err.Error(), "foo")
		a.Empty(ints)
	})

	t.Run("invalid type", func(t *testing.T) {
		a := require.New(t)

		_, err := p.Query(context.TODO(), "SELECT 1::foo", 1)
		a.True(err.IsCode("42704"))
	})
}

func TestConnPool_Close(t *testing.T) {
	ctx := context.Background()
	a := require.New(t)

	p := newTestPool()
	defer p.Close()

	_, err := p.Exec(ctx, "SELECT 1")
	a.NoError(err)

	a.NoError(p.Close())

	_, err = p.Exec(ctx, "SELECT 1")
	a.Error(err)
	a.Regexp("closed", err.Error())
}

func TestConnPool_CopyFrom(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	a := require.New(t)

	p := newTestPool()
	defer p.Close()

	_, err := p.Exec(ctx, "CREATE TABLE copy_table(id serial, c1 text, c2 text)")
	a.NoError(err)
	defer p.Exec(ctx, "DROP TABLE copy_table")

	test := func(a *require.Assertions) {
		var count int

		err = p.QueryRow(ctx,
			"SELECT COUNT(*) FROM copy_table WHERE id=1 AND c1='1' AND c2='2'").Scan(&count)
		a.NoError(err)
		a.Equal(count, 1)

		err = p.QueryRow(ctx,
			"SELECT COUNT(*) FROM copy_table WHERE id=2 AND c1='11' AND c2='22'").Scan(&count)
		a.NoError(err)
		a.Equal(count, 1)
	}

	t.Run("strings slice", func(t *testing.T) {
		aa := require.New(t)

		err := p.CopyFrom(
			[]string{"public", "copy_table"},
			[]string{"c1", "c2"},
			[][]string{[]string{"1", "2"}, []string{"11", "22"}})
		aa.NoError(err)

		test(aa)
	})

	t.Run("interfaces slice", func(t *testing.T) {
		aa := require.New(t)

		err := p.CopyFrom(
			[]string{"copy_table"},
			[]string{"c1", "c2"},
			[][]interface{}{[]interface{}{"1", "2"}, []interface{}{"11", "22"}})
		aa.NoError(err)

		test(aa)
	})
}

func newTestPool() *pool {
	p, err := NewPool(PoolConfig{
		ConnConfig: ConnConfig{
			Host:     "127.0.0.1",
			Database: "test",
		},
		MaxConnections: 3,
	})
	if err != nil {
		panic(err)
	}

	return p.(*pool)
}
