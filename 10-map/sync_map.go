package main

import (
	"fmt"
	"sync"
)

// 使用sync.Map实现并发安全的map
func DescSyncMap() {

	var sm sync.Map
	// 存储值

	sm.Store("name", "wangzhongjie")
	sm.Store("age", 25)
	sm.Store("gender", "male")

	// 获取值
	if value, ok := sm.Load("name"); ok {
		fmt.Printf("name: %v\n", value)
	}

	// 获取或者存储
	// 如果键不存在，则存储值并返回
	value, loaded := sm.LoadOrStore("city", "shanghai")
	fmt.Printf("city: %v, loaded: %v\n", value, loaded)

	// 遍历

	// sync.Map不支持for range遍历，可以使用Range方法遍历
	// for key, value := range sm {
	// }

	sm.Range(func(key, value any) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})

	// // 遍历sync.Map中的所有键值对进行打印
	// sm.Range(func(key, value any) bool {
	// 	fmt.Printf("%v: %v\n", key, value)
	// 	return true
	// })

}
