package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
	fmt.Println("Hello, Golang!")

	var name string = "Golang-1111"
	fmt.Println("Hello, " + name + "!")

	value := 10
	fmt.Println(value)

	fmt.Println("Hello, Golang!", "!")

	fmt.Printf("Hello, %v!\n", "Golang")

	var myName string = "Golang"
	var age int = 25

	fmt.Printf("Hello %s, %d years old!\n", myName, age)

	fmt.Printf("%-15v $%4v\n", "Product1", 100.99)
	fmt.Printf("%-15v $%4.2f\n", "Product2", 100.99)

}

// The new version adds a new line of code to print "Hello, Golang!" to the console.
