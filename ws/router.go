package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/firma/framework-common/httpx"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/olahol/melody"
)

type (
	WebsocketUnaryRouter struct {
		connectRoutes    []RouteNoDataFunc
		routeMap         map[int64]RouteDataFunc
		disconnectRoutes []RouteNoDataFunc

		ConnectionManager *ConnectionManager
	}

	RouteDataFunc func(ctx context.Context, data []byte) MetaData

	RouteNoDataFunc func(ctx context.Context) (next bool)
	RouteFunc       func(ctx context.Context, data []byte) (resp httpx.Response, next bool)

	ConnectionManager struct {
		connections map[int64]*melody.Session
		Register    chan *UserSession
		Unregister  chan *UserSession
	}

	UserSession struct {
		UserID  int64
		OrgId   int64
		Session *melody.Session
	}
)

func NewRouter() *WebsocketUnaryRouter {
	return &WebsocketUnaryRouter{
		connectRoutes:    make([]RouteNoDataFunc, 0),
		routeMap:         map[int64]RouteDataFunc{},
		disconnectRoutes: make([]RouteNoDataFunc, 0),

		ConnectionManager: &ConnectionManager{
			connections: make(map[int64]*melody.Session),
			Register:    make(chan *UserSession),
			Unregister:  make(chan *UserSession),
		},
	}
}

func (r *WebsocketUnaryRouter) AddConnect(fn ...RouteNoDataFunc) {
	r.connectRoutes = append(r.connectRoutes, fn...)
}

func (r *WebsocketUnaryRouter) Add(code int64, fn ...RouteFunc) {
	if _, has := r.routeMap[code]; has {
		panic(fmt.Sprintf("route code %d registered", code))
	}

	r.routeMap[code] = func(ctx context.Context, data []byte) MetaData {
		d := MetaData{
			Code: code,
		}
		for _, f := range fn {
			resp, isNext := f(ctx, data)
			if !isNext {
				d.Data = resp.ToDataString()

				break
			}
		}

		return d
	}
}

func (r *WebsocketUnaryRouter) AddDisconnect(fn ...RouteNoDataFunc) {
	r.disconnectRoutes = append(r.disconnectRoutes, fn...)
}

func (r *WebsocketUnaryRouter) Match(code int64) (RouteDataFunc, error) {
	if _, has := r.routeMap[code]; !has {
		return nil, errors.New("mismatch route")
	}

	return r.routeMap[code], nil
}

func (r *WebsocketUnaryRouter) ConnectHandler(session *melody.Session) {
	ctx := session.Request.Context()

	ctx = context.WithValue(ctx, sessionKey{}, session)

	ctx = context.WithValue(ctx, headerKey{}, session.Request.Header)
	ctx = context.WithValue(ctx, queryKey{}, session.Request.RequestURI)

	for _, fn := range r.connectRoutes {
		isNext := fn(ctx)
		if !isNext {
			if err := session.Close(); err != nil {
				log.Error(err)
			}
		}
	}

	session.Request = session.Request.WithContext(ctx)
}

func (r *WebsocketUnaryRouter) MessageHandler(session *melody.Session, msg []byte) {
	ctx := session.Request.Context()

	// 约定的心跳通讯内容
	if bytes.Compare([]byte("\r\n"), msg) == 0 {
		return
	}

	var data MetaData
	if err := json.Unmarshal(msg, &data); err != nil {
		if ce := session.Write([]byte("unexpected input")); ce != nil {
			log.Error(ce)
		}

		return
	}

	rn, err := r.Match(data.Code)
	if err != nil {
		if ce := session.Write([]byte(err.Error())); ce != nil {
			log.Error(ce)
		}

		return
	}

	resp := rn(ctx, []byte(data.Data))

	res, err := json.Marshal(resp)
	if err != nil {
		if ce := session.Write([]byte(err.Error())); ce != nil {
			log.Error(ce)
		}

		return
	}

	if ce := session.Write(res); ce != nil {
		log.Error(ce)
	}
}

func (r *WebsocketUnaryRouter) DisconnectHandler(session *melody.Session) {
	ctx := session.Request.Context()

	for _, fn := range r.disconnectRoutes {
		isNext := fn(ctx)
		if !isNext {
			if err := session.Close(); err != nil {
				log.Error(err)
			}
		}
	}
}

func (cm *ConnectionManager) Run() {
	for {
		select {
		case v := <-cm.Register:
			cm.connections[v.UserID] = v.Session
			log.Debugf("维系的用户连接 connections %v", cm.connections)

		case s := <-cm.Unregister:
			if _, has := cm.connections[s.UserID]; has {
				delete(cm.connections, s.UserID)
			}
			log.Debugw("维系的用户连接connections %v", cm.connections)
		}
	}
}

func (cm *ConnectionManager) Connections() map[int64]*melody.Session {
	return cm.connections
}
