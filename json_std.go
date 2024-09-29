//go:build !jsoniter

package log

import "encoding/json"

type (
	Marshaler  = json.Marshaler
	RawMessage = json.RawMessage
)

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
