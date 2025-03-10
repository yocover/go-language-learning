package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"
)

// ResourcePool 表示一个资源池（如数据库连接池）
type ResourcePool struct {
	sem     *semaphore.Weighted // 信号量控制并发访问
	maxSize int64               // 资源池最大容量
}

// NewResourcePool 创建一个新的资源池
func NewResourcePool(maxSize int64) *ResourcePool {
	return &ResourcePool{
		sem:     semaphore.NewWeighted(maxSize),
		maxSize: maxSize,
	}
}

// SimulateQuery 模拟执行一个数据库查询
func (p *ResourcePool) SimulateQuery(ctx context.Context, queryID int) error {
	// 获取信号量（获取一个资源）
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("无法获取资源: %v", err)
	}

	// 使用 defer 确保释放信号量
	defer p.sem.Release(1)

	// 模拟查询执行时间（随机 100-300ms）
	queryTime := time.Duration(100+rand.Intn(200)) * time.Millisecond

	fmt.Printf("[查询 %d] 开始执行查询，预计耗时: %v\n", queryID, queryTime)
	time.Sleep(queryTime)
	fmt.Printf("[查询 %d] 查询完成\n", queryID)

	return nil
}

func TestSemaphore(t *testing.T) {
	fmt.Println("=== 测试信号量 ===")

	// 创建一个最多允许3个并发连接的资源池
	pool := NewResourcePool(3)

	// 创建上下文
	ctx := context.Background()

	// 使用等待组来等待所有查询完成
	var wg sync.WaitGroup

	// 模拟10个并发查询
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(queryID int) {
			defer wg.Done()

			// 执行查询
			startTime := time.Now()
			err := pool.SimulateQuery(ctx, queryID)
			if err != nil {
				fmt.Printf("[查询 %d] 错误: %v\n", queryID, err)
				return
			}

			// 打印查询耗时
			fmt.Printf("[查询 %d] 总耗时: %v\n", queryID, time.Since(startTime))
		}(i)

		// 短暂延迟，使输出更清晰
		time.Sleep(time.Millisecond * 100)
	}

	// 等待所有查询完成
	wg.Wait()
	fmt.Println("\n所有查询已完成")
}

// ExampleResourcePool_MultipleQueries 展示多个查询如何使用资源池
func ExampleResourcePool_MultipleQueries() {
	// 创建一个容量为2的资源池
	pool := NewResourcePool(2)
	ctx := context.Background()

	// 创建3个查询
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := pool.SimulateQuery(ctx, id)
			if err != nil {
				fmt.Printf("查询 %d 失败: %v\n", id, err)
			}
		}(i)
	}

	wg.Wait()
	// Output:
}

// BenchmarkResourcePool 基准测试资源池性能
func BenchmarkResourcePool(b *testing.B) {
	pool := NewResourcePool(5)
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		queryID := 0
		for pb.Next() {
			queryID++
			err := pool.SimulateQuery(ctx, queryID)
			if err != nil {
				b.Errorf("查询失败: %v", err)
			}
		}
	})
}
