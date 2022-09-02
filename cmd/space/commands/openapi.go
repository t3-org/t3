package commands

import (
	"github.com/kamva/hexa-echo/hechodoc"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Manage and generate openapi docs",
}

var oaiExtractCmd = &cobra.Command{
	Use:     "extract",
	Short:   "Extract api routes and insert them to the openapi docs file",
	Example: "listen",
	RunE:    withApp(oaExtractCmdF),
}

var oaiTrimCmd = &cobra.Command{
	Use:     "trim",
	Short:   "Remove old routes which we don't have in our echo routes from the doc file",
	Example: "trim",
	RunE:    withApp(oaTrimCmdF),
}

var docsRouteNameConverter = hechodoc.NewDividerNameConverter("::", 1)

func init() {
	openapiCmd.AddCommand(oaiExtractCmd, oaiTrimCmd)

	rootCmd.AddCommand(openapiCmd)
}

func oaExtractCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	cfg := o.Cfg
	if err := registry.Provide(registry.Registry(), provider.HttpServerProvider); err != nil {
		return tracer.Trace(err)
	}

	extractor := hechodoc.NewExtractor(hechodoc.ExtractorOptions{
		Echo:                    o.SP.HttpServer().Echo,
		ExtractDestinationPath:  cfg.ApiDocExportFilePath(),
		SingleRouteTemplatePath: cfg.ApiDocsRouteTemplatePath(),
		Converter:               docsRouteNameConverter,
	})

	return tracer.Trace(extractor.Extract())
}

func oaTrimCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	cfg := o.Cfg
	if err := registry.Provide(registry.Registry(), provider.HttpServerProvider); err != nil {
		return tracer.Trace(err)
	}

	trimmer := hechodoc.NewTrimmer(hechodoc.TrimmerOptions{
		Echo:                   o.SP.HttpServer().Echo,
		ExtractDestinationPath: cfg.ApiDocExportFilePath(),
	})

	return tracer.Trace(trimmer.Trim())
}
