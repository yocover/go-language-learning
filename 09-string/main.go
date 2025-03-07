package main

import (
	"fmt"
	"strings"
)

func main() {
	// 字符串字面值

	// 单引号字符串
	s1 := 'a'
	s2 := '中'

	// 双引号字符串
	s3 := "hello, world"
	s4 := "你好，世界"

	// 多行字符串
	s5 := `This is a multi-line string.
	It can span multiple lines.
	And it can contain "double quotes" and'single quotes'.`

	// 字符串拼接
	s6 := s3 + " " + s4
	s7 := s5 + " " + s6

	fmt.Printf("s1: %c\n", s1)
	fmt.Printf("s2: %c\n", s2)
	fmt.Println("s3:", s3)
	fmt.Println("s4:", s4)
	fmt.Println("s5:", s5)
	fmt.Println("s6:", s6)
	fmt.Println("s7:", s7)

	var str = "hello, world"
	fmt.Println(str)
	fmt.Println(str[0]) // 输出的是字节，不是字符
	fmt.Printf("len(str): %d\n", len(str))
	fmt.Printf("str[0]: %c\n", str[0])
	fmt.Printf("str[1:4]: %s\n", str[1:4])

	var name = "wangzhongjie"
	fmt.Println(name)
	fmt.Println(stringToByteSlice(name))
	var nameLength = stringLength(name)
	fmt.Printf("name length: %d\n", nameLength)

	var name2 = "中文"
	fmt.Println(name2)
	fmt.Println(stringToByteSlice(name2))
	var name2Length = stringLength(name2)
	fmt.Printf("name2 length: %d\n", name2Length)

	// 看起来中文的字符串的长度比英文字符换短
	// 但实际上求的的长度确比英文字符传的长度长
	// 这是因为在unicode编码中，一个汉字在大多数情况下占用两个字节，而英文字符占用一个字节

	// 上面的是错误的
	// 因为在go语言中，一个中文字符通常占用3个字节，这是因为go中默认使用的是UTF-8编码，而UTF-8编码的中文字符占用3个字节

	// 验证中文字符占用的字节数
	fmt.Printf(" 一个中文字符占用的字节数：%d\n", len([]byte("中")))
	fmt.Printf(" 一个英文字符占用的字节数：%d\n", len([]byte("a")))
	fmt.Printf(" 特殊字符字符占用的字节数：%d\n", len([]byte("😊")))

	// 特殊字符 emoji 表情
	// 特殊符号
	// 一些复杂的香型文字
	// unicode 扩展平面中的字符

	// // 字符串复制
	// var str1 = "hello, world"

	// str2 := str1
	// fmt.Printf("str1: %s\n", str1)
	// fmt.Printf("str2: %s\n", str2)

	/**
		strings.Clone() 和使用 copy() 进行字符串复制有一些重要区别：

	1. strings.Clone()

	   - Go 1.18 新增的函数
	   - 直接返回一个新的字符串，内部会分配新的内存空间
	   - 保证返回的字符串有自己独立的底层数组
	   - 实现更简单，更安全
	   - 适用于需要确保字符串完全独立的场景
	2. 使用 copy() 复制

	   - 需要手动分配目标切片
	   - 需要手动转换回字符串
	   - 可以部分复制
	   - 性能上可能略好（因为可以复用已分配的内存）
	   - 实现更灵活，但也更容易出错

		 建议：
			- 如果只是简单复制整个字符串，优先使用 strings.Clone()
			- 如果需要部分复制或者有特殊的内存管理需求，使用 copy()
	*/

	// stringCopy()
	// stringClone()
	stringConcat()
}

// 字符串拼接
func stringConcat() {

	// 直接使用 + 运算符进行字符串拼接
	/**
	- 最简单直观
	- 适合简单的、少量的字符串拼接
	- 每次拼接都会创建新的字符串，性能较差
	- 如果在循环中使用会产生大量临时对象
	*/
	s1 := "hello"
	s2 := "world"
	s3 := s1 + " " + s2
	fmt.Println(s3)

	// 转换为字节再拼接
	/**
	- 需要进行字符串和字节切片的转换
	- 可以预分配内存，减少内存分配次数
	- 适合大量字符串拼接且能预知大小的场景
	- 转换过程会有额外的内存开销
	*/
	str1 := "hello"
	bytes := []byte(str1)
	bytes = append(bytes, " world"...)
	str2 := string(bytes)
	fmt.Println(str2)

	// 可以使用strings.Builder来拼接字符串
	/**
	- 专门用于字符串拼接的高效工具
	- 内部使用 []byte 实现，但对外隐藏实现细节
	- 可以预分配内存（通过 builder.Grow() ）
	- 适合大量字符串拼接的场景
	- 性能最好，内存使用最优
	*/
	builder1 := strings.Builder{}
	builder1.WriteString("hello")
	builder1.WriteString(" world")
	fmt.Printf("builder1: %s\n", builder1.String())

	/**
	建议：

	- 简单拼接用 +
	- 大量拼接用 strings.Builder
	- 需要精细控制内存时可以考虑 []byte 方案
	*/
}

// 字符串clone
func stringClone() {

	var dest, src string
	src = "hello, world"
	dest = strings.Clone(src)

	fmt.Printf("src: %s\n", src)
	fmt.Printf("dest: %s\n", dest)
}

// 字符串copy
func stringCopy() {
	var dest, src string
	src = "hello, world"
	destBytes := make([]byte, len(src))
	// 这行代码会显示的把string 转换为 []byte，需要额外的内存分配来创建临时的字节切片，性能不好
	// copy(destBytes, []byte(src))

	// 优化后的代码，直接使用string 作为源，go编译器内部自动处理类型转换，性能更好，避免了显示转换时的临时内存分配
	copy(destBytes, src)
	dest = string(destBytes)
	fmt.Printf("src: %s\n", src)
	fmt.Printf("dest: %s\n", dest)
}

// 字符串转换为字节切片
func stringToByteSlice(s string) []byte {
	return []byte(s)
}

// 字符串的长度
func stringLength(s string) int {
	return len(s)
}
