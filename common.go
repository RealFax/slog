package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

var bufferPool = sync.Pool{New: func() any { return &bytes.Buffer{} }}

func buildOutput(bs ...[]byte) []byte {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	for _, b := range bs {
		buf.Write(b)
		buf.WriteByte(' ')
	}
	buf.WriteByte('\n')

	b := buf.Bytes()
	return b
}

//func suppressDefaults(
//	next func([]string, slog.Attr) slog.Attr,
//) func([]string, slog.Attr) slog.Attr {
//	return func(groups []string, a slog.Attr) slog.Attr {
//		if a.Method == slog.TimeKey ||
//			a.Method == slog.LevelKey ||
//			a.Method == slog.MessageKey {
//			return slog.Attr{}
//		}
//		if next == nil {
//			return a
//		}
//		return next(groups, a)
//	}
//}

func putAttr(attrs Attrs, attr slog.Attr) {
	var v any
	switch attr.Value.Kind() {
	case slog.KindString:
		v = attr.Value.String()
	case slog.KindInt64:
		v = attr.Value.Int64()
	case slog.KindUint64:
		v = attr.Value.Uint64()
	case slog.KindFloat64:
		v = attr.Value.Float64()
	case slog.KindBool:
		v = attr.Value.Bool()
	case slog.KindDuration:
		v = attr.Value.Duration()
	case slog.KindTime:
		v = attr.Value.Time()
	case slog.KindAny:
		switch x := attr.Value.Any().(type) {
		case json.Marshaler:
			b, _ := json.Marshal(x)
			v = json.RawMessage(b)
		case fmt.Stringer:
			v = x.String()
		case error:
			v = x.Error()
		default:
			b, err := json.Marshal(x)
			if err != nil {
				panic(fmt.Sprintf("bad kind any: %s", attr.Value.Kind()))
			}
			v = json.RawMessage(b)
		}
	default:
		panic(fmt.Sprintf("bad kind: %s", attr.Value.Kind()))
	}

	attrs[attr.Key] = v
}
