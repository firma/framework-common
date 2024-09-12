package ws

import (
	"context"
	"github.com/olahol/melody"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	sessionKey struct {
	}
)

func SessionFromServerContext(ctx context.Context) *melody.Session {
	v, ok := ctx.Value(sessionKey{}).(*melody.Session)
	if !ok {
		return &melody.Session{}
	}

	return v
}

func CloseSessionFromServerContext(ctx context.Context) {
	s := SessionFromServerContext(ctx)

	if s.IsClosed() {
		return
	}

	if err := s.Close(); err != nil {
		logx.WithContext(ctx).Errorw("关闭会话异常", logx.Field("err", err))
	}
}
