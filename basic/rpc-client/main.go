package main

import (
	"example/rpc/client/service"
	"log"
	"net"
	"net/rpc"
)

func main() {

	// 注册服务
	userService := &service.UserServiceImpl{
		Users: make(map[int]*service.User),
	}
	rpc.Register(userService)

	// 启动RPC 服务器
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	log.Println("RPC server is running... PORT: 1234")
	rpc.Accept(listener)
}
