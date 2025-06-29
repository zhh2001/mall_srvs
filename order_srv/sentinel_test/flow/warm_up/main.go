package main

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/core/base"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
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
			Resource:               "warmup-test",
			TokenCalculateStrategy: flow.WarmUp, // 预热/冷启动
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              1000,        // 表示流控阈值
			WarmUpPeriodSec:        30,          // 预热的时间长度，该字段仅仅对 `WarmUp` 的TokenCalculateStrategy生效；
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败：%v", err)
		return
	}

	var globalTotal int
	var blockTotal int
	var passTotal int

	ch := make(chan struct{})

	for i := 0; i < 128; i++ {
		go func() {
			for {
				globalTotal = globalTotal + 1
				e, b := sentinel.Entry("warmup-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					blockTotal = blockTotal + 1
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal = passTotal + 1
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	go func() {
		var oldGlobalTotal int
		var oldBlockTotal int
		var oldPassTotal int

		for {
			oneSecondGlobalTotal := globalTotal - oldGlobalTotal
			oldGlobalTotal = globalTotal
			oneSecondBlockTotal := blockTotal - oldBlockTotal
			oldBlockTotal = blockTotal
			oneSecondPassTotal := passTotal - oldPassTotal
			oldPassTotal = passTotal

			time.Sleep(1 * time.Second)
			fmt.Printf("globalTotal = %d, passTotal = %d, blockTotal = %d\n",
				oneSecondGlobalTotal, oneSecondPassTotal, oneSecondBlockTotal)
		}
	}()

	<-ch
}
