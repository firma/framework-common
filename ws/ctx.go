package ws

import (
	"context"
	"net/http"
	"strings"
)

type (
	headerKey struct {
	}
	queryKey struct {
	}
)

func HeadersFromServerContext(ctx context.Context) http.Header {
	v, ok := ctx.Value(headerKey{}).(http.Header)
	if ok {
		return v
	}

	return http.Header{}
}

func QueryTokenFromServerContext(ctx context.Context) string {
	v, ok := ctx.Value(queryKey{}).(string)
	if ok {

		return strings.Replace(v, "/ws/", "", -1)
	}

	return ""
}
