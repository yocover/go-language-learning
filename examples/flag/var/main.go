package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// 声明变量
	var (
		name    string
		age     int
		married bool
		delay   time.Duration
	)

	// 使用 flag.TypeVar 将命令行参数绑定到变量
	// flag.TypeVar(&变量, 参数名, 默认值, 帮助信息)
	flag.StringVar(&name, "name", "world", "指定名称")
	flag.IntVar(&age, "age", 0, "指定年龄")
	flag.BoolVar(&married, "married", false, "是否已婚")
	flag.DurationVar(&delay, "delay", 0, "延迟时间")

	// 解析命令行参数
	flag.Parse()

	// 直接使用变量（不需要解引用）
	fmt.Printf("Hello, %s!\n", name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Married: %t\n", married)
	fmt.Printf("Delay: %v\n", delay)
}
