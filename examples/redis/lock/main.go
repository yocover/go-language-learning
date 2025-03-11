package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DistributedLock Redis 分布式锁实现
type DistributedLock struct {
	rdb        *redis.Client
	key        string
	value      string
	expiration time.Duration
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(rdb *redis.Client, key string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		rdb:        rdb,
		key:        fmt.Sprintf("lock:%s", key),
		value:      fmt.Sprintf("lock:%d", time.Now().UnixNano()),
		expiration: expiration,
	}
}

// TryLock 尝试获取锁
func (l *DistributedLock) TryLock(ctx context.Context) (bool, error) {
	return l.rdb.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

// Unlock 释放锁
func (l *DistributedLock) Unlock(ctx context.Context) error {
	// 使用 Lua 脚本确保原子性操作
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end`

	result, err := l.rdb.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return fmt.Errorf("锁已过期或被其他进程释放")
	}

	return nil
}

// 模拟业务处理
func processOrder(orderID string) error {
	fmt.Printf("处理订单 %s...\n", orderID)
	// 模拟处理时间
	time.Sleep(2 * time.Second)
	return nil
}

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()

	// 创建分布式锁
	orderID := "12345"
	lock := NewDistributedLock(rdb, fmt.Sprintf("order:%s", orderID), 10*time.Second)

	// 尝试获取锁
	fmt.Printf("尝试获取订单 %s 的处理锁...\n", orderID)
	acquired, err := lock.TryLock(ctx)
	if err != nil {
		panic(err)
	}

	if !acquired {
		fmt.Println("无法获取锁，订单正在被其他进程处理")
		return
	}

	fmt.Println("成功获取锁，开始处理订单")

	// 使用 defer 确保在函数结束时释放锁
	defer func() {
		if err := lock.Unlock(ctx); err != nil {
			fmt.Printf("释放锁时发生错误: %v\n", err)
		} else {
			fmt.Println("锁已释放")
		}
	}()

	// 处理订单
	if err := processOrder(orderID); err != nil {
		fmt.Printf("处理订单时发生错误: %v\n", err)
		return
	}

	fmt.Println("订单处理完成")

	// 演示并发场景
	fmt.Println("\n模拟并发请求...")
	// 创建另一个锁实例（模拟另一个进程）
	lock2 := NewDistributedLock(rdb, fmt.Sprintf("order:%s", orderID), 10*time.Second)

	acquired2, err := lock2.TryLock(ctx)
	if err != nil {
		panic(err)
	}

	if !acquired2 {
		fmt.Println("（并发请求）无法获取锁，订单正在被其他进程处理")
	} else {
		fmt.Println("警告：锁被重复获取，这不应该发生！")
		lock2.Unlock(ctx)
	}
}
