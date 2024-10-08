//go:build jsoniter

package log

import (
	jsoniter "github.com/json-iterator/go"
)

func Marshal(v any) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return jsoniter.ConfigFastest.MarshalIndent(v, prefix, indent)
}
