package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/firma/framework-common/queue"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"time"
)

type publisher struct {
	sender *RabbitMqSender
}

func MustNewSender(rabbitMqConf rabbitmq.RabbitSenderConf) queue.IProducer {
	sender := &RabbitMqSender{ContentType: rabbitMqConf.ContentType}
	conn, err := amqp.Dial(getRabbitURL(rabbitMqConf.RabbitConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %v", err)
	}

	sender.conn = conn
	channel, err := sender.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel, error: %v", err)
	}

	sender.channel = channel
	return &publisher{
		sender: sender,
	}
}

func (q *publisher) Publish(ctx context.Context, topic string, data any) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	logx.WithContext(ctx).Debugw("发布消息", logx.Field("topic", topic), logx.Field("data", data))

	return q.sender.Send(topic, "", res)
}

func (q *publisher) DeferredPublish(ctx context.Context, topic string, data any, delay time.Duration) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	logx.WithContext(ctx).Debugw("发布延迟消息", logx.Field("topic", topic), logx.Field("data", data))
	return q.sender.SendDelay(topic, "", res, delay)
}

func (q *publisher) Stop() {
	q.sender.Stop()
}
