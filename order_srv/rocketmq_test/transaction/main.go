package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct{}

func (listener *OrderListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(10 * time.Second)
	fmt.Println("本地逻辑执行失败")
	return primitive.UnknowState
}

func (listener *OrderListener) CheckLocalTransaction(message *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("RocketMQ的消息回查")
	time.Sleep(12 * time.Second)
	fmt.Println("RocketMQ的回查成功")
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(&OrderListener{}, producer.WithNameServer([]string{"10.120.21.77:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	res, err := p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage("TransTopic", []byte("this is a transaction message.")),
	)
	if err != nil {
		fmt.Printf("发送失败：%s\n", err)
	} else {
		fmt.Printf("发送成功：result=%s\n", res.String())
	}

	time.Sleep(30 * time.Minute)
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败: %s")
	}
}
