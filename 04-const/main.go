package main

import "fmt"

const (
	a = 10
	b = 20
	c = a + b
) // 常量声明

const (
	Count = 10
	Name  = "wangzhongjie"
	Age   = 25
	Num   = iota * 2
	_
	NUm1
	Num2
	Num3
)

func main() {

	const name = "wangzhongjie" // 常量声明

	const numExpression = 10 + 20 // 常量表达式

	fmt.Println(a, b, c, name, numExpression)
	fmt.Println(Num, NUm1, Num2, Num3) // iota 常量
	fmt.Println(Count)
}

// Output: 10 20 30 wangzhongjie 30 0 1 2 3 10
