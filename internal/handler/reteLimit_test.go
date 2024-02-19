package handler

import (
	"log"
	"sync"
	"testing"
	"time"
)

func addToQueue(item *RateLimitResource) {
	RateLimitQue.AddRateLimited(item)
}

func TestRateLimitQueue(t *testing.T) {
	// 创建一个等待组以等待所有方法完成
	var wg sync.WaitGroup

	// 要测试的方法数量
	methodCount := 5

	// 每个方法在开始时都向队列中添加一个RateLimitResource
	wg.Add(methodCount)
	for i := 0; i < methodCount; i++ {
		go func(methodIndex int) {
			defer wg.Done()
			resource := NewRateLimitResource("testType", methodIndex)
			addToQueue(resource)
		}(i)
	}

	// 等待所有方法完成向队列的添加
	wg.Wait()

	// 检验队列中的项目处理速率
	itemsProcessed := 0
	startTime := time.Now()
	log.Println("start:", startTime.Unix())

	for itemsProcessed < methodCount {
		item, _ := RateLimitQue.Get()
		// 处理项目，这里仅模拟处理
		RateLimitQue.Done(item)

		itemsProcessed++
		if itemsProcessed < methodCount {
			// 等待一会儿，以便观察下一个项目是否被限流
			time.Sleep(500 * time.Millisecond)
		}
	}

	since := time.Since(startTime)

	log.Println("since", since.String(), time.Now().Unix())
	// 验证是否满足限流条件（处理完所有项目应该至少需要 methodCount 秒）
	if since < (time.Duration(methodCount) * time.Second) {
		t.Errorf("Rate limiting did not work as expected. Items were processed too quickly.")
	}
}
