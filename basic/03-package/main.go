package main

import (
	h "example/hello"
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	h.SayHello()
}

func init() {
	fmt.Println("Main package init() is called")
}
