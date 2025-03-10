package main

import (
	"fmt"
	"time"
)

func TestSelect() {

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	defer func() {
		close(ch1)
		close(ch2)
		close(ch3)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("ch1 done")

		ch1 <- 10

	}()

	select {
	case n, ok := <-ch1:
		fmt.Println(n, ok)
	case n, ok := <-ch2:
		fmt.Println(n, ok)
	case n, ok := <-ch3:
		fmt.Println(n, ok)
		// default:
		// 	fmt.Println("default")
	}
}
