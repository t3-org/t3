package hechodoc

import (
	"bufio"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"text/template"
)

const (
	BeginRoutePrefix = "// route:begin: "
	EndRoutePrefix   = "// route:end: "
)

var DefaultSingleRouteTemplatePath = path.Join(gutil.SourcePath(), "default_template.tpl")
var beginRouteRegex = regexp.MustCompile(fmt.Sprintf("%s(.+)", BeginRoutePrefix))

func oldExtractedRoutes(f []byte) []string {
	allMatches := beginRouteRegex.FindAllStringSubmatch(string(f), -1)
	routes := make([]string, len(allMatches))
	for i, v := range allMatches {
		routes[i] = v[1]
	}
	return routes
}

type Extractor struct {
	echo           *echo.Echo
	singleRouteTpl *template.Template
	dst            string // Destination path
	converter      RouteNameConverter
}

type ExtractorOptions struct {
	Echo                    *echo.Echo
	ExtractDestinationPath  string // Destination path of extract filePath
	SingleRouteTemplatePath string
	Converter               RouteNameConverter
}

func NewExtractor(o ExtractorOptions) *Extractor {
	fName := path.Base(o.SingleRouteTemplatePath)

	return &Extractor{
		echo:           o.Echo,
		singleRouteTpl: template.Must(template.New(fName).ParseFiles(o.SingleRouteTemplatePath)),
		dst:            o.ExtractDestinationPath,
		converter:      o.Converter,
	}
}

func (e *Extractor) Extract() error {
	file, err := os.OpenFile(e.dst, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return tracer.Trace(err)
	}
	defer file.Close()

	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tracer.Trace(err)
	}

	buf := bufio.NewWriter(file)
	oldRoutes := oldExtractedRoutes(fBytes)

	// TODO: check if echo has repetitive route name, return error, we must not have any repetitive routes.
	// TODO: on each route check if route name is valid (using name converter), if its not valid => log and ignore it.
	// append new routes
	for _, r := range SortEchoRoutesByName(echoRoutes(e.echo)) {
		if !gutil.Contains(oldRoutes, r.Name) {
			if err := e.addRoute(r,buf); err != nil {
				return tracer.Trace(err)
			}
		}
	}
	return tracer.Trace(buf.Flush())
}

func (e *Extractor) addRoute(r *echo.Route,w io.Writer) error {
	val :=newRoute(r,e.converter)
	return e.singleRouteTpl.Execute(w, val)
}
