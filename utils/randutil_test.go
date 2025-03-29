package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 测试并发生成订单号的唯一性
func TestGenerateOrderNumberConcurrency(t *testing.T) {
	fmt.Println(len("1605475281009647616"))
	// 并发数
	numGoroutines := 5000000
	//numGoroutines := 20000

	// 用于存储生成的订单号
	var orderNumbers sync.Map

	// 用于等待所有goroutine完成
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 记录冲突数量
	var conflicts int64

	// 启动多个goroutine同时生成订单号
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			// 生成订单号
			//orderNum := GenerateShortIdWithPrefix("T")
			orderNum := GenerateAbsoluteUniqueOrderNumber(int64(i))
			if numGoroutines == 20000 {
				fmt.Println(orderNum, len(orderNum))
			}

			// 检查是否已存在相同的订单号
			if _, loaded := orderNumbers.LoadOrStore(orderNum, true); loaded {
				// 存在冲突，原子增加计数器
				// atomic.AddInt64(&conflicts, 1)
				t.Logf("冲突: %s, %d", orderNum, len(orderNum))
				conflicts++
			}

			// 模拟一些其他操作，增加并发冲突的可能性
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 计算生成的唯一订单号数量
	uniqueCount := 0
	orderNumbers.Range(func(_, _ interface{}) bool {
		uniqueCount++
		return true
	})

	// 输出结果
	t.Logf("生成 %d 订单号 %d 重复冲突", numGoroutines, conflicts)
	t.Logf("唯一订单号: %d", uniqueCount)

	// 断言没有冲突
	if conflicts > 0 {
		t.Errorf("得到%d 冲突", conflicts)
	}

	// 断言唯一订单号数量等于预期
	if uniqueCount != numGoroutines {
		t.Errorf("预期有 %d 个唯一的订单号，但实际得到 %d ", numGoroutines, uniqueCount)
	}
}
