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
	LayerNameTracing         = "tracing"
	LayerNameModelDescriptor = "descriptor"
)

const (
	TracingLayerPackageName = "store"
	ModelsPackageName       = "model"
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

	layers := []string{LayerNameTracing, LayerNameModelDescriptor}

	if len(os.Args) > 1 {
		layers = strings.Split(os.Args[1], ",")
	}

	for _, v := range layers {
		l := strings.TrimPrefix(strings.TrimSuffix(v, " "), " ")
		switch l {
		case LayerNameTracing:
			if err := generateTracingLayer(); err != nil {
				fmt.Println(tracer.StackAsString(tracer.Trace(err)))
				os.Exit(1)
			}
		case LayerNameModelDescriptor:
			if err := generateModelDescriptorLayer(); err != nil {
				fmt.Println(tracer.StackAsString(tracer.Trace(err)))
				os.Exit(1)
			}
		}
	}

}

type SubStore struct {
	*lg.TemplateData
	StoreMethod   *lg.Method // method that we use in store to return this subStore.
	InterfaceName string     // Interface that this subStore will be implemented
	InterfaceType string
}

type TracingLayerData struct {
	Store     *lg.TemplateData
	SubStores map[string]*SubStore // map [storeMethodName] SubStoreData
}

func generateTracingLayer() error {
	tmpl := path.Join(here, "tracing_layer.tmpl")
	src := path.Join(here, "..", "..", "model")
	output := path.Join(here, "..", "tracing_layer_gen.go")

	pkg, err := lg.SinglePackageFromDir(path.Join(modulePath, "internal", ModelsPackageName), src)
	if err != nil {
		return tracer.Trace(err)
	}

	if err := lg.NewEmbeddedResolver(pkg).Resolve(); err != nil {
		return tracer.Trace(err)
	}

	_, store := pkg.FindInterface("Store")

	subStores := make(map[string]*SubStore)

	for _, m := range store.Methods {
		if m.Annotations.Lookup("subStore") != nil { // If the method returns a subStore
			ifaceName := m.Results[0].Type
			_, subStore := pkg.FindInterface(ifaceName)
			subStores[m.Name] = &SubStore{
				StoreMethod:   m,
				InterfaceName: ifaceName,
				InterfaceType: lg.SetPackageOnType(ModelsPackageName, ifaceName),
				TemplateData: &lg.TemplateData{
					Package:   TracingLayerPackageName,
					Name:      "tracingLayer" + ifaceName,
					Interface: lg.UseInterfaceInPackage(pkg, subStore),
				},
			}
		}
	}

	data := &TracingLayerData{
		Store: &lg.TemplateData{
			Package:   TracingLayerPackageName,
			Name:      "tracingLayerStore",
			Interface: lg.UseInterfaceInPackage(pkg, store),
		},
		SubStores: subStores,
	}

	return tracer.Trace(lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true))
}

type ModelDescriptor struct {
	Name   string
	Type   string
	Cols   []string // column names. Please note order of columns MUST be the same as fields order.
	Fields []string // Fields name.
}

func generateModelDescriptorLayer() error {
	tmpl := path.Join(here, "descriptor_layer.tmpl")
	src := path.Join(here, "..", "..", "model")
	output := path.Join(here, "..", "sqlstore", "descriptor_gen.go")

	pkg, err := lg.SinglePackageFromDir(path.Join(modulePath, "internal/model"), src)
	if err != nil {
		return tracer.Trace(err)
	}

	// resolve embedded structs:
	// Right now we allow all embedded fields to be resolved, but if
	// we wanted to ignore fields with specific annotation, we can use
	// lg.EmbeddedResolveFilter.
	if err := lg.NewEmbeddedResolver(pkg).Resolve(); err != nil {
		return tracer.Trace(err)
	}

	var modelsFromStore []string
	_, store := pkg.FindInterface("Store")
	for _, m := range store.Methods {
		if a := m.Annotations.Lookup("subStore"); a != nil {
			modelsFromStore = append(modelsFromStore, m.Name)
		}
	}

	var descriptors []*ModelDescriptor
	// Find all models with @model annotation.
	for _, f := range pkg.Files {
		for _, s := range f.Structs {
			if a := s.Annotations.Lookup("model"); a != nil || gutil.Contains(modelsFromStore, s.Name) {
				descriptors = append(descriptors, newDescriptor(ModelsPackageName, s))
			}
		}
	}

	data := map[string]any{"Descriptors": descriptors}
	return tracer.Trace(lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true))
}

func newDescriptor(structPkg string, s *lg.Struct) *ModelDescriptor {
	var columns []string
	var fields []string
	for _, field := range s.Fields {
		if lg.IsPrivate(field.Name) {
			continue
		}

		col, ok := lg.Lookup(field.Tag, "sql", "json")
		if ok && (col == "" || col == "-") { // skip for existed tag with emtpy or dash
			continue
		}

		columns = append(columns, gutil.StringDefault(col, gutil.ToSnakeCase(field.Name)))
		fields = append(fields, field.Name)
	}

	return &ModelDescriptor{
		Name:   gutil.LowerFirst(s.Name),
		Type:   lg.SetPackageOnType(structPkg, s.Name),
		Cols:   columns,
		Fields: fields,
	}
}
