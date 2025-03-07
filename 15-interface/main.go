package main

import "fmt"

// go 1.18 加入了反省
// 因此接口其实分为两类
// 基本接口 ：值包含方法集的接口，就是基本接口
// 通用接口：只要包含类型集的接口就是通用接口

// 方法集：就是一组方法的集合
// 类型集：就是一组类型的集合

// 定义接口
type Animal interface {
	Speak() string
	Move() string
}

// 定义实现接口的结构体
type Dog struct {
	Name string
}

// Dog 实现 Animal 接口
func (d Dog) Speak() string {
	return "汪汪!"
}
func (d Dog) Move() string {
	return "狗在跑!"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "喵喵!"
}
func (c Cat) Move() string {
	return "猫在跑!"
}

// 通用接口
func AnimalSpeak(a Animal) {
	fmt.Printf("动物叫声: %s\n", a.Speak())
}
func AnimalMove(a Animal) {
	fmt.Printf("动物移动: %s\n", a.Move())
}

func main() {

	// 创建实例
	dog := Dog{Name: "旺财"}
	cat := Cat{Name: "花椒"}

	// 通过接口调用方法
	AnimalSpeak(dog)
	AnimalMove(dog)
	AnimalSpeak(cat)
	AnimalMove(cat)

	// 空接口使用
	var i interface{}
	i = 43
	fmt.Printf("类型 %T, 值 %v\n", i, i)

	// 接口的重要特点
	// 1. 接口实现：不需要显示声明实现了某个接口
	// 2.组合：接口额可以组合其它接口
	// 3.灵活性：一个类型可以实现多个接口
	// 4.面向接口编程：以来抽象而不是具体实现

	// 注意实现
	// 1.接口尽量保持小巧，单一职责
	// 2.使用类型断言时要做错误检查
	// 3.空接口可以接收任意类型，但要谨慎使用
	// 4.接口名称通常以er结尾，比如 Reader、Writer、Closer等
	// 5.指针接收者和值接收者的区别要注意

	println("---------------------------")
	company := ConstructionCompany{CraneA{}}
	company.BuildCrane()
	println()
	fmt.Println("更换起重机B")
	println()
	company.Crane = CraneB{}
	company.BuildCrane()

	// any 接口内部没有方法集合，根据实现定义
	// 所有类型都是any 接口的实现，因为所有类型的方法集
	// 都是空集的超集，所以any接口可以保存任何类型的值
	// var anything any
	// anything = company

	// 在比较空接口时 会对底层类型进行比较
	// 如果类型不匹配的话 则为flase
	var a any
	var b any
	a = 1
	b = "hello"
	fmt.Println(a == b)
	a = 12
	b = 12
	fmt.Println(a == b)
}

// type any = interface{} // 使用any代替空接口
func DoSomething(a any) interface{} {
	// interface{} 类型可以保存任何类型的值
	// 因此可以将 any 接口转换为 interface{} 类型
	return a
}
