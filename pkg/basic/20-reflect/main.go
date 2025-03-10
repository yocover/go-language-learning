package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) SayHello(msg string) {
	fmt.Printf("Hello, my name is %s, %s\n", p.Name, msg)
}

func main() {

	//1. 基础类型的反射
	num := 43
	numValue := reflect.ValueOf(num)
	numType := reflect.TypeOf(num)
	fmt.Printf("num value: %v, type: %v\n", numValue, numType)

	// 2. 结构体的反射
	p := Person{Name: "Alice", Age: 25}
	pValue := reflect.ValueOf(p)
	pType := reflect.TypeOf(p)
	fmt.Printf("p value: %v, type: %v\n", pValue, pType)

	// 遍历结构体字段
	for i := range pType.NumField() {
		filed := pType.Field(i)
		value := pValue.Field(i)

		// 下面两种打印值方法的区别
		// value.interface()
		// 将反射值，转换回接口类型，返回原始的、实际的值
		// 可以用于后续的类型断言
		// 性能稍差一些，因为设计类型转换

		// value
		// 直接使用reflect.value  类型的值
		// 保持反射状态
		// 性能较好，因为不需要类型转换

		// 在特定的打印场景中，两种方式的输出结果是相同的，因为printf 会在处理reflect value类型时 自动调用其
		// string() 方法，将其转换为字符串。
		fmt.Printf("字段名：%s, 字段类型：%s, 字段值：%v\n", filed.Name, filed.Type, value.Interface())
		fmt.Printf("字段名：%s, 字段类型：%s, 字段值：%v\n", filed.Name, filed.Type, value)
	}

	// 3. 方法的反射
	method := pValue.MethodByName("SayHello")
	args := []reflect.Value{reflect.ValueOf("I'm a gopher")}
	method.Call(args)

	// 4. 修改值
	str := "hello world"
	strPtr := reflect.ValueOf(&str)

	// Elm 是解引用
	strPtr.Elem().SetString("goodbye world")
	fmt.Printf("str: %s\n", str)
}
