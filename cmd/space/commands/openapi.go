package commands

import (
	"github.com/kamva/hexa-echo/hechodoc"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Manage and generate openapi docs",
}

var oaiExtractCmd = &cobra.Command{
	Use:     "extract",
	Short:   "Extract api routes and insert them to the openapi docs file",
	Example: "listen",
	RunE:    withApp(oaiExtractCmdF),
}

var oaiTrimCmd = &cobra.Command{
	Use:     "trim",
	Short:   "Remove old routes which we don't have in our echo routes from the doc file",
	Example: "trim",
	RunE:    withApp(oaiTrimCmdF),
}

var docsRouteNameConverter = hechodoc.NewDividerNameConverter("::", 1)

func init() {
	openapiCmd.AddCommand(oaiExtractCmd, oaiTrimCmd)

	rootCmd.AddCommand(openapiCmd)
}

func oaiExtractCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	cfg := o.Cfg
	e := bootEcho(cfg, o.SP, o.App)

	extractor := hechodoc.NewExtractor(hechodoc.ExtractorOptions{
		Echo:                    e.Echo,
		ExtractDestinationPath:  cfg.ApiDocExportFilePath(),
		SingleRouteTemplatePath: cfg.ApiDocsRouteTemplatePath(),
		Converter:               docsRouteNameConverter,
	})

	return tracer.Trace(extractor.Extract())
}

func oaiTrimCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	cfg := o.Cfg
	e := bootEcho(cfg, o.SP, o.App)

	trimmer := hechodoc.NewTrimmer(hechodoc.TrimmerOptions{
		Echo:                   e.Echo,
		ExtractDestinationPath: cfg.ApiDocExportFilePath(),
	})

	return tracer.Trace(trimmer.Trim())
}
