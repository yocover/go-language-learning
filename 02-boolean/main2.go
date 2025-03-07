package main

import (
	"fmt"
)

func main() {
	// 1. 布尔类型
	// Go语言中只有两个布尔类型：true和false。
	// 布尔类型的值可以直接使用，也可以作为条件表达式的判断条件。
	// 布尔类型的值可以进行逻辑运算，包括与(&&)、或(||)、非(!)运算。
	// 布尔类型的值也可以进行比较运算，包括等于(==)、不等于(!=)、大于(>)、小于(<)、大于等于(>=)、小于等于(<=)运算。

	// var a, b, c bool = true, false, true
	// fmt.Println(a && b) // false

	var age = 10
	var isAdult bool = age >= 18
	fmt.Println(isAdult) // true

	// 逻辑运算
}
