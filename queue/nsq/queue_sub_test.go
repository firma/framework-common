package nsq

import (
	"context"
	"fmt"
	"github.com/firma/framework-common/queue"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestSub(t *testing.T) {

	q := NewQueue(
		NsqConfig{
			NsqdAddr:       "127.0.0.1:4150",
			NsqLookupdAddr: "127.0.0.1:4161",
			EnableLookup:   false,
		},
	)

	q.RegisterSubscribe(
		queue.NewDefaultSubscribe(
			"PlanWorkStatusSet", func(ctx context.Context, resp []byte) error {
				fmt.Println("test 1: ", string(resp))
				return nil
			},
		),
	)
	q.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	os.Exit(0)
}
