package main

import (
	"context"
	"fmt"
	"time"
)

func UseContextTest() {

	// timeoutExample()
	// cancelTimeoutExample()
	valueExample()
}

type ContextKey string

const (
	userKey ContextKey = "user"
	ageKey  ContextKey = "age"
)

// 携带值的示例
func valueExample() {

	ctx := context.WithValue(context.Background(), userKey, "张三")
	ctx = context.WithValue(ctx, ageKey, 25)

	// 模拟传递context
	hanleRequest(ctx)
}

func hanleRequest(ctx context.Context) {
	// user := ctx.Value("user").(string)
	// age := ctx.Value("age").(int)
	// fmt.Printf("用户 %s，年龄 %d\n", user, age)

	if user, ok := ctx.Value(userKey).(string); ok {
		if age, ok := ctx.Value(ageKey).(int); ok {
			fmt.Printf("用户 %s，年龄 %d\n", user, age)
		}
	}

}

// 基础超时控制
func timeoutExample() {
	// 创建一个2秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		time.Sleep(2 * time.Second) // 模拟耗时操作
		fmt.Println("timeoutExample done")
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Select: timeoutExample timeout")
	case <-time.After(3 * time.Second):
		fmt.Println("Select: timeoutExample done")
	}
}

// 取消超时操作
func cancelTimeoutExample() {
	ctx, cancel := context.WithCancel(context.Background())

	// 启动工作协程
	go worker(ctx, "工作协程1")
	go worker(ctx, "工作协程2")

	// 让工作协程运行一会
	time.Sleep(time.Second)

	// 取消所有工作
	cancel()

	time.Sleep(time.Second)
}

func worker(ctx context.Context, name string) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到取消信号，退出", name)
			return

		default:
			fmt.Printf(" %s 工作中...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
