package main

import (
	"fmt"
	"log"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化 Sentinel 异常：%v", err)
	}

	// 配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "qps-test-1",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 匀速通过
			MaxQueueingTimeMs:      100,             // 匀速排队的最大等待时间，该字段仅仅对 `Throttling` ControlBehavior生效
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "qps-test-2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              10,
			StatIntervalInMs:       5000,
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败：%v", err)
		return
	}

	for i := 0; i < 15; i++ {
		e, b := sentinel.Entry("qps-test-1", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Println("限流了")
		} else {
			fmt.Println("检查通过")
			e.Exit()
		}
		time.Sleep(10 * time.Millisecond)
	}
}
