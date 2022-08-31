package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/lg"
	"github.com/kamva/tracer"
	"golang.org/x/mod/modfile"
)

const (
	LayerNameTx      = "tx"
	LayerNameTracing = "tracing"
)

var modulePath string
var here = gutil.SourcePath()

func main() {
	modSrc, err := gutil.ReadAll(path.Join(gutil.SourcePath(), "..", "..", "..", "go.mod"))
	if err != nil {
		fmt.Println(tracer.StackAsString(tracer.Trace(err)))
		os.Exit(1)
	}
	modulePath = modfile.ModulePath(modSrc)

	layers := []string{LayerNameTx, LayerNameTracing}

	if len(os.Args) > 1 {
		layers = strings.Split(os.Args[1], ",")
	}

	// trim list
	for _, v := range layers {
		l := strings.TrimPrefix(strings.TrimSuffix(v, " "), " ")
		switch l {
		case LayerNameTx:
			//generateTxLayer()
		case LayerNameTracing:
			if err := generateTracingLayer(); err != nil {
				fmt.Println(tracer.StackAsString(tracer.Trace(err)))
				os.Exit(1)
			}
		}
	}
}

//func generateTxLayer() {
//	tmpl := path.Join(gutil.SourcePath(), "tx_layer_mongo.tmpl")
//	src := path.Join(gutil.SourcePath(), "..", "app_iface.go")
//	output := path.Join(gutil.SourcePath(), "..", "tx_layer_mongo_gen.go")
//
//	metadata, err := lg.ExtractInterfaceMetadata(src, "App")
//	if err != nil {
//		fmt.Println(tracer.StackAsString(tracer.Trace(err)))
//		os.Exit(1)
//	}
//
//	data := &lg.TemplateData{
//		Package:   "app",
//		Name:      "mongoTxLayer",
//		Interface: metadata,
//	}
//
//	if err := lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true); err != nil {
//		fmt.Println(tracer.StackAsString(tracer.Trace(err)))
//		os.Exit(1)
//	}
//}

func generateTracingLayer() error {
	tmpl := path.Join(here, "tracing_layer.tmpl")
	src := path.Join(here, "..", "app_iface.go")
	output := path.Join(here, "..", "tracing_layer_gen.go")

	pkg, err := lg.NewPackageFromFilenames(path.Join(modulePath, "internal/app"), src)
	if err != nil {
		return tracer.Trace(err)
	}

	if err := lg.NewEmbeddedResolver(pkg).Resolve(); err != nil {
		return tracer.Trace(err)
	}

	_, app := pkg.FindInterface("App")

	data := &lg.TemplateData{
		Package:   "app",
		Name:      "tracingLayer",
		Interface: app,
	}

	return lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true)
}
