package main

import (
	"context"
	"log"
	"net"

	"example/examples/kitex_demo/kitex_gen/api"
	"example/examples/kitex_demo/kitex_gen/api/calculator"

	"github.com/cloudwego/kitex/server"
)

// CalculatorImpl 实现生成的接口
type CalculatorImpl struct{}

// Add 实现加法
func (s *CalculatorImpl) Add(ctx context.Context, req *api.CalcRequest) (resp *api.CalcResponse, err error) {
	result := req.A + req.B
	return &api.CalcResponse{
		Result: result,
	}, nil
}

// Subtract 实现减法
func (s *CalculatorImpl) Subtract(ctx context.Context, req *api.CalcRequest) (resp *api.CalcResponse, err error) {
	result := req.A - req.B
	return &api.CalcResponse{
		Result: result,
	}, nil
}

func main() {
	impl := &CalculatorImpl{}
	svr := calculator.NewServer(impl, server.WithServiceAddr(&net.TCPAddr{Port: 8888}))
	err := svr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
