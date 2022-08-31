package gutil

import (
	"bytes"
	"text/template"
)

// RenderText parses the provided text and then execute it.
func RenderText(text string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := template.Must(template.New("").Parse(text)).Execute(&buf, data)
	return buf.String(), err
}
