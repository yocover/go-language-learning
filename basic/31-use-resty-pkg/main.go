package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"
	"fmt"

	"git.neuxnet.com/ai/global-kit/net/resty"
	"git.neuxnet.com/ai/global-kit/tools"
	"go.uber.org/zap"
)

type FileParserRequest struct {
	FileId     string `json:"fileId"`
	FilePath   string `json:"filePath"`
	Filename   string `json:"filename"`
	FileFormat string `json:"fileFormat"`
	IsFixTitle bool   `json:"isFixTitle"`
}

type DoParserResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *DoParser `json:"data"`
}

type DoParser struct {
	Id int64 `json:"id"`
}

var Header = map[string]string{
	"Content-Type": "application/json",
}

func main() {
	// DoParserFunc("123", "adagpt/tmpfile/1872217666842804224", "test.pdf", "pdf")
	DoExcelParser()
}

func DoExcelParser() {
	// 获取文件路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filePath := filepath.Join(currentDir, "files", "Diagram_en.xlsx")

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open file error:", err)
	}
	defer file.Close()

	// 发送请求
	resp, err := resty.File(
		"http://10.18.101.42:8011/api/excel/parse", // URL
		nil,               // 额外的表单数据
		nil,               // 请求头
		"file",            // 文件字段名
		"Diagram_en.xlsx", // 文件名
		file,              // 文件读取器
	)
	if err != nil {
		log.Fatal("upload file error:", err)
	}

	// 处理响应
	fmt.Println(string(resp))
}

// DoParser
/**
* @Description
* @Author zhouxl 2024-11-18 16:10:42
* @Param fileId
* @Param fileKey
* @Param fileName
**/
func DoParserFunc(fileId, fileKey, fileName, fileType string) string {
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex != -1 && lastDotIndex+1 < len(fileName) {
		fileType = fileName[lastDotIndex+1:]
	}
	r := &FileParserRequest{
		FileId:     fileId,
		FilePath:   fileKey,
		Filename:   fileName,
		FileFormat: fileType,
		IsFixTitle: true,
	}
	resp, err := resty.Post("http://localhost:21003"+"/api/pdf2markdown", r, Header)
	if err != nil {
		zap.L().Error("upload file parser error", zap.Any("resp", resp), zap.Error(err))
	}
	var obj *DoParserResponse
	_ = json.Unmarshal(resp, &obj)
	if obj != nil && obj.Data != nil {
		return tools.Int64ToString(obj.Data.Id)
	}
	return ""
}
