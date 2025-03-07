package main

import (
	"fmt"
	"os"
)

func main() {
	// 文件操作reafFile的基础数据类型支持是[]byte，即字节切片。
	// reafFile1()
	reafFile2()
	getFileInfo()
	reafFile3()
}

func reafFile3() {
	bytes, err := os.ReadFile("1.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("文件内容：%s\n", string(bytes))
	}

}

func getFileInfo() {
	fileInfo, err := os.Lstat("1.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("文件名：%s\n", fileInfo.Name())
		fmt.Printf("文件大小：%d bytes\n", fileInfo.Size())
		fmt.Printf("文件权限：%v\n", fileInfo.Mode())
		fmt.Printf("最后修改时间：%v\n", fileInfo.ModTime())
		fmt.Printf("是否是目录：%v\n", fileInfo.IsDir())
	}
}

// 读写模式打开
func reafFile2() {
	file, error := os.OpenFile("readme.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if os.IsNotExist(error) {
		fmt.Println("文件不存在", error)
	} else if error != nil {
		fmt.Println("打开文件失败", error)
	} else {
		fmt.Println("文件打开成功", file.Name())
		file.Close()
	}
}

func reafFile1() {
	file, err := os.Open("11.txt")
	if os.IsNotExist(err) {
		fmt.Println(err)
	} else if err != nil {
		fmt.Println(err)
	} else {
		content := make([]byte, 1024)
		n, err := file.Read(content)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("读取了 %d 字节\n", n)
			fmt.Printf("文件内容： %s\n", string(content[:]))
		}
	}
	defer file.Close()
}
