package hechodoc

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/labstack/echo/v4"
	"path"
	"regexp"
	"sort"
	"strings"
)

type RouteNameConverter interface {
	Tags(name string) []string // give tags from route name
	CamelCase(name string) string
}

// DefaultRouteNameConverter is the default route name Converter.
var DefaultRouteNameConverter = NewDividerNameConverter("::", 0)

// dividedNameConverter convert names which their format is just
// multiple words which joined with a div like , or :.
type dividedNameConverter struct {
	div string
	// specify maximum tags count which we can return. -1 means unlimited.
	maxTag int
}

func NewDividerNameConverter(div string, maxTag int) RouteNameConverter {
	return &dividedNameConverter{
		div:    div,
		maxTag: maxTag,
	}
}

func (c *dividedNameConverter) Tags(name string) []string {
	tags := strings.Split(name, c.div)
	if c.maxTag == -1 {
		return tags
	}
	return tags[:c.maxTag]
}

func (c *dividedNameConverter) CamelCase(name string) string {
	return camelCaseFromStringList(strings.Split(name, c.div))
}

var _ RouteNameConverter = &dividedNameConverter{}

func echoRoutes(e *echo.Echo) []*echo.Route {
	routes := make([]*echo.Route, 0)
	for _, r := range e.Routes() {
		if !isEchoInternalRoute(r) {
			routes = append(routes, r)
		}
	}
	return routes
}

// isEchoInternalRoute detect whether route is for echo or defined by user.
// By default name of each route in echo is its handler name, so if
// route name is begins with "github.com/labstack/echo" its echo itself method.
func isEchoInternalRoute(r *echo.Route) bool {
	return strings.Index(r.Name, "github.com/labstack/echo") == 0
}

// echoRawRouteNames returns string list of echo route names.
func echoRawRouteNames(e *echo.Echo) []string {
	routes := make([]string, len(echoRoutes(e)))
	for i, r := range echoRoutes(e) {
		routes[i] = r.Name
	}

	return routes
}

// OpenApiPathFromEchoPath convert echo path to openapi path.
// e.g, a/:id/:code => a/{id}/{code}
func OpenApiPathFromEchoPath(path string) string {
	for {
		colonIndex := strings.Index(path, ":")
		if colonIndex == -1 {
			return path
		}
		path = gutil.ReplaceRune(path, '{', colonIndex)

		slashPath := path[colonIndex:]
		slashIndex := strings.Index(slashPath, "/")
		if slashIndex == -1 {
			slashIndex = len(slashPath)
		}
		slashPath = gutil.ReplaceAt(slashPath, "}", slashIndex, slashIndex)
		path = path[:colonIndex] + slashPath
	}
}

type Route struct {
	BeginRouteVal string
	EndRouteVal   string
	Name          string
	RawName       string
	Method        string
	Path          string
	PathParams    []PathParam
	TagsString    string
	ParamsId      string
	SuccessRespId string
}

type PathParam struct {
	Name         string
	ExportedName string
}

func newRoute(r *echo.Route, c RouteNameConverter) Route {
	p := path.Join("/", OpenApiPathFromEchoPath(r.Path))
	return Route{
		BeginRouteVal: beginRouteVal(r.Name),
		EndRouteVal:   endRouteVal(r.Name),
		Name:          c.CamelCase(r.Name),
		RawName:       r.Name,
		Method:        r.Method,
		Path:          p,
		PathParams:    PathParams(p),
		TagsString:    strings.Join(c.Tags(r.Name), " "),
		ParamsId:      fmt.Sprintf("%sParams", c.CamelCase(r.Name)),
		SuccessRespId: fmt.Sprintf("%sSuccessResponse", c.CamelCase(r.Name)),
	}
}

func beginRouteVal(rName string) string {
	return BeginRoutePrefix + rName
}

func endRouteVal(rName string) string {
	return EndRoutePrefix + rName
}

var pathRegex = regexp.MustCompile("{.*?}")

func PathParams(p string) []PathParam {
	pList := pathRegex.FindAllString(p, -1)
	l := make([]PathParam, len(pList))
	for i, v := range pList {
		v = strings.Trim(v, "{")
		v = strings.Trim(v, "}")

		l[i] = PathParam{
			Name:         v,
			ExportedName: strings.ToUpper(v[:1]) + v[1:],
		}
	}
	return l
}

func SortEchoRoutesByName(routes []*echo.Route) []*echo.Route {
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Name < routes[j].Name
	})
	return routes
}
