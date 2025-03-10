package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 创建子命令
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)

	// server 子命令的参数
	serverPort := serverCmd.Int("port", 8080, "服务器端口")
	serverHost := serverCmd.String("host", "localhost", "服务器主机名")
	serverMode := serverCmd.String("mode", "development", "服务器模式 (development/production)")

	// client 子命令的参数
	clientServer := clientCmd.String("server", "localhost:8080", "服务器地址")
	clientTimeout := clientCmd.Int("timeout", 30, "连接超时时间(秒)")
	clientDebug := clientCmd.Bool("debug", false, "启用调试模式")

	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("请指定子命令: server 或 client")
		fmt.Println("使用示例:")
		fmt.Println("  ./subcommand server -port=9090 -mode=production")
		fmt.Println("  ./subcommand client -server=example.com:9090 -debug")
		os.Exit(1)
	}

	// 根据第一个参数确定子命令
	switch os.Args[1] {
	case "server":
		// 解析 server 子命令的参数
		serverCmd.Parse(os.Args[2:])
		fmt.Println("启动服务器:")
		fmt.Printf("  主机: %s\n", *serverHost)
		fmt.Printf("  端口: %d\n", *serverPort)
		fmt.Printf("  模式: %s\n", *serverMode)

	case "client":
		// 解析 client 子命令的参数
		clientCmd.Parse(os.Args[2:])
		fmt.Println("启动客户端:")
		fmt.Printf("  服务器: %s\n", *clientServer)
		fmt.Printf("  超时: %d秒\n", *clientTimeout)
		fmt.Printf("  调试模式: %t\n", *clientDebug)

	default:
		fmt.Printf("未知子命令: %s\n", os.Args[1])
		fmt.Println("可用子命令: server, client")
		os.Exit(1)
	}

	// 显示未解析的参数
	var args []string
	if serverCmd.Parsed() {
		args = serverCmd.Args()
		fmt.Println("未解析的 server 参数:", args)
	} else if clientCmd.Parsed() {
		args = clientCmd.Args()
		fmt.Println("未解析的 client 参数:", args)
	}
}
