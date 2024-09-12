package nsq

import (
	"context"
	"github.com/firma/framework-common/queue"
	"github.com/nsqio/go-nsq"
	"github.com/zeromicro/go-zero/core/logx"
)

// 1.defaultHandler 是一个消费者类型
type defaultHandler struct {
	fn queue.MessageHandler
}

func newDefaultHandler(fn queue.MessageHandler) nsq.Handler {
	return defaultHandler{
		fn: fn,
	}
}

func (h defaultHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	if err := h.fn(context.TODO(), m.Body); err != nil {
		return err
	}

	return nil
}

var _ queue.Queue = (*NSQQueue)(nil)

type NSQQueue struct {
	cfg         NsqConfig
	subscribers map[string]queue.ISubscribe
	consumers   map[string]*nsq.Consumer
}

func NewQueue(config NsqConfig) queue.Queue {
	return NSQQueue{
		cfg:         config,
		subscribers: make(map[string]queue.ISubscribe, 0),
		consumers:   make(map[string]*nsq.Consumer, 0),
	}
}

func (n NSQQueue) RegisterSubscribe(subscribe queue.ISubscribe) {
	n.subscribers[subscribe.TopicName()] = subscribe
}

func (n NSQQueue) Start() error {
	c := nsq.NewConfig()
	if len(n.cfg.AuthSecret) > 0 {
		c.AuthSecret = n.cfg.AuthSecret
	}

	for _, v := range n.subscribers {
		consumer, err := nsq.NewConsumer(v.TopicName(), v.Channel(), c)
		if err != nil {
			return err
		}
		logx.Infow("注册消费队列", logx.Field("topic", v.TopicName()), logx.Field("channel", v.Channel()))

		consumer.AddHandler(
			newDefaultHandler(v.Handler()),
		)

		if n.cfg.EnableLookup {
			if e := consumer.ConnectToNSQLookupd(n.cfg.NsqLookupdAddr); e != nil {
				return e
			}
		} else {
			if e := consumer.ConnectToNSQD(n.cfg.NsqdAddr); e != nil {
				return e
			}
		}
		n.consumers[v.TopicName()+":"+v.Channel()] = consumer
	}

	return nil
}

func (n NSQQueue) Stop() error {
	for _, v := range n.consumers {
		v.Stop()
	}

	return nil
}
