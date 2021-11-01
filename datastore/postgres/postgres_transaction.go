package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pghq/go-museum/museum/diagnostic/errors"

	"github.com/pghq/go-datastore/datastore/client"
)

// Transaction creates a database transaction for Postgres.
func (c *Client) Transaction(ctx context.Context) (client.Transaction, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	t := transaction{ctx: ctx, tx: tx}
	return &t, err
}

type transaction struct {
	ctx context.Context
	tx  pgx.Tx
}

func (t *transaction) Execute(statement client.Encoder) (int, error) {
	sql, args, err := statement.Statement()
	if err != nil {
		return 0, errors.BadRequest(err)
	}

	tag, err := t.tx.Exec(t.ctx, sql, args...)
	if err != nil {
		if IsIntegrityConstraintViolation(err) {
			return 0, errors.BadRequest(err)
		}
		return 0, errors.Wrap(err)
	}

	return int(tag.RowsAffected()), nil
}

func (t *transaction) Commit() error {
	if err := t.tx.Commit(t.ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil

}
func (t *transaction) Rollback() error {
	if err := t.tx.Rollback(t.ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
