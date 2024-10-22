package utils

import (
	"fmt"
	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	node   *sf.Node
	node31 *sf.Node
	err    error

	mux sync.Mutex
)

var ID_PRE [9]int = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

// SnowflakeInit 雪花算法工具初始化
// param: 机器节点(用于生成分布式id)
func SnowflakeInit(machineID int64) error {
	node, err = sf.NewNode(machineID)
	if err != nil {
		return fmt.Errorf("snowflake.GenerateID err %v", zap.Error(err))
	}

	node31, _ = sf.NewNode(31)

	return nil
}

// GenerateID 生成分布式id
func GenerateID() int64 {
	return node.Generate().Int64()
}

func GenerateOrderSn() string {
	return node.Generate().String()
}

func GenerateID2() sf.ID {
	return node31.Generate()
}

func GenerateShortUid() string {
	mux.Lock()
	defer mux.Unlock()
	rand.Seed(int64(time.Now().UnixNano()))

	sid := GenerateID2().String()
	// node31, _ = sf.NewNode(31)
	// sid := node31.Generate().String()
	uid := fmt.Sprintf("%v%s", ID_PRE[rand.Intn(9)], sid[len(sid)-9:])
	return uid
}

func GenerateShortId() string {
	mux.Lock()
	defer mux.Unlock()
	rand.Seed(int64(time.Now().UnixNano()))

	now := time.Now()
	timeFormat := now.Format("060102")
	genCode := RandStr(8)
	return fmt.Sprintf("%s-%s", timeFormat, strings.ToUpper(genCode))
}

func GenerateShortIdWithPrefix2(prefix string) string {
	mux.Lock()
	defer mux.Unlock()
	rand.Seed(int64(time.Now().UnixNano()))

	now := time.Now()
	timeFormat := now.Format("060102")
	genCode := RandStr(8)
	return fmt.Sprintf("%s%s-%s", prefix, timeFormat, strings.ToUpper(genCode))
}

func GenerateShortIdWithPrefix(prefix string) string {
	mux.Lock()
	defer mux.Unlock()
	rand.Seed(int64(time.Now().UnixNano()))

	genCode := RandStr(8)
	return fmt.Sprintf("%s-%s", prefix, strings.ToUpper(genCode))
}
