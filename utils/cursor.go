package utils

import (
	"encoding/base64"
	"encoding/json"
)

func UnmarshalCursor[T any](s string) *T {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return new(T)
	}

	var cursor T
	if err = json.Unmarshal(res, &cursor); err != nil {
		return new(T)
	}

	return &cursor
}

func MarshalCursor[T any](t T) string {
	res, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(res)
}
