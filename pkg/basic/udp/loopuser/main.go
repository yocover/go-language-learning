package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 服务器地址和端口
	serverAddr := "127.0.0.1:9999"
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to send (or type 'exit' to quit): ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1] // 去除换行符

		if text == "exit" {
			fmt.Println("Exiting...")
			break
		}

		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
		fmt.Println("Message sent to server")

		// 准备接收服务器的回复
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println("Server response:", string(buffer[:n]))
	}
}
