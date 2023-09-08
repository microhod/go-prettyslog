package prettyslog

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/microhod/go-prettyslog/colour"
)

const (
	TemplateMultilineColourised = `{{ levelcolour (padleft .Level 5) . }}[{{ .Time }}] {{ colour .Message 1 }}
{{- range .Attrs }}
{{ spaces (len (padleft $.Level 5)) (len $.Time) 3 }}{{ colour .Key 30 }}{{ colour "=" 30 }}{{ colour .Value 30 }}
{{- end }}
`
)

type TemplateRecordWriter struct {
	Name          string
	Template      string
	TemplateFuncs template.FuncMap
}

func (t TemplateRecordWriter) WriteRecord(w io.Writer, record Record) error {
	recordTemplate, err := template.New(t.Name).Funcs(t.TemplateFuncs).Parse(t.Template)
	if err != nil {
		return fmt.Errorf("failed to parse log template: %w", err)
	}
	return recordTemplate.Execute(w, record)
}

var defaultTemplateFuncs = template.FuncMap{
	"levelcolour": LevelColourise,
	"colour":      Colourise,
	"padleft":     PadLeft,
	"padright":    PadRight,
	"spaces":      Spaces,
}

func LevelColourise(value string, record Record) string {
	return record.LevelColours.Sprint(value)
}

func Colourise(value string, colours ...colour.Colour) string {
	return colour.Colours(colours).Sprint(value)
}

func PadLeft(value string, length int) string {
	return pad("+", value, length)
}

func PadRight(value string, length int) string {
	return pad("-", value, length)
}

func pad(direction, value string, length int) string {
	format := `%` + direction + fmt.Sprint(length) + `s`
	return fmt.Sprintf(format, value)
}

func Spaces(lengths ...int) string {
	var length int
	for _, l := range lengths {
		length += l
	}
	return strings.Repeat(" ", length)
}
