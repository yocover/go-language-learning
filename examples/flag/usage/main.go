package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 定义命令行参数
	configFile := flag.String("config", "config.json", "配置文件路径")
	port := flag.Int("port", 8080, "服务器端口")
	verbose := flag.Bool("verbose", false, "启用详细日志")

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "一个演示自定义帮助信息的示例程序\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s -config=app.json -port=9090\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -verbose\n", os.Args[0])
	}

	// 解析命令行参数
	flag.Parse()

	// 显示解析结果
	fmt.Println("配置文件:", *configFile)
	fmt.Println("端口:", *port)
	fmt.Println("详细模式:", *verbose)

	// 提示用户如何获取帮助
	fmt.Println("\n提示: 使用 -h 或 -help 参数查看帮助信息")
}
