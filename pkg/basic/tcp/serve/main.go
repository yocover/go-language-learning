package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Print("Message Received:", string(message))
		conn.Write([]byte("Message received.\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
