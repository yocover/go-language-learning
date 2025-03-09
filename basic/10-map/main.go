package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

func main() {

	// map 映射表的使用方式

	// createMap()
	// accessMap()
	// lenMap()
	// replaceMapValue()
	// deleteMapValue()
	traverseMap()

	// 创建并发安全的map
	sm := NewSafeMap()

	// 等待创建组
	var wg sync.WaitGroup

	// 启动多个gouroutine 并发写入

	for i := range 10 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sm.Set(fmt.Sprintf("key-%d", n), i)
		}(i)
	}

	// 等待所有gouroutine 结束
	wg.Wait()

	// 读取map

	for i := range 10 {

		if val, exists := sm.Get(fmt.Sprintf("key-%d", i)); exists {
			fmt.Printf("key-%d value is %d\n", i, val)
		}
	}

	fmt.Printf("-----------\n")
	DescSyncMap()
}

func createMap() {
	// 创建一个空的映射表
	mp := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	mp2 := map[string]int{
		"name": 1,
		"age":  2,
		"city": 3,
	}

	var m3 = make(map[string]int)
	m3["name"] = 1
	m3["age"] = 2
	// m3["age"] = 2
	// m3["city"] = 3

	// 打印映射表
	fmt.Printf("mp: %v\n", mp)
	for k, v := range mp {
		println(k, v)
	}

	fmt.Printf("mp2: %v\n", mp2)
	for k, v := range mp2 {
		println(k, v)
	}

	fmt.Printf("m3: %v\n", m3)
	for k, v := range m3 {
		println(k, v)
	}

	/**
	* 注意：
	map在使用前必须初始化
	可以使用make 函数初始化
	可以使用字面量语法初始化
	直接声明的map 时 nil map 不能直接使用
	*/
	// var m4 map[string]int
	// m4["a"] = 1
}

// 访问map
func accessMap() {

	var mp = map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	mp[1] = "four" // 修改元素

	// 通过索引访问
	println(mp[1])
	println(mp[2])

	// 对于不存在的键，返回值将是零值
	println(mp[21])

	/**

	这段代码返回两个值

	val, exist := mp[21]
		val 时健对应的值，如果不存在，则为该类型的零值
		exist 时一个布尔值，表示键是否存在，true表示存在 false表示不存在
	*/

	// 这种写法比较好，可以明确知道键是否存在
	if val, exist := mp[21]; exist {
		fmt.Println("key exist, value is ", val)
	} else {
		fmt.Println("key not exist")
	}
}

// 求map 的长度
func lenMap() {
	var mp = map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	fmt.Printf("map: %v\n", mp)
	println(len(mp))
}

func replaceMapValue() {
	// 不能直接修改map的值，需要先获取到值，然后重新赋值
	var mp = map[string]string{
		"name":  "wangzhongjie",
		"email": "123456@qq.com",
	}

	if _, exist := mp["name"]; exist {
		mp["name"] = "wangzhongjie-1"
	} else {
		mp["name"] = "Alex"
	}

	fmt.Printf("map: %v\n", mp)

	mathMp := make(map[float64]string, 10)
	fmt.Printf("mathMp: %v\n", mathMp)
	fmt.Printf("mathMp length: %d\n", len(mathMp))

	for index := range 10 {
		mathMp[math.NaN()] = strconv.Itoa(index) // 或者使用 fmt.Sprintf("%d", index)
		// / 直接打印整个map，而不是尝试访问特定键
		fmt.Printf("当前map内容: %v，长度: %d\n", mathMp, len(mathMp))
	}

	fmt.Printf("mathMp length: %d\n", len(mathMp))
	fmt.Printf("mathMp : %v\n", mathMp)

}

// 删除map中的元素
// 删除map中的元素，需要使用delete()函数
func deleteMapValue() {
	var mp = map[string]string{
		"name":  "wangzhongjie",
		"email": "123456@qq.com",
	}

	fmt.Printf("map: %v\n", mp)
	delete(mp, "name")
	fmt.Printf("map: %v\n", mp)
}

// map的遍历
// 遍历map，可以使用range关键字，遍历的结果是键值对
// 遍历的结果是无序的
func traverseMap() {

	var mp = map[string]string{
		"name":  "wangzhongjie",
		"email": "123456@qq.com",
		"age":   "25",
		"1":     "one",
		"2":     "two",
		"3":     "three",
	}

	for k, v := range mp {
		println(k, v)
	}
}
