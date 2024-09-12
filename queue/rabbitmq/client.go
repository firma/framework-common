package rabbitmq

import (
	"github.com/firma/framework-common/queue"
	"github.com/firma/framework-common/queue/nsq"
	"github.com/firma/framework-common/stores/gormx"
	"github.com/firma/framework-common/stores/redisx"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

var Producer queue.IProducer

type Config struct {
	rest.RestConf
	WsPort     int
	LogConf    logx.LogConf
	NsqConfig  nsq.NsqConfig
	Rabbitmq   rabbitmq.RabbitConf
	DB         gormx.Config
	Redis      redisx.Config
	UserRpc    zrpc.RpcClientConf
	ProjectRpc zrpc.RpcClientConf
}

func (conf *Config) InitConfig() {
	logx.MustSetup(conf.LogConf)
}

func MustSetup(config Config) {
	Producer = MustNewSender(
		rabbitmq.RabbitSenderConf{
			RabbitConf:  config.Rabbitmq,
			ContentType: "application/json",
		},
	)
}

func Queue() queue.IProducer {
	return Producer
}
