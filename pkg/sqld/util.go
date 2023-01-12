package sqld

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
)

const countField = "count(*) as count"

const OnConflictIDSET = "On Conflict (id) DO UPDATE SET"

func OnConflictSet(field string) string {
	return fmt.Sprintf("On Conflict (%s) DO UPDATE SET", field)
}

func ReplaceNotFound(err, notFoundErr error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFoundErr
	}

	return tracer.Trace(err)
}

func idEq(id any) squirrel.Eq {
	return squirrel.Eq{"id": id}
}

// ScanRows scans the db rows and returns list of it.
// I think we need to check performance of this method
// because of "T any". Maybe we should change any to Model
// and define a standard interface for Models or just
// autogenerate scanRows function per each model in the app.
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
