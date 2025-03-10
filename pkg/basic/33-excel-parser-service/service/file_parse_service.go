package service

import (
	"bufio"
	"example/use/pkg/configs"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getOutPutPath() string {

	// // 获取当前工作目录
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Printf("获取当前工作目录失败：%v\n", err)
	// 	return "."
	// }
	// log.Printf("当前工作目录：%s", currentDir)

	// outputpath := filepath.Join(currentDir, "output")
	// fmt.Println(outputpath)

	return filepath.Join(configs.Cfg.BusinessConfig.LocalOutputPath)
}

func writeFile(path, content string) error {
	fileA, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(fileA *os.File) {
		_ = fileA.Close()
	}(fileA)
	writer := bufio.NewWriter(fileA)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	err = writer.Flush()
	return err
}

func ParseFeature(filePath string) {
	excelService := NewExcelService()
	excelService.ProcessExcel(filePath)
}

func uploadMarkdownToHwObs(localFileDir string) error {
	err := filepath.WalkDir(localFileDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 检查是否为文件
		if !d.IsDir() {
			objectKey := strings.TrimPrefix(path, configs.Cfg.BusinessConfig.LocalOutputPath+"/")
			uploadObjectKey := configs.Cfg.BusinessConfig.UploadPath + objectKey

			fmt.Printf("objectKey: %s\n", objectKey)
			fmt.Printf("UploadObjectKey: %s\n", uploadObjectKey)
			fmt.Println("-------------------------")
		}
		time.Sleep(1 * time.Second)
		return nil
	})
	return err
}

func deleteLocalFile(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		log.Fatal("删除文件失败：", err)
	}
	fmt.Printf("删除文件成功：%s\n", filePath)
	return nil
}
