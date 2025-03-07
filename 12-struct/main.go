package main

import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

type Color struct {
	r, g, b int
}

type Size struct {
	width, height, depth int
}

func main() {
	// go 抛弃了类与继承，同时也抛弃了构造方法
	// 刻意弱化了面向对象的功能
	// go 并非是一个oop的语言，但是go依然有着oop的影子
	// 通过结构体可以模拟出一个类
	// 结构体可以存储一组类型不同的数据，是一种复合类型

	// 结构体的声明
	type Person struct {
		name      string
		job, city string
		age       int
	}

	p := Person{
		name: "Alice",
		job:  "Engineer",
		city: "Beijing",
		age:  25,
	}

	fmt.Printf("Name: %s, Job: %s, City: %s, Age: %d\n", p.name, p.job, p.city, p.age)
	fmt.Printf("Struct size: %d bytes\n", unsafe.Sizeof(p))
	type Rectangle struct {
		width, height, area int
		color               string
	}

	type Box struct {
		name  string
		size  Size
		color Color
	}

	box := Box{
		name: "Box1",
		size: Size{
			width:  10,
			height: 20,
			depth:  30,
		},
		color: Color{
			r: 255,
			g: 0,
			b: 0,
		},
	}
	fmt.Printf("Box : %T, %#v\n", box, box)
	fmt.Printf("Box : %+v\n", box)

	fmt.Println(box.color.g)

	/**
	color1 是值类型
	直接在栈上分配内存
	赋值或者传递时会复制整个结构体
	适合小型结构
	函数修改不会影响原始值
	*/
	color1 := Color{
		r: 0,
		g: 255,
		b: 0,
	}

	/**
	color2 是指针类型
	在堆上分配内存
	赋值或者传递时只复制指针
	适合大型结构体
	修改会影响原始值
	*/
	color2 := &Color{
		r: 0,
		g: 255,
		b: 0,
	}

	/**
	小型结构体和大型结构体的区分
	小型结构体
	1. 小于等于 2 - 3个字段 用值类型
	2. 需要修改原始值的时候 使用指针类型
	3. 并发安全考虑时 根据具体场景选择

	*/

	modifyColor(color1)
	fmt.Println(color1.r)
	modifyColorPointer(color2)
	fmt.Println(color2.r)

	// labeluse()    // 标签的使用
	emptyStruct() // 空结构体的应用场景
}

func modifyColor(c Color) {
	c.r = 123 // 不会影响原始值
}

func modifyColorPointer(c *Color) {
	c.r = 456 // 会影响原始值
}

// 标签的使用
/*
JSON/XML 序列化
type Person struct {
    Name string `json:"name" xml:"name"`
    Age  int    `json:"age,omitempty" xml:"age"`
}
- 定义序列化后的字段名
- omitempty 表示值为零值时忽略该字段

表单验证
type User struct {
    Username string `validate:"required,min=3,max=20"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"gte=0,lte=130"`
}

数据库映射
type Product struct {
    ID    int64  `gorm:"primary_key"`
    Name  string `gorm:"column:product_name;type:varchar(100)"`
    Price float64 `gorm:"not null"`
}

配置文件映射
type Config struct {
    Host string `yaml:"host" env:"SERVER_HOST"`
    Port int    `yaml:"port" env:"SERVER_PORT"`
}

标签通常需要配合反射机制使用，常用于：

- 数据格式转换
- 数据验证
- 配置管理
- ORM 映射
*/
func labeluse() {
	// 标签时一种元编程的形式，结合反射
	// 标签可以用来给结构体字段添加额外的元信息
	// 标签的语法是以反引号开头，后面跟着标签名和内容
	type Person struct {
		Name string `json:"name" validate:"required"`
		Job  string `json:"job"`
		City string `json:"city"`
		Age  int    `json:"age,omitempty"`
	}

	p := Person{
		Name: "Alice",
		Job:  "Engineer",
		City: "Beijing",
		Age:  0,
	}

	fmt.Printf("Person: %+v\n", p)

	// 序列华为JSON
	jsonData, _ := json.Marshal(p)
	fmt.Printf("JSON: %s\n", jsonData)

}

// 空结构体的应用场景
func emptyStruct() {

	// 1. 作为集合的实现 Set
	set := make(map[string]struct{}) // 空结构体作为键值类型
	set["Alice"] = struct{}{}
	set["Bob"] = struct{}{}
	fmt.Println(set)

	// 检查元素是否存在
	if _, exists := set["Alice"]; exists {
		fmt.Println("Alice exists")
	}

	// // 2. 作为信号通道
	// done := make(chan struct{})
	// go func() {
	// 	// do something

	// 	// <- 这个是go 语言中的通道操作符号，用于接收和发送数据 它有两种用法
	// 	// 箭头指向通道，表示发送数据
	// 	done <- struct{}{} // 发送信号

	// 	// 箭头指向变量，表示接收数据
	// 	value := <-done // 接收信号

	// 	// 仅接收但不使用
	// 	<-done // 接收信号
	// }()

	// 3. 占位符
	type Node struct {
		data int
		_    struct{} // 空结构体作为占位符
	}

	// 空结构体不占用内存
	fmt.Printf("Empty struct size: %d bytes\n", unsafe.Sizeof(struct{}{}))

	channeluse() // 调用通道的使用方式
}

// 通道肩头的使用方式
func channeluse() {
	done := make(chan struct{})

	// 启动
	go func() {
		fmt.Println("go routine start!")
		time.Sleep(4 * time.Second) // 模拟耗时操作
		done <- struct{}{}          // 发送信号
		fmt.Println("go routine end!")
	}()

	// 主 go routine 等待信号
	<-done // 从通道接收数据
	fmt.Println("main routine end!")
}
