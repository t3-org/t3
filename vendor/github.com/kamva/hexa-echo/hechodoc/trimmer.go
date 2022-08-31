package hechodoc

import (
	"errors"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
	"strings"
)

type TrimmerOptions struct {
	Echo                   *echo.Echo
	ExtractDestinationPath string
}

type Trimmer struct {
	echo     *echo.Echo
	filePath string
}

func NewTrimmer(o TrimmerOptions) *Trimmer {
	return &Trimmer{
		echo:     o.Echo,
		filePath: o.ExtractDestinationPath,
	}
}

func (t *Trimmer) Trim() error {
	file, err := os.OpenFile(t.filePath, os.O_RDWR, 0666)
	if err != nil {
		return tracer.Trace(err)
	}
	defer file.Close()

	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tracer.Trace(err)
	}
	builder := strings.Builder{}
	builder.Write(fBytes)

	oldRoutes := oldExtractedRoutes(fBytes)
	newRoutes := echoRawRouteNames(t.echo)

	for _, r := range oldRoutes {
		if !gutil.Contains(newRoutes, r) {
			if err := t.removeRoute(&builder, r); err != nil {
				hlog.With(hlog.String("route", r), hlog.Err(err)).Error("can not remove old route")
				return tracer.Trace(err)
			}
		}
	}

	if err := file.Truncate(0); err != nil {
		return tracer.Trace(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return tracer.Trace(err)
	}

	_, err = file.Write([]byte(builder.String()))
	return tracer.Trace(err)
}

func (t *Trimmer) removeRoute(b *strings.Builder, rName string) error {
	content := b.String()
	from := strings.Index(content, beginRouteVal(rName))
	to := strings.Index(content, endRouteVal(rName))
	if from == -1 || to == -1 {
		return errors.New("can not find beginning or ending tag")
	}
	b.Reset()
	b.WriteString(content[:from] + content[to+len(endRouteVal(rName)):])
	return nil
}
