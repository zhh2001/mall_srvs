package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testMall"),
		consumer.WithNameServer([]string{"10.120.21.77:9876"}),
	)

	if err := c.Subscribe("TransTopic", consumer.MessageSelector{},
		func(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range messages {
				fmt.Printf("获取到值：%v\n", messages[i])
			}
			return consumer.ConsumeSuccess, nil
		},
	); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	time.Sleep(10 * time.Minute)
	_ = c.Shutdown()
}
