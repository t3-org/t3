//go:generate easyjson -all
//--------------------------------
// We Got this package (pagination) from
// [go-rest-api](https://github.com/qiangxue/go-rest-api)
//--------------------------------

// Package pagination provides support for pagination requests and responses.
package pagination

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kamva/gutil"
)

var (
	// DefaultPageSize specifies the default page size
	DefaultPageSize = 30
	// MaxPageSize specifies the maximum page size
	MaxPageSize = 100
	// PageVar specifies the query parameter name for page number
	PageVar = "page"
	// PageSizeVar specifies the query parameter name for page size
	PageSizeVar = "per_page"
)

// Pagination represents pagination details of items.
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

type Pages struct {
	Pagination `json:"pagination"`
	Items      any `json:"items"`
}

// New creates a new Pages instance.
// The page parameter is 1-based and refers to the current page index/number.
// The perPage parameter refers to the number of items on each page.
// And the total parameter specifies the total number of data items.
// If total is less than 0, it means total is unknown.
func New(page, perPage, total int) *Pages {
	if perPage <= 0 {
		perPage = DefaultPageSize
	}
	if perPage > MaxPageSize {
		perPage = MaxPageSize
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Pages{
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			TotalCount: total,
			PageCount:  pageCount,
		},
	}
}

// NewFromRequest creates a Pages object using the query parameters found in the given HTTP request.
// count stands for the total number of items. Use -1 if this is unknown.
func NewFromRequest(req *http.Request, count int) *Pages {
	page := gutil.ParseInt(req.URL.Query().Get(PageVar), 1)
	perPage := gutil.ParseInt(req.URL.Query().Get(PageSizeVar), DefaultPageSize)
	return New(page, perPage, count)
}

// Offset returns the OFFSET value that can be used in a SQL statement.
func (p *Pages) Offset() int64 {
	return (int64(p.Page) - 1) * int64(p.PerPage)
}

// Limit returns the LIMIT value that can be used in a SQL statement.
func (p *Pages) Limit() int64 {
	return int64(p.PerPage)
}

// BuildLinkHeader returns an HTTP header containing the links about the pagination.
func (p *Pages) BuildLinkHeader(baseURL string, defaultPerPage int) string {
	links := p.BuildLinks(baseURL, defaultPerPage)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf("<%v>; rel=\"first\", ", links[0])
		header += fmt.Sprintf("<%v>; rel=\"prev\"", links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("<%v>; rel=\"next\"", links[2])
		if links[3] != "" {
			header += fmt.Sprintf(", <%v>; rel=\"last\"", links[3])
		}
	}
	return header
}

// BuildLinks returns the first, prev, next, and last links corresponding to the pagination.
// A link could be an empty string if it is not needed.
// For example, if the pagination is at the first page, then both first and prev links
// will be empty.
func (p *Pages) BuildLinks(baseURL string, defaultPerPage int) [4]string {
	var links [4]string
	pageCount := p.PageCount
	page := p.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, 1)
		links[1] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
		links[3] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
	}
	if perPage := p.PerPage; perPage != defaultPerPage {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&%v=%v", PageSizeVar, perPage)
			}
		}
	}

	return links
}
