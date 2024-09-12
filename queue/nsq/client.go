package nsq

import "github.com/firma/framework-common/queue"

var Producer queue.IProducer

func InitProducer(con NsqConfig) {
	producer, err := NewProducer(&con)
	if err != nil {
		panic("Producer初始化失败")
	}
	Producer = producer
}

func Queue() queue.IProducer {
	return Producer
}
