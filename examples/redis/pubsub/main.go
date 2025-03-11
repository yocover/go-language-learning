package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func subscribe(ctx context.Context, rdb *redis.Client, channel string, done chan bool) {
	// 订阅频道
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	// 等待确认订阅成功
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// 接收消息的通道
	ch := pubsub.Channel()

	fmt.Printf("开始订阅频道: %s\n", channel)

	// 处理消息
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				fmt.Println("频道已关闭")
				return
			}
			fmt.Printf("收到消息: %s\n", msg.Payload)
			if msg.Payload == "quit" {
				done <- true
				return
			}
		case <-ctx.Done():
			fmt.Println("上下文已取消")
			return
		}
	}
}

func publish(ctx context.Context, rdb *redis.Client, channel string, messages []string) {
	for _, msg := range messages {
		// 发布消息
		err := rdb.Publish(ctx, channel, msg).Err()
		if err != nil {
			panic(err)
		}
		fmt.Printf("发布消息: %s\n", msg)
		time.Sleep(time.Second) // 间隔发送
	}
}

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()

	// 创建一个通道用于同步
	done := make(chan bool)

	// 定义要发送的消息
	messages := []string{
		"Hello, Redis!",
		"这是一条测试消息",
		"发布/订阅模式示例",
		"quit", // 结束信号
	}

	// 启动订阅者
	channel := "news"
	go subscribe(ctx, rdb, channel, done)

	// 等待一秒确保订阅者已经准备好
	time.Sleep(time.Second)

	// 启动发布者
	go publish(ctx, rdb, channel, messages)

	// 等待结束信号
	<-done
	fmt.Println("程序结束")
}
