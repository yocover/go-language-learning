package main

import (
	"flag"
	"fmt"
)

func main() {
	// 标准库 flag 不直接支持短标志形式，但可以通过定义多个标志实现类似功能
	// 定义长标志
	namePtr := flag.String("name", "world", "指定名称")
	helpPtr := flag.Bool("help", false, "显示帮助信息")

	// 定义短标志（与长标志指向同一个变量）
	flag.StringVar(namePtr, "n", "world", "指定名称 (简写)")
	flag.BoolVar(helpPtr, "h", false, "显示帮助信息 (简写)")

	// 解析命令行参数
	flag.Parse()

	// 如果指定了帮助标志，显示使用说明
	if *helpPtr {
		fmt.Println("使用说明:")
		fmt.Println("  -name, -n: 指定名称")
		fmt.Println("  -help, -h: 显示帮助信息")
		return
	}

	fmt.Printf("Hello, %s!\n", *namePtr)

	// 提示：可以使用 -name=value 或 -n=value 的形式
	fmt.Println("\n使用示例:")
	fmt.Println("  ./shorthand -name=John")
	fmt.Println("  ./shorthand -n=John")
	fmt.Println("  ./shorthand -h")
}
