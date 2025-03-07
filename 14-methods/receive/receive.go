package receive

import "fmt"

type Person struct {
	Name  string //大写字母开头 ，可导出字段
	Age   int    // 可导出字段
	count int    // 小写字母开头，不可导出字段
}

func (p *Person) SetName(name string) {
	p.Name = name
}

func (p *Person) SetAge(age int) {
	p.Age = age
}

func (p *Person) SetCount(count int) {
	p.count = count
}

func (p *Person) GetCount() int {
	return p.count
}

func init() {
	initMethods()
}

func initMethods() {
	// 调用方法
	p := &Person{}
	p.SetName("Tom")
	p.SetAge(25)
	p.SetCount(100)
	count := p.GetCount()
	fmt.Printf("Init Methods: %s, %d, %d\n", p.Name, p.Age, count)
}

func GeneratePerson() (person *Person) {
	// p := new(Person)
	// p.Name = "Tom"
	// p.Age = 25
	// p.count = 100

	p := &Person{
		Name:  "Alice",
		Age:   20,
		count: 10,
	}
	person = p
	return
}
