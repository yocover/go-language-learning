package main

import (
	"flag"
	"fmt"
)

func main() {
	// 定义命令行参数
	// flag.Type(参数名, 默认值, 帮助信息)
	name := flag.String("name", "world", "指定名称")
	age := flag.Int("age", 0, "指定年龄")
	married := flag.Bool("married", false, "是否已婚")
	delay := flag.Duration("delay", 0, "延迟时间")

	// 解析命令行参数
	flag.Parse()

	// 使用参数值（注意使用指针）
	fmt.Printf("Hello, %s!\n", *name)
	fmt.Printf("Age: %d\n", *age)
	fmt.Printf("Married: %t\n", *married)
	fmt.Printf("Delay: %v\n", *delay)

	// 获取未解析的参数
	fmt.Println("未解析的参数:", flag.Args())
	fmt.Printf("未解析的参数数量: %d\n", flag.NArg())
	fmt.Printf("已解析的参数数量: %d\n", flag.NFlag())
}
