package pointer

import "fmt"

func init() {
	fmt.Println("Pointer init function")
}

// 演示指针使用的函数

func DemoPointer() {
	var a int = 4
	var b *int = &a
	fmt.Println("a =", a)
	fmt.Println("b =", b)
	fmt.Println("*b =", *b)
	*b = 5
	fmt.Println("a =", a)
	fmt.Println("b =", b)
	fmt.Println("*b =", *b)
}
