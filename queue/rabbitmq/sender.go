package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-queue/rabbitmq"
	"time"
)

type Sender interface {
	Send(exchange string, routeKey string, msg []byte) error
	SendDelay(exchange string, routeKey string, msg []byte, delay time.Duration) error
	Stop() error
}
type RabbitMqSender struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	ContentType string
}

func (q *RabbitMqSender) Send(exchange string, routeKey string, msg []byte) error {
	return q.channel.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: q.ContentType,
			Body:        msg,
		},
	)
}

func (q *RabbitMqSender) SendDelay(exchange string, routeKey string, msg []byte, delay time.Duration) error {
	return q.channel.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: q.ContentType,
			Body:        msg,
			Headers: amqp.Table{
				"x-delay": delay.Milliseconds(),
			},
		},
	)
}

func (q *RabbitMqSender) Stop() error {
	q.conn.Close()
	q.channel.Close()

	return nil
}

func getRabbitURL(rabbitConf rabbitmq.RabbitConf) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s", rabbitConf.Username, rabbitConf.Password,
		rabbitConf.Host, rabbitConf.Port, rabbitConf.VHost,
	)
}
