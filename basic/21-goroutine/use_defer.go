package main

import "fmt"

/*
defer 是go语言中的一个关键字
用于延迟函数或者方法的执行

直到当前函数返回之前

 1. 执行顺序
    defer 语句会按照先进先出的顺序执行
    多个 defer 语句会形成一个栈，最后声明的先执行

 2. 常见用途
    资源清理，关闭文件、数据库连接等
    解锁互斥锁
    错误处理
*/
func UseDefer() {

	// 1. 基本使用
	defer fmt.Println("我是最先 defer 的，但是最后执行")
	defer fmt.Println("我是第二个 defer 的，第二个执行")
	defer fmt.Println("我是最后  defer 的，最先执行")
	fmt.Println("正常代码，最先执行")

	// 2. 资源清理
	demoFile()
}

func demoFile() {

	// 模拟打开文件
	fmt.Println("打开文件")

	// defer 确保文件最后一定会关闭
	defer fmt.Println("关闭文件")

	// 使用文件
	fmt.Println("使用文件")
}
