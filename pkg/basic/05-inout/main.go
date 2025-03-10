package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	fmt.Printf("%s\n", "Hello, world!")
	fmt.Printf("%q\n", "Hello, world!")
	fmt.Printf("%d\n", 20)
	fmt.Printf("%f\n", 3.1415926)
	fmt.Printf("%e\n", 12i)
	fmt.Printf("%g\n", 12i)
	fmt.Printf("%b\n", 128)
	fmt.Printf("%#b\n", 128)
	fmt.Printf("%o\n", 128)
	fmt.Printf("%#o\n", 128)
	fmt.Println("--------")
	fmt.Printf("%x\n", 128)
	fmt.Printf("%#x\n", 128)
	fmt.Printf("%#X\n", 128)
	fmt.Printf("%U\n", 'A')

	println("address -------")

	count := 10
	var num = 3.141592
	fmt.Printf("%d\n", count)
	fmt.Printf("%p\n", &count)
	fmt.Printf("%p\n", &num)

	println("struct -------")

	type person struct {
		name    string
		age     int
		address string
	}

	p := person{"Alice", 25, "Shanghai"}
	fmt.Printf("%s is %d years old and lives in %s\n", p.name, p.age, p.address)

	fmt.Printf("%T\n", p)
	fmt.Printf("%v\n", p)
	fmt.Printf("%+v\n", p)
	fmt.Printf("%#v\n", p)

	println("Scan -------")

	var s, s2 string

	fmt.Println("Enter a string:")
	fmt.Scan(&s, &s2)
	fmt.Println("You entered:", s, s2)

	var str1, str2 string
	fmt.Scanln(&str1, &str2)
	fmt.Println("You entered:", str1, str2)

	println("Bufio -------")

	var scanner = bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a string: ")
	scanner.Scan()
	fmt.Println("You enterend:", scanner.Text())
}
