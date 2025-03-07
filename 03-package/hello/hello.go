package hello

import (
	_ "example/say" // 导入say包
	"fmt"           // 导入fmt包
)

// init()函数在包被导入时执行一次，一般用于初始化包内的变量或函数。
func init() {
	println("Hello Package init() is called")
}

func SayHello() {

	var fnum = 123i

	num := 2*9 + 1

	println("Num:", num)

	// fmt.Println("Hello, world! str: %v", str)

	fmt.Println("Hello, world! fnum:", fnum)

	var count = 100_000_000 // 声明变量并赋值
	fmt.Println("Hello, world! count:", count)

	fmt.Println("Hello package!")
	// 调用SayHello()函数时，会自动调用init()函数。
}

// 注意：
// 1. 包名必须与文件名相同。
