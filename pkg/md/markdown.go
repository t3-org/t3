package md

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Markdown struct {
	parserExtensions parser.Extensions
	rendererOptions  html.RendererOptions
}

func NewMarkdown(extensions parser.Extensions, rendererOptions html.RendererOptions) *Markdown {
	return &Markdown{
		parserExtensions: extensions,
		rendererOptions:  rendererOptions,
	}
}

func (m *Markdown) Render(md []byte) []byte {
	p := parser.NewWithExtensions(m.parserExtensions)
	r := html.NewRenderer(m.rendererOptions)
	return markdown.Render(p.Parse(md), r)
}

func (m *Markdown) RenderString(md string) string {
	return string(m.Render([]byte(md)))
}
