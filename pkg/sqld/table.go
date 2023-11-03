package sqld

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
)

// TODO: move table implementation into an external go package. e.g., github.com/kamva/sqld (sql database)

type BuilderProvider func(ctx context.Context) sq.StatementBuilderType
type Table struct {
	name    string
	fields  []string
	builder BuilderProvider
}

func NewTable(name string, fields []string, builder BuilderProvider) *Table {
	return &Table{
		name:    name,
		fields:  fields,
		builder: builder,
	}
}

func (t *Table) FindByID(ctx context.Context, id any, dest []any) error {
	return tracer.Trace(t.builder(ctx).Select(t.fields...).From(t.name).Where(idEq(id)).Scan(dest...))
}

func (t *Table) First(ctx context.Context, dest []any, pred any, args ...any) error {
	return tracer.Trace(t.builder(ctx).Select(t.fields...).From(t.name).Where(pred, args...).Scan(dest...))
}

func (t *Table) Create(ctx context.Context, dest []any) (sql.Result, error) {
	return t.builder(ctx).Insert(t.name).Columns(t.fields...).Values(dest...).ExecContext(ctx)
}

func (t *Table) Upsert(ctx context.Context, dest []any, upsertPrefix string) (sql.Result, error) {
	expr, err := UpsertSuffix(upsertPrefix, Clauses(t.fields, dest)...)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return t.builder(ctx).Insert(t.name).Columns(t.fields...).Values(dest...).SuffixExpr(expr).ExecContext(ctx)
}

func (t *Table) CreateMany(ctx context.Context, dest ...[]any) (sql.Result, error) {
	b := t.builder(ctx).Insert(t.name).Columns(t.fields...)
	for _, d := range dest {
		b.Values(d...)
	}

	return b.ExecContext(ctx)
}

func (t *Table) Update(ctx context.Context, id any, dest []any) (sql.Result, error) {
	return t.UpdateBuilder(ctx, dest).Where(idEq(id)).ExecContext(ctx)
}

func (t *Table) UpdateBuilder(ctx context.Context, dest []any) sq.UpdateBuilder {
	update := t.builder(ctx).Update(t.name)
	for i, field := range t.fields {
		update = update.Set(field, dest[i])
	}

	return update
}

func (t *Table) Select(ctx context.Context, fields ...string) sq.SelectBuilder {
	if len(fields) == 0 {
		fields = t.fields
	}
	return t.builder(ctx).Select(fields...).From(t.name)
}

func (t *Table) Count(ctx context.Context) sq.SelectBuilder {
	return t.builder(ctx).Select(countField).From(t.name)
}

func (t *Table) Delete(ctx context.Context, id any) (sql.Result, error) {
	return t.DeleteBuilder(ctx).Where(idEq(id)).ExecContext(ctx)
}

func (t *Table) DeleteBuilder(ctx context.Context) sq.DeleteBuilder {
	return t.builder(ctx).Delete(t.name)
}
