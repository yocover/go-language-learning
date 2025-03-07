package main

import (
	"fmt"
	"slices"
)

func Outer() {
	println("------------ 函数调用开始 ------------")

	println("------------ 1.栈上分配内存的变量 ------------")
	GetAddress()
	println("------------ 1.栈上分配内存的变量 ------------")

	println("------------ 2.切片相关 ------------")
	println("------------ 切片的扩容 ------------")
	SliceAppend()
	println("------------ 切片和底层数组的关系 ------------")
	ShowSliceStructure()
	println("------------ 切片的删除 ------------")
	SliceDeleteItem()

	println("------------ 切片的copy ------------")
	SliceCopyItem()

	println("------------ 切片和数组的便利 ------------")
	SliceAndArrayReverse()

	println("------------ 2.切片相关 ------------")
	println("------------ 函数调用结束 ------------")
}

func GetAddress() {

	// 1. 基本变量类型 - 通常在栈上
	x := 10

	// 2. 常量 - 编译时确定
	// const y = 1290

	// 3.字符类型 - 通常在栈上
	ch := 'a'

	// 4. 字符串 - 结构体在栈上、内容在堆上
	str := "Hello, World!"

	// 5. 数组 - 通常在栈上
	arr := [5]int{1, 2, 3, 4, 5}

	// 6. 切片 - 结构体在栈上、内容在堆上
	slice := []int{1, 2, 3, 4, 5}

	// 7. 映射 - 结构体在栈上、内容在堆上
	m := map[string]int{"one": 1, "two": 2}

	// 8. 结构体 - 通常在栈上
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alex", Age: 30}

	// 9. 指针 - 通常在栈上
	var ptr *int = &x

	// 打印地址观察内存分配
	fmt.Printf("x 的地址: %p\n", &x)
	// fmt.Printf("y 的地址: %p\n", y)
	fmt.Printf("ch 的地址: %p\n", &ch)
	fmt.Printf("str 的地址: %p\n", &str)
	fmt.Printf("arr 的地址: %p\n", &arr)
	fmt.Printf("slice 的地址: %p\n", &slice)
	fmt.Printf("slice 底层数组的地址：%p\n", slice)
	fmt.Printf("m 的地址: %p\n", &m)
	fmt.Printf("p 的地址: %p\n", &p)
	fmt.Printf("ptr 的地址: %p\n", &ptr)
}

func ObserSlice() {
	// 1. 声明切片
	var s []int = []int{1, 2, 3}

	fmt.Printf(" s长度： %d, 容量： %d\n", len(s), cap(s))

	// 创建切片 长度为5 容量为10
	s2 := make([]int, 5, 10)

	// 打印切片的长度和容量
	fmt.Printf(" s2长度： %d, 容量： %d\n", len(s2), cap(s2))

	// 2. 切片的扩容
	fmt.Printf(" s2的地址： %p\n", s2)
	fmt.Printf(" s2的底层数组地址： %p\n", &s2)
	fmt.Printf(" s2的内容： %v\n", s2)

	// 添加新元素
	s2 = append(s2, 100)
	fmt.Printf(" s2长度： %d, 容量： %d\n", len(s2), cap(s2))
	fmt.Printf(" s2的地址： %p\n", s2)
	fmt.Printf(" s2的底层数组地址： %p\n", &s2)
	fmt.Printf(" s2的内容： %v\n", s2)

	// 扩展容量
	s2 = append(s2, 200)
	fmt.Printf(" s2长度： %d, 容量： %d\n", len(s2), cap(s2))

}

// 1. 切片的扩容
// 2. 切片的传递
func SliceAppend() {
	var s = make([]int, 0, 3)

	fmt.Printf(" s初始状态 s长度： %d, 容量： %d\n", len(s), cap(s))

	// 逐个添加元素 观察扩容
	for i := range 5 {
		s = append(s, i)
		fmt.Printf(" 添加元素 %d 后， s长度： %d, 容量： %d， 切片结构体地址： %p, 底层数组地址: %p\n", i, len(s), cap(s), &s, s)
	}
}

