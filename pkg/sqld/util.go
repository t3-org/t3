package sqld

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
)

const countField = "count(*) as count"

func ReplaceNotFound(err, notFoundErr error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFoundErr
	}

	return tracer.Trace(err)
}

func idEq(id any) squirrel.Eq {
	return squirrel.Eq{"id": id}
}

func ScanRows[T any](rows *sql.Rows, fields func(m *T) []any) ([]*T, error) {
	defer rows.Close()
	var res []*T

	for rows.Next() {
		var m T
		if err := rows.Scan(fields(&m)...); err != nil {
			return nil, tracer.Trace(err)
		}
		res = append(res, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, tracer.Trace(err)
	}

	return res, nil
}
