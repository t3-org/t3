package api

import (
	"strconv"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/pagination"
	"github.com/labstack/echo/v4"
	apperr "space.org/space/internal/error"
)

// QueryVar is name of the query parameter that provide query value to us.
const QueryVar = "query"
const SortVar = "sort_by"

// Resource is a base Resource struct to use in other resources.
type Resource struct {
}

// ID extracts ID(and convert it to hexa id) from route path
func (r Resource) ID(c echo.Context) int64 {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(apperr.ErrInvalidIDValue.SetData(hexa.Map{"id": c.Param("id"), "in": "API path"}))
	}
	return id
}

// QueryAndPaginationParams gets pagination prams from request
func (r Resource) QueryAndPaginationParams(c echo.Context) (query string, page, pageSize int) {
	query = c.QueryParam(QueryVar)
	page = gutil.ParseInt(c.QueryParam(pagination.PageVar), 1)

	pageSize = gutil.Min(40, gutil.ParseInt(c.QueryParam(pagination.PageSizeVar), pagination.DefaultPageSize))
	return
}
