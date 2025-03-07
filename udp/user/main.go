package main

import (
	"fmt"
	"net"
)

func main() {
	// 服务器地址和端口
	serverAddr := "127.0.0.1:9999"
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 建立连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	message := []byte("Hello, UDP Server!")
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Println("Message sent to server")

	// 接收服务器响应
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Server response:", string(buffer[:n]))
}
