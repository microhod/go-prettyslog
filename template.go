package prettyslog

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/microhod/go-prettyslog/colour"
)

const (
	TemplateMultilineColourised = `{{ levelcolour . }}{{ padleft .Level 5 }}{{ uncolour }}[{{ .Time }}] {{ colour 1 }}{{ .Message }}{{ uncolour }}
{{- range .Attrs }}
	{{- printf "\n" }}
	{{- spaces (len (padleft $.Level 5)) (len $.Time) 3 }}
	{{- colour 30 }}
		{{- join .Groups "." }}
		{{- if .Groups }}.{{ end }}
		{{- .Key }}={{ .Value }}
	{{- uncolour }}
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
	"levelcolour": LevelColourStart,
	"colour":      ColourStart,
	"uncolour":    ColourEnd,
	"padleft":     PadLeft,
	"padright":    PadRight,
	"spaces":      Spaces,
	"join":        strings.Join,
}

func LevelColourStart(record Record) string {
	return record.LevelColours.Format()
}

func ColourStart(colours ...colour.Colour) string {
	return colour.Colours(colours).Format()
}

func ColourEnd() string {
	return colour.Unformat()
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
