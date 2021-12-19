package sql

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pghq/go-tea"

	"github.com/pghq/go-ark/db"
)

func (tx txn) Get(table string, k db.Key, v interface{}, opts ...db.QueryOption) error {
	if tx.err != nil {
		return tea.Error(tx.err)
	}

	query := db.QueryWith(opts)
	stmt, args, err := squirrel.StatementBuilder.
		Select().
		From(table).
		Limit(1).
		Columns(query.Fields...).
		PlaceholderFormat(placeholder(query.SQLPlaceholder)).
		Where(squirrel.Eq{k.Name: k.Value}).
		ToSql()
	if err != nil {
		return tea.NewError(err)
	}

	if err := tx.unit.GetContext(tx.ctx, v, stmt, args...); err != nil {
		if err == sql.ErrNoRows {
			return tea.NotFound(err)
		}
		return err
	}

	return nil
}
