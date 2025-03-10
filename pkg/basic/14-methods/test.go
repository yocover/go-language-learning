package main

import (
	"fmt"
	"methods/receive"
)

func init() {
	testFunc()
}

func testFunc() {

	person := receive.Person{
		Name: "John",
		Age:  30,
	}

	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
