package main

import "fmt"

// 起重机接口
type Crane interface {
	JackUp() string
	Hoist() string
}

// 起重机A
// ----------------
type CraneA struct {
	work int // 内部字段不同代表 内部实现细节不一样
}

func (c CraneA) Work() {
	fmt.Println("使用技术A")
}
func (c CraneA) JackUp() string {
	c.Work()
	return "起重机A起吊"
}
func (c CraneA) Hoist() string {
	c.Work()
	return "起重机A提起"
}

// ----------------
// 起重机B
type CraneB struct {
	boot int
}

func (c CraneB) Boot() {
	fmt.Println("使用技术B")
}
func (c CraneB) JackUp() string {
	c.Boot()
	return "起重机B起吊"
}
func (c CraneB) Hoist() string {
	c.Boot()
	return "起重机B提起"
}

type ConstructionCompany struct {
	Crane Crane // 指根据Crane 类型 来存放起重机
}

func (c *ConstructionCompany) BuildCrane() {
	fmt.Println("正在建造起重机")
	fmt.Println(c.Crane.JackUp())
	fmt.Println(c.Crane.Hoist())
	fmt.Println("建造完成")
}
