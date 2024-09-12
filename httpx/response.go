package httpx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/firma/framework-common/errno"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type (
	Response      errno.Error
	ContextFunc   func(ctx context.Context) Response
	contextKey    struct{}
	ginContextKey struct{}
)

var (
	ginCtxNotFoundErr = fmt.Errorf("nil gin ctx")
)

func Json(e ContextFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()

		ctx = context.WithValue(ctx, contextKey{}, ginCtx.Request)

		ctx = context.WithValue(ctx, ginContextKey{}, ginCtx)

		f := e(ctx)

		f.RenderJson(ginCtx)
	}
}

func RequestFromServerContext(ctx context.Context) *http.Request {
	val := ctx.Value(contextKey{})

	request, ok := val.(*http.Request)
	if !ok {
		return &http.Request{
			Body: io.NopCloser(bytes.NewBuffer([]byte{})),
		}
	}

	return request
}

func RequestHeaderFromServerContext(ctx context.Context) http.Header {
	req := RequestFromServerContext(ctx)

	return req.Header
}

func RawDataFromServerContext(ctx context.Context) ([]byte, error) {
	req := RequestFromServerContext(ctx)

	return io.ReadAll(req.Body)
}

func GinCtxFromServerContext(ctx context.Context) *gin.Context {
	val := ctx.Value(ginContextKey{})

	ginCtx, ok := val.(*gin.Context)
	if ok {
		return ginCtx
	}

	return nil
}

func ShouldBindURIFromServerContext(ctx context.Context, obj interface{}) error {
	ginCtx := GinCtxFromServerContext(ctx)
	if ginCtx == nil {
		return ginCtxNotFoundErr
	}

	return ginCtx.ShouldBindUri(obj)
}

func ShouldBindFromServerContext(ctx context.Context, obj interface{}) error {
	ginCtx := GinCtxFromServerContext(ctx)
	if ginCtx == nil {
		return ginCtxNotFoundErr
	}

	return ginCtx.ShouldBind(obj)
}
