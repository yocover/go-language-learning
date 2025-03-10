package main

import "fmt"

// go 是一个静态强类型语言
// 静态指的是go编译器在编译代码时就已经确定了变量的类型，而动态语言则是在运行时才确定变量的类型
// 强类型语言的好处是可以避免很多运行时错误，比如类型不匹配导致的崩溃，以及类型转换导致的错误。

// 类型声明

type MyInt int

func main() {
	var a MyInt = 10
	var b int = 20

	// error
	// sum := a + b

	// 类型转换
	// sum := int(a) + b
	sum := a + MyInt(b)
	fmt.Printf("sum is %d\n", sum)

	fmt.Printf("a is %d, a type: %T\n", a, a)

	// 类型断言
	var c any = "hello"
	// var d int = 1

	// 类型断言语法
	// vale, ok := x.(T)
	// x 是接口类型变量
	// T 是要断言的类型
	// val 是转换后的值
	// ok 是布尔值，表示是否成功转换

	// 方式1
	if val, ok := c.(int); ok {
		fmt.Printf("c is int, value is %d\n", val)
	} else {
		fmt.Printf("c is not int\n")
	}

	// // 方式2 不推荐 会导致panic
	// val := c.(int)
	// if val != 0 {
	// 	fmt.Printf("c is int, value is %d\n", val)
	// }

	switch c.(type) {
	case int:
		fmt.Printf("c is int\n")
	case string:
		fmt.Printf("c is string\n")
	default:
		fmt.Printf("c is not int or string\n")
	}
}
