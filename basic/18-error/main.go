package main

import (
	"fmt"
	"os"
)

func main() {

	// go 中的异常有三种级别
	// error: 部分流程出错 ，需要处理
	// panic:严重错误，程序崩溃，需要修复
	// fatal：致命错误，程序崩溃，不能恢复

	// 准确来说，go并没有异常 也没有 try catch finally 这样的异常处理机制，go的异常处理机制是通过channel来实现的。
	// 异常处理机制的实现依赖于channel，channel是go中的一个重要的概念。

	// 打开一个文件
	if file, err := os.Open("test.txt"); err != nil {
		fmt.Println("文件打开失败", err)
	} else {
		fmt.Println("文件打开成功", file.Name())

		// 读取文件
		//1. Read 方法需要一个缓冲区 来存储读取的数据
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if err != nil {
			fmt.Println("文件读取失败", err)
		} else {
			fmt.Printf("文件读取成功 \n %s ", string(buffer[:n]))
		}

		// 2. 使用完文件后需要关闭
		defer file.Close()

		// panic 异常处理

		// 初始化map
		dic := make(map[string]int)
		dic["a"] = 1
	}

	// 调用一个危险操作
	dangerOperation("111")
	fmt.Println("程序执行完成")
}

func dangerOperation(str string) {
	if len(str) == 0 {
		fmt.Println("字符串为空")
		os.Exit(1) // 退出程序
	}
	fmt.Println("字符串长度为", len(str))
}
