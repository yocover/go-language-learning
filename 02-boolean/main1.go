package main

import (
	"fmt"
	"strings"
)

func main1() {

	count := 10

	var he bool = false

	he = count >= 10
	if he {
		fmt.Println("he is true")
	}

	fmt.Println("he =", he)

	fmt.Println("true && true =")

	var commond = "walk outside"

	var exit = strings.Contains(commond, "outside")

	fmt.Println("You leave the house:", exit)

}
