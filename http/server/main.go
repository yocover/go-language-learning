package main

import (
	"fmt"
	"net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	// 设置响应头信息
	w.Header().Set("Content-Type", "text/plain")
	// 写入响应体
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// 注册处理器函数到路由"/"
	http.HandleFunc("/", helloWorld)
	// 监听并服务在8080端口
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
