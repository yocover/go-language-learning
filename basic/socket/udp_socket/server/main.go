package main

import (
	"fmt"
	"net"
)

func main() {
	addr := "127.0.0.1:9999"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on", addr)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		fmt.Printf("Received %d bytes from %v: %s\n", n, remoteAddr, string(buf[:n]))

		_, err = conn.WriteToUDP([]byte("Message received"), remoteAddr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
		}
	}
}
