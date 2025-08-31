package transaction

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type key struct{}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(key{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
