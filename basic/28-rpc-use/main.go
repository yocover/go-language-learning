package main

import (
	"example/rpc/client/service"
	"flag"
	"log"
	"net/rpc"
	"time"
)

func main() {

	// 连接RPC Server

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	// 调用远程服务
	user := service.User{Id: 1, Name: "张三", Age: 20}

	var response bool
	err = client.Call("UserServiceImpl.CreateUser", user, &response)
	if err != nil {
		log.Fatal("调用失败:", err)
	}

	time.Sleep(2 * time.Second)

	// 获取用户
	var result service.User
	err = client.Call("UserServiceImpl.GetUser", 1, &result)
	if err != nil {
		log.Fatal("调用失败:", err)
	}
	log.Printf("获取用户: %+v", result)

	var name string
	var age int
	var sex bool
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 15, "年龄")
	flag.BoolVar(&sex, "sex", true, "性别")
}
