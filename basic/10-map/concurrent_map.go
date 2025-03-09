package main

import "sync"

// 定义一个并发安全的map结构
type SafeMap struct {
	sync.Mutex
	data map[string]int
}

// 创建新的 SafeMap

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

// 设置值
func (m *SafeMap) Set(key string, value int) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = value
}

// 获取值
func (m *SafeMap) Get(key string) (int, bool) {
	m.Lock()
	defer m.Unlock()
	value, exist := m.data[key]
	return value, exist
}

// 删除值
func (m *SafeMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
}
