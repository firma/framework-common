package nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/firma/framework-common/queue"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
	"time"
)

var _ queue.IProducer = (*nsqProducer)(nil)

type nsqProducer struct {
	producer *nsq.Producer
}

func NewProducer(config *NsqConfig) (queue.IProducer, error) {
	c := nsq.NewConfig()
	if len(config.AuthSecret) > 0 {
		c.AuthSecret = config.AuthSecret
	}
	producer, err := nsq.NewProducer(config.NsqdAddr, c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize nsq producer: %v", err)
	}

	return &nsqProducer{producer: producer}, nil
}

func (p *nsqProducer) Publish(ctx context.Context, topic string, intf interface{}) error {
	data, err := json.Marshal(intf)
	if err != nil {
		return err
	}
	if err := p.producer.Publish(topic, data); err != nil {
		log.Errorw(
			"NsqPublish", "err", err, "topic", topic, "topic", string(data),
		)
		return err
	} else {
		return nil
	}
}

func (p *nsqProducer) DeferredPublish(ctx context.Context, topic string, intf interface{}, delay time.Duration) error {

	data, err := json.Marshal(intf)
	if err != nil {
		return err
	}
	if err := p.producer.DeferredPublish(topic, delay, data); err != nil {
		log.Errorw(
			"NsqDeferredPub", "err", err, "topic", topic, "topic", string(data),
		)
		return err
	} else {
		return nil
	}
}

func (p *nsqProducer) Stop() {
	p.producer.Stop()
}
