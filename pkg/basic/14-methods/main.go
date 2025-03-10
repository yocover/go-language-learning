package main

import (
	"fmt"
)

// 函数：独立的功能模块，不属于任何类型
func add(a, b int) int {
	return a + b
}

// 定一个类型
type Rectangle struct {
	width, height float64
}

// 方法： 属于特定类型的函数，有接收者 receiver
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// 指针接收者的方法：可以修改接受者的值
func (r *Rectangle) SetWidth(width float64) {
	r.width = width
}

func main() {
	// 方法与函数的区别

	// 函数掉哟ing
	sum := add(2, 3)
	fmt.Println("sum:", sum)

	// 方法调用
	rect := Rectangle{width: 10, height: 5}
	area := rect.Area()
	fmt.Println("area:", area)
	fmt.Printf("%+v\n", rect)

	// 修改指针接收者的值
	rect.SetWidth(20)
	fmt.Println("rect width:", rect.width)

	area = rect.Area()
	fmt.Println("area:", area)

	// 函数与方法的区别
	// 定义方式
	// 调用方式
	// 作用于
	// 函数：包级别
	// 方法：绑定到特定类型

	// 功能特点
	// 函数：通用功能
	// 方法：OOP思想，实现类型的行为

	// 接收者
	// 函数：无接收者
	// 方法：有接收者，（值接收者、指针接收者）

	// 什么事接收者？？？
	// 在go 中 接收者：是指方法所属的类型实例

	ReceiverUse()
}

type Person struct {

	// 首字母小写 在类型中 表明是私有字段，只能在当前包内访问
	name string
	age  int
}

// 值接收者：接收者是 Person 类型
func (p Person) GetInfo() string {
	return fmt.Sprintf("name:%s,age:%d", p.name, p.age)
}

// 指针接收者：接收者是 *Person 类型
func (p *Person) SetName(name string) {
	p.name = name // 可以修改接收者的值
}

// 在go中，标识符，包括方法名、变量名、类型名的首字母，大小写有重要含义

// 大写字母开头
// 可以被其它包访问和使用
func ReceiverUse() {

	person := Person{name: "wangzhongjie", age: 25}
	fmt.Println(person.GetInfo())

}
