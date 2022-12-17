package sqld

import (
	"context"
	"database/sql"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type contextTxKey string

const (
	keyTx       contextTxKey = "_tx"
	keyNestedTx contextTxKey = "_nested_tx"
)

type TxWrapper struct {
	db *sql.DB
}

func NewTxWrapper(db *sql.DB) *TxWrapper {
	return &TxWrapper{
		db: db,
	}
}

func (t *TxWrapper) Begin(ctx context.Context, opts *sql.TxOptions) (context.Context, error) {
	if isNestedTx(ctx) { // If another nested transaction is already started.
		return ctx, nil
	}

	if TxFromCtx(ctx) != nil { // If another transaction is already started.
		return context.WithValue(ctx, keyNestedTx, struct{}{}), nil
	}

	tx, err := t.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return context.WithValue(ctx, keyTx, tx), nil
}

// Rollback the Transaction.
// If you had nested transaction, rollback on the nested
// transactions does nothing, so you MUST return an error
// to the method which is called your method, to get the
// error and rollback the parent transaction.
func (t *TxWrapper) Rollback(ctx context.Context) error {
	// If is a nested Tx, ignore this call
	if isNestedTx(ctx) {
		return nil
	}

	err := TxFromCtx(ctx).Rollback()
	if err != nil && err != sql.ErrTxDone {
		hexa.CtxLogger(ctx).Error("can not rollback tx", hlog.Err(err))
	}

	return tracer.Trace(err)
}

// Commit the Transaction.
// If you had nested transaction, commit on the nested
// transaction does nothing.
func (t *TxWrapper) Commit(ctx context.Context) error {
	// If is a nested Tx, ignore this call
	if isNestedTx(ctx) {
		return nil
	}

	return tracer.Trace(TxFromCtx(ctx).Commit())
}

// Finalize tries to commit if the returned error is nil, and if the
// commit had error, it'll assign the error to the returned error.
// But if the returned error is not nil, it'll roll back the transaction.
func (t *TxWrapper) Finalize(ctx context.Context, returnedErr *error) {
	if *returnedErr != nil {
		if err := t.Rollback(ctx); err != nil && err != sql.ErrTxDone {
			hexa.CtxLogger(ctx).Error("can not rollback")
		}
	}

	if err := t.Commit(ctx); err != nil {
		*returnedErr = tracer.Trace(err)
	}
}

func TxFromCtx(ctx context.Context) *sql.Tx {
	if val := ctx.Value(keyTx); val != nil {
		return val.(*sql.Tx)
	}

	return nil
}

func isNestedTx(ctx context.Context) bool {
	return ctx.Value(keyNestedTx) != nil
}
