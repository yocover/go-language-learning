package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 使用互斥锁的数据结构
type MutexCounter struct {
	mu    sync.Mutex
	count int
}

// 使用读写锁的数据结构
type RWMutexCounter struct {
	mu    sync.RWMutex
	count int
}

// 基准测试：互斥锁，读写比例 9:1（读多写少）
func BenchmarkMutex_ReadMostly(b *testing.B) {
	counter := &MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Intn(10) == 0 { // 10% 的写操作
				counter.mu.Lock()
				counter.count++
				counter.mu.Unlock()
			} else { // 90% 的读操作
				counter.mu.Lock()
				_ = counter.count
				counter.mu.Unlock()
			}
		}
	})
}

// 基准测试：读写锁，读写比例 9:1（读多写少）
func BenchmarkRWMutex_ReadMostly(b *testing.B) {
	counter := &RWMutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Intn(10) == 0 { // 10% 的写操作
				counter.mu.Lock()
				counter.count++
				counter.mu.Unlock()
			} else { // 90% 的读操作
				counter.mu.RLock()
				_ = counter.count
				counter.mu.RUnlock()
			}
		}
	})
}

// 基准测试：互斥锁，读写比例 1:1（读写均衡）
func BenchmarkMutex_EqualReadWrite(b *testing.B) {
	counter := &MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Intn(2) == 0 { // 50% 的写操作
				counter.mu.Lock()
				counter.count++
				counter.mu.Unlock()
			} else { // 50% 的读操作
				counter.mu.Lock()
				_ = counter.count
				counter.mu.Unlock()
			}
		}
	})
}

// 基准测试：读写锁，读写比例 1:1（读写均衡）
func BenchmarkRWMutex_EqualReadWrite(b *testing.B) {
	counter := &RWMutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Intn(2) == 0 { // 50% 的写操作
				counter.mu.Lock()
				counter.count++
				counter.mu.Unlock()
			} else { // 50% 的读操作
				counter.mu.RLock()
				_ = counter.count
				counter.mu.RUnlock()
			}
		}
	})
}
