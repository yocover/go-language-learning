package main

import (
	"fmt"
)

// 泛型就是为了解决执行逻辑与类型无关的问题
// 如果一个问题需要根据不同类型做出不同的逻辑
// 那么根本就不应该使用泛型
// 而应该使用interface 或者 any

// 1. 泛型函数
func Sum[T int | float64](a, b T) T {
	return a + b
}

// 2. 泛型结构
type GenericSlice[T int | int32 | int64] []T

func (s GenericSlice[T]) Sum() T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// 3. 泛型哈希表
type GenericMap[k comparable, v int | string | byte] map[k]v

func (m GenericMap[k, v]) Get(key k) (v, bool) {
	val, ok := m[key]
	return val, ok
}

// 泛型结构体
type GeneraicStruct[T int | string] struct {
	Name string
	Age  T
}

// 泛型结构推荐的写法
type Company[T int | string, S int | string] struct {
	Name string
	Id   T
	Data S
}

// 泛型接口
type SayAble[T int | string] interface {
	Say() T
}

type Person[T int | string] struct {
	msg T
}

func (p Person[T]) Say() T {
	return p.msg
}

// 类型集合
// 含有类型集的接口，称为General interfaces 既通用接口
// 类型集主要用于类型约束，不能用做类型声明
// 既然是集合，就会有空集，并集，交集等概念

// 并集
type SingedInt interface {
	int | int32 | int64
}

type Integer interface {
	int | int32 | int64 | uint | uint32 | uint64
}

// 交集
type Number interface {
	SingedInt
	Integer
}

func SumNumber[T Number](a, b T) T {
	return a + b
}

func main() {

	sum2 := Sum[int](10, 20.0)
	fmt.Println(sum2)

	s := GenericSlice[int64]{1, 2, 3, 4, 5}
	fmt.Println(s.Sum())

	m := GenericMap[string, int]{"a": 1, "b": 2, "c": 3}
	val, ok := m.Get("a")
	fmt.Println(val, ok)

	// 泛型结构体
	s1 := GeneraicStruct[int]{"Alice", 25}
	s2 := GeneraicStruct[string]{"Bob", "30"}
	fmt.Println(s1.Name, s1.Age)
	fmt.Println(s2.Name, s2.Age)

	// 推荐的泛型结构体使用方法
	c1 := Company[int, string]{"Google", 123456, "Some data"}
	c2 := Company[string, int]{"Apple", "123456789", 1000000}
	fmt.Println(c1.Name, c1.Id, c1.Data)
	fmt.Println(c2.Name, c2.Id, c2.Data)

	// 泛型接口的使用方式
	var p1 SayAble[int] = Person[int]{100}
	var p2 SayAble[string] = Person[string]{"Hello"}
	fmt.Println(p1.Say())
	fmt.Println(p2.Say())

	// 类型约束
	// 接口你类型不能直接用做变量类型，只能用做约束
	// 类型约束可以让代码更加灵活，更加易读，更加安全
	// var a Number = 10
	// var b Number = 20
	fmt.Println(SumNumber(10, 20))

	// 大多数情况下可以省略类型参数，让编译器自动推导
	fmt.Println(SumNumber[int](10, 20))
}
