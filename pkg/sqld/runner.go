package sqld

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Runner is an interface with common methods between a regular
// connection and a transaction.
type Runner interface {
	sq.Preparer
	sq.PreparerContext

	sq.Queryer
	sq.QueryerContext

	sq.Execer
	sq.ExecerContext

	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
