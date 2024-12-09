package rabbitmq

import (
	"context"
	"github.com/firma/framework-common/queue"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rabbitmq/amqp091-go"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-queue/rabbitmq"
	mq "github.com/zeromicro/go-zero/core/queue"
	"github.com/zeromicro/go-zero/core/threading"
)

type handler struct {
	fn queue.MessageHandler
}

func (h handler) Consume(message string) error {
	return h.fn(context.TODO(), []byte(message))
}

var _ queue.Queue = (*hub)(nil)

type hub struct {
	cfg rabbitmq.RabbitConf

	subscribers map[string]queue.ISubscribe

	consumers map[string]mq.MessageQueue

	admin *rabbitmq.Admin
}

func NewQueue(cfg rabbitmq.RabbitConf) queue.Queue {
	return &hub{
		cfg:         cfg,
		subscribers: make(map[string]queue.ISubscribe),
		consumers:   make(map[string]mq.MessageQueue),
	}
}

func (q *hub) RegisterSubscribe(subscribe queue.ISubscribe) {
	q.subscribers[subscribe.TopicName()] = subscribe
}

func (q *hub) Start() error {
	q.admin = rabbitmq.MustNewAdmin(q.cfg)

	log.Infow("开启队列")

	for _, v := range q.subscribers {

		//if err := q.admin.DeclareExchange(
		//	rabbitmq.ExchangeConf{
		//		ExchangeName: v.TopicName(),
		//		Type:         "x-delayed-message",
		//		Durable:      true,
		//		AutoDelete:   false,
		//		Internal:     false,
		//		NoWait:       false,
		//	},
		//	amqp091.Table(
		//		amqp.Table{
		//			"x-delayed-type": v.Type(), // 这个可能需要安装插件
		//		},
		//	),
		//); err != nil {
		//	return err
		//}

		if err := q.admin.DeclareQueue(
			rabbitmq.QueueConf{
				Name:       v.TopicName(),
				Durable:    true,
				AutoDelete: false,
				Exclusive:  false,
				NoWait:     false,
			}, amqp091.Table(amqp.Table{}),
		); err != nil {
			return err
		}
		if err := q.admin.Bind(
			v.TopicName(), "", v.TopicName(), false, amqp091.Table(amqp.Table{}),
		); err != nil {
			return err
		}

		s := rabbitmq.MustNewListener(
			rabbitmq.RabbitListenerConf{
				RabbitConf: q.cfg,
				ListenerQueues: []rabbitmq.ConsumerConf{
					{
						Name:      v.TopicName(),
						AutoAck:   true,
						Exclusive: false,
						NoLocal:   false,
						NoWait:    false,
					},
				},
			}, handler{
				fn: v.Handler(),
			},
		)

		log.Infow("注册队列", "name", v.TopicName())

		q.consumers[v.TopicName()+":"+v.Channel()] = s

		threading.GoSafe(s.Start)
	}

	return nil
}

func (q *hub) Stop() error {
	for _, v := range q.consumers {
		v.Stop()
	}

	log.Infow("结束队列")

	return nil
}
