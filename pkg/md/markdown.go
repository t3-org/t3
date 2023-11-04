package md

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Markdown struct {
	p *parser.Parser
	r markdown.Renderer
}

func NewMarkdown(parser *parser.Parser, renderer markdown.Renderer) *Markdown {
	return &Markdown{
		p: parser,
		r: renderer,
	}
}

func (m *Markdown) Render(md []byte) []byte {
	return markdown.Render(m.p.Parse(md), m.r)
}

func (m *Markdown) RenderString(md string) string {
	return string(m.Render([]byte(md)))
}
