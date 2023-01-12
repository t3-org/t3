package main

import (
	"log"
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

var here = gutil.SourcePath()
var modPath = path.Join(gutil.SourcePath(), "../../../go.mod")
var modulePath = modfile.ModulePath(gutil.Must(os.ReadFile(modPath)).([]byte))
var layers = map[string]func() error{
	//LayerNameTx:      generateTxLayer,
	LayerNameTracing: generateTracingLayer,
}

func main() {
	layerNames := []string{LayerNameTx, LayerNameTracing}
	if len(os.Args) > 1 {
		layerNames = strings.Split(os.Args[1], ",")
	}

	if err := lg.RunLayerFns(layers, layerNames...); err != nil {
		log.Fatal(tracer.Trace(err))
	}
}

//func generateTxLayer() error { // TODO: generate tx layer from both interface and the app instance.
//	tmpl := path.Join(gutil.SourcePath(), "tx_layer_mongo.tmpl")
//	src := path.Join(gutil.SourcePath(), "..", "app_iface.go")
//	output := path.Join(gutil.SourcePath(), "..", "tx_layer_mongo_gen.go")
//
//	metadata, err := lg.ExtractInterfaceMetadata(src, "App")
//	if err != nil {
//		return tracer.Trace(err)
//	}
//
//	data := &lg.TemplateData{
//		Package:   "app",
//		Name:      "mongoTxLayer",
//		Interface: metadata,
//	}
//
//	return tracer.Trace(lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true))
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

	return tracer.Trace(lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true))
}
