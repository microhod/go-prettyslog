package prettyslog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"

	"github.com/microhod/go-prettyslog/colour"
)

type Handler struct {
	w       io.Writer
	options HandlerOptions

	attrs  []Attr
	groups []string
}

type HandlerOptions struct {
	Level     slog.Leveler
	AddSource bool

	LevelColours    map[slog.Level]colour.Colours
	TimestampFormat string

	RecordWriter RecordWriter
}

type RecordWriter interface {
	WriteRecord(w io.Writer, record Record) error
}

type HandlerOptionsFunc func(opts *HandlerOptions)

var defaultHandlerOptions = HandlerOptions{
	Level:     slog.LevelInfo,
	AddSource: false,

	LevelColours: map[slog.Level]colour.Colours{
		slog.LevelDebug: {colour.FgBlack},
		slog.LevelInfo:  {colour.FgCyan},
		slog.LevelWarn:  {colour.FgYellow},
		slog.LevelError: {colour.FgRed},
	},
	TimestampFormat: "2006-01-02T15:04:05.000Z",

	RecordWriter: TemplateRecordWriter{
		Name:          "prettyslog-log-template",
		Template:      TemplateMultilineColourised,
		TemplateFuncs: defaultTemplateFuncs,
	},
}

func NewHandler(w io.Writer, options ...HandlerOptionsFunc) *Handler {
	handler := Handler{
		w:       w,
		options: defaultHandlerOptions,
	}

	for _, option := range options {
		option(&handler.options)
	}
	return &handler
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := h.clone()
	for _, attr := range attrs {
		handler.attrs = append(handler.attrs, Attr{
			Key:   attr.Key,
			Value: attr.Value.String(),
			Groups: handler.groups,
		})
	}
	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	handler := h.clone()
	handler.groups = append(handler.groups, name)
	return handler
}

func (h *Handler) clone() *Handler {
	return &Handler{
		options: h.options,
		w:       h.w,
		attrs:   h.attrs,
		groups:  h.groups,
	}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	// use a buffer to ensure we only ever write a full log line at once to the handler's writer
	buffer := new(bytes.Buffer)
	err := h.options.RecordWriter.WriteRecord(buffer, h.recordFromSlogRecord(record))
	if err != nil {
		return fmt.Errorf("failed to execute log template: %w", err)
	}

	_, err = h.w.Write(buffer.Bytes())
	return err
}

type Record struct {
	LevelColours colour.Colours

	Time    string
	Level   string
	Message string
	Attrs   []Attr
}

type Attr struct {
	Key, Value string
	Groups []string
}

func (h *Handler) recordFromSlogRecord(slogRecord slog.Record) Record {
	var record Record

	if !slogRecord.Time.IsZero() {
		record.Time = slogRecord.Time.Format(h.options.TimestampFormat)
	}
	record.Level = slogRecord.Level.String()
	record.Message = slogRecord.Message
	record.Attrs = h.attrsFromSlogRecord(slogRecord)

	// colour
	var exists bool
	record.LevelColours, exists = h.options.LevelColours[slogRecord.Level]
	if !exists {
		// default to black
		record.LevelColours = colour.Colours{colour.FgBlack}
	}
	return record
}

func (h *Handler) attrsFromSlogRecord(slogRecord slog.Record) []Attr {
	var attrs []Attr

	// source
	if h.options.AddSource && slogRecord.PC != 0 {
		attrs = append(attrs, Attr{Key: slog.SourceKey, Value: h.source(slogRecord.PC)})
	}

	// logger attrs
	for _, attr := range h.attrs {
		attrs = append(attrs, attr)
	}
	// record attrs
	slogRecord.Attrs(func(attr slog.Attr) bool {
		// skip empty attrs
		if attr.Equal(slog.Attr{}) {
			return true
		}
		attrs = append(attrs, Attr{Key: attr.Key, Value: attr.Value.String(), Groups: h.groups})
		return true
	})
	return attrs
}

// source returns a source code location in the form 'file:linenumber' for the given program counter.
func (h *Handler) source(programCounter uintptr) string {
	fs := runtime.CallersFrames([]uintptr{programCounter})
	f, _ := fs.Next()
	return fmt.Sprintf("%s:%d", f.File, f.Line)
}