// 切片和底层数组的关系
func ShowSliceStructure() {
	// 切片本身是一个结构体，包含指向底层数组的指针
	// 多个切片可以共享一个底层数组
	// 对切片的修改会反应到底层数组上

	// 创建一个底层数组
	array := [5]int{1, 2, 3, 4, 5}

	// 创建两个切片，指向同一个底层数组
	slice1 := array[1:3] // [2,3]
	slice2 := array[1:4] // [2,3,4]

	// 修改slice1 会影响到底层数组 slice2
	slice1[0] = 20

	fmt.Printf("底层数组： %v\n", array)
	fmt.Printf("slice1： %v\n", slice1)
	fmt.Printf("slice2： %v\n", slice2)

	// 打印地址信息
	fmt.Printf("底层数组的地址： %p\n", &array)
	fmt.Printf("slice1的地址： %p\n", &slice1)
	fmt.Printf("slice2的地址： %p\n", &slice2)

	// 切片是对底层数组的引用
	// 多个切片可以共享一个底层数组
	// 通过一个切片修改数据，会影响到其他共享该底层数组的切片
	// 切片的底层数组地址可以通过直接答应切片变量来查看
	// 所以切片被称为引用类型，因为它本质上是通过指针引用底层数组的数据

	var nums []int
	fmt.Printf(" nums 的值： %v\n", nums)
	fmt.Printf(" nums 的类型： %t\n", nums == nil)

	var nums1 = []int{1, 2, 3, 4, 5}
	fmt.Printf("nums1 的值：%v，长度：%d，容量：%d\n", nums1, len(nums1), cap(nums1))

	var nums2 = make([]int, 0, 100)
	fmt.Printf("追加前 nums2 的值：%v，长度：%d，容量：%d\n", nums2, len(nums2), cap(nums2))

	nums2 = append(nums2, nums1...)
	fmt.Printf("追加后 nums2 的值：%v，长度：%d，容量：%d\n", nums2, len(nums2), cap(nums2))

	nums2 = append(nums2, 99, 1024)
	fmt.Printf("追加后 nums2 的值：%v，长度：%d，容量：%d\n", nums2, len(nums2), cap(nums2))
}

func SliceDeleteItem() {

	// // 删除单个元素
	// slice = slices.Delete(slice, index, index+1)

	// // 删除多个连续元素
	// slice = slices.Delete(slice, startIndex, endIndex)

	var nums = []int{1, 2, 3, 4, 5}
	fmt.Printf(" 原始切片： %v\n", nums)

	// 1.使用slices 包删除
	// 2 开始删除的索引位置 （包含）
	// 4 结束删除的索引位置（不包含）
	nums = slices.Delete(nums, 2, 3)
	fmt.Printf("使用 slices.Delete 2-4 后：%v\n", nums)

	// 2. 使用append 的方式删除
	nums = []int{1, 2, 3, 4, 5}
	i := 2 // 要删除的索引
	nums = append(nums[:i], nums[i+1:]...)
	// nums = slices.Delete(nums, i, i+1)
	// fmt.Printf(" 使用 append 删除索引2: %v\n", nums)
	fmt.Printf(" 使用 append 删除索引2: %v\n", nums)
	// 3. 使用copy 的方式删除
	nums = []int{1, 2, 3, 4, 5}
	i = 2                      // 要删除的索引
	copy(nums[i:], nums[i+1:]) // 复制后面的元素到前面
	nums = nums[:len(nums)-1]  // 切片长度减1
	fmt.Printf(" 使用 copy 删除索引2: %v\n", nums)

	// 删除元素
	// 从头部删除n个元素
	nums = []int{1, 2, 3, 4, 5}
	n := 2
	nums = nums[n:]
	fmt.Printf(" 从头部删除2个元素后：%v\n", nums)

	// 从尾部删除n个元素
	nums = []int{1, 2, 3, 4, 5}
	m := 2
	nums = nums[:len(nums)-m]
	fmt.Printf(" 从尾部删除2个元素后：%v\n", nums)

	// 从中间指定下标的位置删除元素
	nums = []int{1, 2, 3, 4, 5}
	index := 2
	// nums = slices.Delete(nums, index, index+1)
	// 通过跳过要删除的元素将其前后的元素重新组合来实现从中间删除元素 可以使用上面的slices.Delete来说明
	nums = append(nums[:index], nums[index+1:]...)
	fmt.Printf(" 从中间删除索引2的元素后：%v\n", nums)
}

