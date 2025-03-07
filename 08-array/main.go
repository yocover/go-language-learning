package main

import "fmt"

// func main() {
// 	nums := [5]int{1, 2, 3, 4, 5}

// 	for i := range nums {
// 		println(i)
// 	}

// 	println("------------")
// 	// print type of nums
// 	println("Type nums : %", nums)

// 	var counts = new([10]int)
// 	for i := range counts {
// 		println(i)
// 	}
// }

func main() {
	nums := [5]int{1, 2, 3, 4, 5}
	var counts = new([10]int)

	// 打印 nums 的值
	// fmt.Printf("nums 的值: %v\n", nums)
	// fmt.Printf("nums 的类型: %T\n", nums)

	fmt.Printf("nums 的值： %v\n", nums)
	fmt.Printf("nums 的类型： %T\n", nums)

	// 打印 counts 的值
	fmt.Printf("counts 的值: %v\n", *counts)
	fmt.Printf("counts 的类型: %T\n", counts)

	// 栈上分配内存和堆上分配内存的区别

	// 栈上分配内存
	/**
		  1. 栈上分配内存的变量在函数返回后会被自动释放
		  2. 栈上分配内存的变量的生命周期和函数的生命周期相同
		  3. 栈上分配内存的变量的大小是固定的
		  4. 栈上分配内存的变量的速度比堆上分配内存的变量的速度快
		  5. 栈上变量的地址会随着函数调用而变化
		  6. 同一函数内的栈变量地址是连续分配的

	- 栈内存的分配是从高地址向低地址增长的
	- 每次函数调用都会创建新的栈帧（stack frame）
	- 栈变量的具体地址取决于：
	- 程序运行时的栈指针位置
	- 函数调用的深度
	- 操作系统的内存布局
	- 程序的加载位置等因素
	*/

	// 给出栈上分配内存的示例
	var stackArray = [3]int{}
	fmt.Printf("栈上数组： %v\n", stackArray)

	for i := range stackArray {
		stackArray[i] = i * 10
		addr := &stackArray[i]
		// 打印每个栈变量的地址
		fmt.Printf("stackArray[%d] 的地址： %p, value: %d\n", i, addr, stackArray[i])
	}

	println("------------ 调用 GetAddress 函数 ------------")
	// 查看内存地址
	Outer()
	// 堆上分配内存
	/**
	  1. 堆上分配内存的变量在函数返回后不会被自动释放
	  2. 堆上分配内存的变量的生命周期和函数的生命周期不同
	  3. 堆上分配内存的变量的大小是不固定的
	  4. 堆上分配内存的变量的速度比栈上分配内存的变量的速度慢
	  5. 堆上变量的地址在分配后保持稳定，直到被垃圾回收
	  6. 同一函数内的堆变量地址是不连续分配的

	  补充说明：
	  1. 在 Go 中堆内存由 GC 自动管理，不需要手动释放
	  2. 堆内存分配和回收的开销较大，可能导致内存碎片
	  3. 堆内存的大小受系统可用内存限制
	*/
}
