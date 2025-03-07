package main

import "fmt"

// 指针介绍
// go 保留的指针，在一定程度上保证了性能，同时为零更好的gc和安全考虑
// 又限制了指针的使用
func main() {

	// 指针有两个常用的操作符，
	// 取地址 &
	// 解引用 *

	num := 10
	// 对一个变量取地址，会返回对应类型的指针
	p := &num // 取地址
	fmt.Println("num 的地址是", p)

	// 解引用 有两个用途，第一个是访问指针所指向的元素
	rawNum := *p // 解引用
	fmt.Println("rawNum 的值是", rawNum)

	// 另一个用途就是声明一个指针

	// 在go 中为初始化的的指针默认nil
	// 直接解引用nil的指针 会导致 panic
	// 所以使用指针前，应该进行 nil检查
	// 另一个用途就是声明一个指针

	// 声明并直接初始化指针
	var numPtr *int = &num
	// 或者使用简短声明
	// numPtr := &num

	fmt.Println("numPtr 的值是", numPtr)
	fmt.Println("numPtr 指向的值是", *numPtr)

	// 短变量声明指针
	numPtr2 := new(int)
	*numPtr2 = 20
	fmt.Println("numPtr2 的值是", numPtr2)
	fmt.Printf("numPtr2 指向的值是 %d\n", *numPtr2)

	// new 函数只有一个参数，那就是类型，并返回一个对应类型的指针
	// 函数会位该指针分配内存，并且指针指向对应类型的零值

	fmt.Println(*new(string))
	fmt.Println(*new(int))
	fmt.Println(*new(bool))
	fmt.Println(*new([2]int))
	fmt.Println(*new([]int))

	NewAndMake()
}

// new 和 make 函数的区别
func NewAndMake() {

	// new func new(Type) *Type
	// make func make(Type, size...IntegerType) Type

	// new 返回值类型指针
	// 接收参数是类型
	// 专门给指针分配内存空间

	// make 返回值是值，不是类型
	// 接收的第一个参数是类型，不定长参数根据传入的类型的不同而不同
	// 专门用于给切片，映射表，通道分配内存

	a := new(int)    // int 指针
	b := new(string) // string 指针
	c := new([]int)  // int 切片指针

	d := make([]int, 10, 100)     // 长度为10， 容量为100的整形切片
	e := make(map[string]int, 10) // 容量为10的映射表
	f := make(chan int, 10)       // 容量为10的通道

	fmt.Println(a, b, c, d, e, f)
}
