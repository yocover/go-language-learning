package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 创建一个新的HTTP GET请求
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应内容
	fmt.Println("Response:", string(body))
}
