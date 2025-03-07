package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {

	// 不推荐的写法： 先声明 赋值
	// var sum func(int, int) int
	// sum = func(a, b int) int {
	// 	return a + b
	// }

	// 推荐写法 直接使用函数字面量进行赋值
	sum2 := func(a, b int) int {
		return a + b
	}
	fmt.Println(sum2(1, 2))

	// 调用函数
	sayHello()

	numList := [5]int{1, 2, 3, 4, 5}
	fmt.Println(max(numList[:]...))

	// 调用函数并返回值
	result, err := divFunc(0, 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	res := sumFunc(1, 2)
	fmt.Println(res)

	// 匿名函数
	// 匿名函数只能在函数内部存在

	func(a, b int) int {
		fmt.Println("anonymous function", a, b)
		return a + b
	}(1, 2)

	// 闭包的使用
	c := counter()                  // 定义一个闭包
	fmt.Println("count mian:", c()) // 调用闭包
	fmt.Println("count mian:", c()) // 调用闭包
	fmt.Println("count mian:", c()) // 调用闭包
}

// 函数闭包
func counter() func() int {
	count := 0 // 外部变量
	return func() int {
		count++ // 内部函数可以访问外部变量
		fmt.Println("count:", count)
		return count
	}
}

func sumFunc(a, b int) (count int) {
	return a + b
}

// 对于一个类型相同的参数而言，可以只声明一次类型
func max(args ...int) int {
	max := math.MinInt64
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}

// 返回值
func divFunc(a, b float64) (float64, error) {
	if a == 0 {
		return math.NaN(), errors.New("division by zero")
	}
	return a / b, nil
}

// 命名返回值
func divide(a, b float64) (quotinet float64, err error) {
	if b == 0 {
		err = errors.New("division by zero")
		return
	}
	quotinet = a / b
	return
}

// 定义一个函数
func sayHello() {
	println("Hello, world!")
}

// 直接声明
func sayHello2() {}

// 字面量
var sayHello3 = func() {}

// 类型
type sayHello4 func()
