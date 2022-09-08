package sqld

import (
	"bytes"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
)

type SetClause struct {
	column string
	value  any
}

func Clauses(names []string, dest []any) []*SetClause {
	l := make([]*SetClause, len(names))

	for i, v := range names {
		l[i] = &SetClause{column: v, value: dest[i]}
	}

	return l
}

func UpsertSuffix(prefix string, clauses ...*SetClause) (sq.Sqlizer, error) {
	sql, args, err := upsertSuffix(prefix, clauses...)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return sq.Expr(sql, args...), nil
}

func upsertSuffix(prefix string, clauses ...*SetClause) (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}
	sql.WriteString(prefix + " ")

	setSqls := make([]string, len(clauses))
	for i, setClause := range clauses {
		var valSql string
		if vs, ok := setClause.value.(sq.Sqlizer); ok {
			vsql, vargs, err := vs.ToSql()
			if err != nil {
				return "", nil, err
			}
			if _, ok := vs.(sq.SelectBuilder); ok {
				valSql = fmt.Sprintf("(%s)", vsql)
			} else {
				valSql = vsql
			}
			args = append(args, vargs...)
		} else {
			valSql = "?"
			args = append(args, setClause.value)
		}
		setSqls[i] = fmt.Sprintf("%s = %s", setClause.column, valSql)
	}

	sql.WriteString(strings.Join(setSqls, ", "))
	return sql.String(), args, nil
}
