package nsq

import (
	"context"
	"fmt"
	"github.com/firma/framework-common/utils"

	"testing"
	"time"
)

type PlanEventName string
type PlanEvent struct {
	PlanId    int64         `json:"plan_id"`
	Uid       int64         `json:"uid"`
	SignTime  int64         `json:"sign_time"`
	Latitude  float64       `json:"latitude" form:"latitude"`   // 经度
	Longitude float64       `json:"longitude" form:"longitude"` // 纬度
	LeaderUid int64         `json:"leader_uid"`
	EventName PlanEventName `json:"event_name"`
}

type GbGatewayEvent struct {
	DeviceInfo interface{} `json:"device_info"`
	Resource   string      `json:"resource"`
	EventName  int64       `json:"event_name"`
}

func TestPubDevice(t *testing.T) {
	hexString := "404060c401011b2208110a172ad49a3b0000000000000000300002020101052839000d000300fe056f82000400050c120000000823100917000000000000000000000000001b2208110a17092323"
	fmt.Println("hexString", hexString)
	dataRaw, _ := utils.HexToBytes(hexString)

	//deviceInfo, _ := data.RowData(dataRaw, len(dataRaw))

	data := GbGatewayEvent{
		//DeviceInfo: deviceInfo,
		Resource:  utils.HexTo16String(dataRaw),
		EventName: 2,
	}
	m, err := NewProducer(
		&NsqConfig{
			NsqdAddr:       "42.192.91.25:4150",
			NsqLookupdAddr: "42.192.91.25:4161",
			AuthSecret:     "jV22WdmaXxHWAiAh",
			EnableLookup:   false,
		},
	)
	fmt.Println(err)
	ctx := context.Background()
	err = m.Publish(ctx, "GbWorkStatusSet", data)
	fmt.Println(err)
}

func TestPub(t *testing.T) {
	m, err := NewProducer(
		&NsqConfig{
			NsqdAddr:       "42.192.91.25:4150",
			NsqLookupdAddr: "42.192.91.25:4161",
			AuthSecret:     "jV22WdmaXxHWAiAh",
			EnableLookup:   false,
		},
	)
	ctx := context.Background()

	data := PlanEvent{
		PlanId:    936,
		Uid:       227,
		SignTime:  1683339184,
		Latitude:  0,
		Longitude: 0,
		LeaderUid: 0,
		EventName: "plan_end",
	}

	for i := 0; i < 1; i++ {
		err = m.Publish(ctx, "PlanWorkStatusSet", data)
		if err != nil {
			fmt.Println("Message queue error: ", err.Error())
		}
	}

}

func TestDarftPub(t *testing.T) {
	m, err := NewProducer(
		&NsqConfig{
			NsqdAddr:       "127.0.0.1:4150",
			NsqLookupdAddr: "127.0.0.1:4161",

			EnableLookup: false,
		},
	)
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		err = m.DeferredPublish(ctx, "RoomSendMessage", []byte(fmt.Sprintf("order %d", i)), 20*time.Second)
		if err != nil {
			fmt.Println("Message queue error: ", err.Error())
		}
	}

}
