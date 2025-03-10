package main

import "fmt"

// Person 结构体用于演示方法接收者的区别
type Person struct {
	name string
	age  int
}

// 值接收者的方法
// 1. 使用值接收者时，方法内部会获得结构体的副本
// 2. 对副本的修改不会影响原始值
func (p Person) modifyByValue() {
	p.age += 1 // 这个修改只会影响副本
	fmt.Printf("在值接收者方法内: age = %d\n", p.age)
}

// 指针接收者的方法
// 1. 使用指针接收者时，方法内部会获得结构体的指针
// 2. 对指针指向的值的修改会影响原始值
func (p *Person) modifyByPointer() {
	p.age += 1 // 这个修改会影响原始值
	fmt.Printf("在指针接收者方法内: age = %d\n", p.age)
}

// 值接收者的方法 - 只读操作
func (p Person) printInfo() {
	fmt.Printf("姓名: %s, 年龄: %d\n", p.name, p.age)
}

// 指针接收者的方法 - 修改操作
func (p *Person) setAge(age int) {
	p.age = age
}

func testMethodReceiver() {
	fmt.Println("=== 测试方法接收者的区别 ===")

	// 1. 使用值类型变量
	fmt.Println("\n1. 使用值类型变量:")
	person1 := Person{name: "Alice", age: 25}

	fmt.Printf("调用方法前: age = %d\n", person1.age)
	person1.modifyByValue() // 传递的是副本
	fmt.Printf("值接收者方法后: age = %d\n", person1.age)
	person1.modifyByPointer() // Go会自动获取地址
	fmt.Printf("指针接收者方法后: age = %d\n", person1.age)

	// 2. 使用指针类型变量
	fmt.Println("\n2. 使用指针类型变量:")
	person2 := &Person{name: "Bob", age: 30}

	fmt.Printf("调用方法前: age = %d\n", person2.age)
	person2.modifyByValue() // Go会自动解引用
	fmt.Printf("值接收者方法后: age = %d\n", person2.age)
	person2.modifyByPointer() // 直接传递指针
	fmt.Printf("指针接收者方法后: age = %d\n", person2.age)

	// 3. 方法集合的区别
	fmt.Println("\n3. 方法集合的使用示例:")
	p1 := Person{name: "Charlie", age: 35}
	p2 := &Person{name: "David", age: 40}

	// 值类型变量可以调用值接收者的方法
	p1.printInfo()
	// 值类型变量也可以调用指针接收者的方法（Go会自动获取地址）
	p1.setAge(36)
	p1.printInfo()

	// 指针类型变量可以调用指针接收者的方法
	p2.setAge(41)
	// 指针类型变量也可以调用值接收者的方法（Go会自动解引用）
	p2.printInfo()
}

/* 方法接收者的重要区别：

1. 值接收者 (t Test)：
   - 方法内部使用的是值的副本
   - 不会修改原始值
   - 适合只读操作
   - 可以用值类型或指针类型调用
   - 值接收者的方法可以被接口值类型调用

2. 指针接收者 (t *Test)：
   - 方法内部使用的是值的指针
   - 可以修改原始值
   - 适合需要修改接收者的操作
   - 可以用值类型或指针类型调用
   - 指针接收者的方法只能被接口指针类型调用

3. 使用场景：
   - 需要修改接收者时，使用指针接收者
   - 接收者是大型结构体时，使用指针接收者（避免复制开销）
   - 需要保持一致性时，所有方法都使用相同的接收者类型
   - 只读操作且结构体较小时，可以使用值接收者

4. 注意事项：
   - 不能在同一个类型上定义同名的值接收者和指针接收者方法（编译错误）
   - 方法集合的规则在接口实现时很重要
   - 值类型变量可以调用指针接收者的方法（Go自动获取地址）
   - 指针类型变量可以调用值接收者的方法（Go自动解引用）
*/
