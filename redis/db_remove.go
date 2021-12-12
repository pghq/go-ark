package redis

import (
	"fmt"

	"github.com/pghq/go-ark/db"
)

func (tx txn) Remove(table, k string, _ ...db.CommandOption) error {
	tx.unit.Del(tx.ctx, fmt.Sprintf("%s.%s", table, k))
	return nil
}
