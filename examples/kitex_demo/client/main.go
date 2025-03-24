package main

import (
	"context"
	"log"
	"time"

	"example/examples/kitex_demo/kitex_gen/api"
	"example/examples/kitex_demo/kitex_gen/api/calculator"

	"github.com/cloudwego/kitex/client"
)

func main() {
	client, err := calculator.NewClient(
		"calculator-service",
		client.WithHostPorts("127.0.0.1:8888"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 测试加法
	addReq := &api.CalcRequest{
		A: 10,
		B: 20,
	}
	addResp, err := client.Add(context.Background(), addReq)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Add result: %d + %d = %d\n", addReq.A, addReq.B, addResp.Result)

	// 测试减法
	subReq := &api.CalcRequest{
		A: 30,
		B: 15,
	}
	subResp, err := client.Subtract(context.Background(), subReq)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Subtract result: %d - %d = %d\n", subReq.A, subReq.B, subResp.Result)

	time.Sleep(time.Second) // 等待日志输出
}
