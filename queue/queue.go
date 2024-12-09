package queue

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"

	"time"
)

type (
	Logger interface {
		log.Logger
	}
	Serialize interface {
		Marshal() ([]byte, error)
	}
	CodeType int64
)

type IProducer interface {
	Publish(ctx context.Context, topic string, intf interface{}) error
	DeferredPublish(ctx context.Context, topic string, intf interface{}, delay time.Duration) error
	Stop()
}

type Queue interface {
	RegisterSubscribe(subscribe ISubscribe)
	Start() error
	Stop() error
}

type MessageHandler func(ctx context.Context, data []byte) error

type ISubscribe interface {
	TopicName() string
	Channel() string
	Handler() MessageHandler
	Type() string
}

type defaultSubscribe struct {
	topicName   string
	h           MessageHandler
	messageType string // rabbitmq里会用到：direct|fanout|topic|headers
}

func (d defaultSubscribe) TopicName() string {
	return d.topicName
}

func (d defaultSubscribe) Channel() string {
	return "default"
}

func (d defaultSubscribe) Handler() MessageHandler {
	return d.h
}
func (d defaultSubscribe) Type() string {
	return "direct"
}

func NewDefaultSubscribe(topicName string, handler MessageHandler) ISubscribe {
	return defaultSubscribe{
		topicName: topicName,
		h:         handler,
	}
}

func NewTypeSubscribe(topicName string, mt string, handler MessageHandler) ISubscribe {
	return defaultSubscribe{
		topicName:   topicName,
		h:           handler,
		messageType: mt,
	}
}