func SliceCopyItem() {

	// 切片在copy时 要确保目标切片有足够的长度
	dest := make([]int, 0)
	src := []int{1, 2, 3, 4, 5}

	fmt.Printf(" 原始切片： %v\n", src)
	fmt.Printf(" 目标切片： %v\n", dest)

	// 1.使用copy函数
	dest1 := make([]int, len(src))
	n := copy(dest1, src)
	fmt.Printf(" 使用copy函数后：%v, 复制的元素个数：%d\n", dest1, n)

	// 2. 使用slices.Clone 方法
	dest2 := slices.Clone(src)
	fmt.Printf(" 使用slices.Clone 方法后：%v\n", dest2)

	// 3. 使用append方法
	var dest3 []int
	dest3 = append(dest3, src...)
	fmt.Printf(" 使用append方法后：%v\n", dest3)
}

// 遍历
func SliceAndArrayReverse() {
	// 数组
	var arr = [5]int{1, 2, 3, 4, 5}

	// 1. 使用普通for循环遍历
	fmt.Printf("\n 普通for循环遍历数组：")
	for i := 0; i < len(arr); i++ {
		fmt.Printf(" 数组元素： %d\n", arr[i])
	}

	// 2. 使用 for rang 遍历
	fmt.Printf("\n for range 遍历数组")
	for i, v := range arr {
		fmt.Printf(" 数组元素 arr[%d]： %d\n", i, v)
	}

	// 3. 值需要索引
	fmt.Printf(" \n只遍历索引")
	for i := range arr {
		fmt.Printf(" 索引： %d\n", i)
	}

	// 4. 只需要值
	fmt.Printf(" \n只遍历值")
	for _, v := range arr {
		fmt.Printf(" 值： %d\n", v)
	}

	// 切片的遍历方式与数组相同
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("\n 遍历切片：")
	for i, v := range slice {
		fmt.Printf(" 切片元素 slice[%d]： %d\n", i, v)
	}

	// 多维切片
	var nums_lsit [4][4]int

	for i := range 4 {
		for j := range 4 {
			nums_lsit[i][j] = i*4 + j
		}
	}

	for _, arr := range nums_lsit {
		fmt.Printf("\n 遍历二维数组：%v", arr)
	}

	// 使用嵌套 for range 遍历二维数组
	fmt.Println("\n完整遍历二维数组：")
	for i, row := range nums_lsit {
		for j, val := range row {
			fmt.Printf(" nums_lsit[%d][%d]： %d\n", i, j, val)
		}
		fmt.Printf("\n")
	}

	fmt.Println("\n 按行遍历二维数组：")
	for i, row := range nums_lsit {
		fmt.Printf("第 %d 行：%v\n", i, row)
	}

	// 拓展表达式
	s1 := []int{1, 2, 3, 4, 5} // cap 5
	s2 := s1[1:4]              // cap 5-1 = 4

	s2 = append(s2, 66) // 扩容
	// 由于s2的容量组否，所以扩容的时候 把s1的值也修改了
	// 输出结果是 [1 2 3 4 66] [2 3 4 66] 5 4
	fmt.Println(s1, s2, cap(s1), cap(s2))

	// 解决办法 使用扩展表达式，注意扩展表达式只能适用于切片
	// slice[low:high:max]
	// low 切片的起始索引
	// high 切片的结束索引
	// max 指最大容量

	/**
			s2 := s1[3:4:4]

		- 使用三索引切片表达式： slice[start:end:cap]
		- start=3：从索引3开始
		- end=4：到索引4（不含）
		- cap=4：新切片的容量到原切片的索引4为止

		s2 = append(s2, 1)

	- 因为 s2 的容量已经是1，没有更多空间
	- append 会创建新的底层数组
	- 不会影响原始切片 s1

	*/
	ss1 := []int{1, 2, 3, 4, 5} // cap 5
	ss2 := ss1[1:4:4]           // cap 4-1 = 3
	ss2 = append(ss2, 66)
	// [1 2 3 4 5] [2 3 4 66] 5 6
	fmt.Println(ss1, ss2, cap(ss1), cap(ss2))
}
