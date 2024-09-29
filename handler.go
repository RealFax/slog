package log

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var TimeFormat = "[15:04:05.000]"

type (
	Attrs map[string]any

	groupOrAttrs struct {
		group string
		attrs []slog.Attr
	}

	Handler struct {
		addSource bool
		level     slog.Leveler
		mu        *sync.Mutex
		w         io.Writer

		goas []groupOrAttrs
	}

	Record struct {
		Group   string     `json:"group,omitempty"`
		Level   slog.Level `json:"level"`
		Time    time.Time  `json:"time"`
		Message string     `json:"message"`
		Attrs   Attrs      `json:"attrs"`
		Source  *string    `json:"source,omitempty"`
	}
)

func (a Attrs) String() string {
	if len(a) == 0 {
		return ""
	}

	if ReleaseMode {
		b, _ := Marshal(a)
		return string(b)
	}
	b, _ := MarshalIndent(a, "", "  ")
	return string(b)
}

func (r Record) Bytes() []byte {
	if ReleaseMode {
		b, _ := Marshal(r)
		return append(b, '\n')
	}
	bs := [][]byte{
		colorize(lightGray, r.Time.Format(TimeFormat)),
	}

	if r.Group != "" {
		bs = append(bs, colorize(yellow, "@"+r.Group+"@"))
	}

	bs = append(
		bs,
		colorize(getColor(r.Level), r.Level.String()+":"),
		colorize(white, r.Message),
		colorize(lightMagenta, r.Attrs.String()),
	)
	return buildOutput(bs...)
}

func (h *Handler) withGroupOrAttrs(goa groupOrAttrs) *Handler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *Handler) Handle(_ context.Context, record slog.Record) (err error) {
	output := Record{
		Level:   record.Level,
		Time:    record.Time,
		Message: record.Message,
		Attrs:   make(Attrs),
	}

	switch {
	case record.NumAttrs() > 0:
		record.Attrs(func(attr slog.Attr) bool {
			putAttr(output.Attrs, attr)
			return true
		})
	}

	group := make([]string, 0, len(h.goas))
	for _, goa := range h.goas {
		if goa.group != "" {
			group = append(group, goa.group)
		} else {
			for _, a := range goa.attrs {
				output.Attrs[a.Key] = a.Value
			}
		}
	}

	if len(group) != 0 {
		output.Group = strings.Join(group, ".")
	}

	if h.addSource {
		_, filename, line, ok := runtime.Caller(4)
		if ok {
			src := filename + ":" + strconv.Itoa(line)
			output.Source = &src
		}
	}

	_, err = h.w.Write(output.Bytes())
	return err
}

func NewHandler(opts *slog.HandlerOptions, w io.Writer) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	h := &Handler{
		addSource: opts.AddSource,
		mu:        &sync.Mutex{},
		w:         w,
		level:     opts.Level,
	}

	return h
}
