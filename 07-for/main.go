package main

import "fmt"

func main() {

	var count int = 10

	for i := range count {
		fmt.Println(i)
	}

	fmt.Println("------------")
	var list = []int{1, 2, 3, 4, 5}

	for i := range list {
		fmt.Println("list[", i, "]:", list[i])
	}

	for i := range list {
		println("list[", i, "]:", list[i])
	}

	fmt.Println("------------")

	num := 100
	for num < 1000 {
		num *= 2
		fmt.Println(num)
	}

	// for range 可以用来遍历数组、切片、map、字符串、通道等。

	// 遍历数组
	arr := [5]int{1, 2, 3, 4, 5}
	for i := range arr {
		print("arr[", i, "]:", arr[i], " ")
	}

	println("")
	// 遍历切片
	slice := []int{1, 2, 3, 4, 5}
	for i := range slice {
		print("slice[", i, "]:", slice[i], " ")
	}

	// 遍历map
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	var map2 = map[string]string{"a": "apple", "b": "banana", "c": "cherry"}
	for k, v := range m {
		fmt.Printf("key:%s,value:%d\n", k, v)
	}
	for k, v := range map2 {
		fmt.Printf("key:%s,value:%s\n", k, v)
	}

	// 遍历字符串

}

// Output:
// 0
// 1
// 2
