package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// 简化的响应结构体
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Sheets []Sheet `json:"sheets"`
	} `json:"data"`
}

type Sheet struct {
	SheetName string  `json:"sheet_name"`
	Tables    []Table `json:"tables"`
}

type Table struct {
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
}

func main() {

	// 获取当前工作mulu
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 使用filepath.Join()方法拼接文件路径
	filepath := filepath.Join(currentDir, "files", "Diagram_en.xlsx")

	// 创建表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 打开文件
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("open file error:", err)
	}
	defer file.Close()

	// 创建表单文件字段
	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		log.Fatal("create form file error:", err)
	}

	// 复制文件内容到表单
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal("copy file to form error:", err)
	}

	// 关闭表单
	writer.Close()

	// 创建请求
	req, err := http.NewRequest("POST", "http://10.18.101.42:8011/api/excel/parse", body)
	if err != nil {
		log.Fatal("create request error:", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("send request error:", err)
	}
	defer resp.Body.Close()

	// 读取相应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response body error:", err)
	}
	// fmt.Printf("相应状态码: %d\n", resp.StatusCode)
	// fmt.Printf("相应内容: %s\n", string(respBody))

	// 解析json
	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Fatal("JSON unmarshal error:", err)
	}

	// 打印结果，
	fmt.Printf("状态码： %d\n", result.Code)
	fmt.Printf("消息： %s\n", result.Message)

	// 打印第一个表格的数据
	if len(result.Data.Sheets) > 0 && len(result.Data.Sheets[0].Tables) > 0 {
		table := result.Data.Sheets[0].Tables[0]

		fmt.Println("列名：", table.Columns)
		fmt.Println("数据：", table.Values)
	}
}
