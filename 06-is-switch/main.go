package main

func main() {
	a, b := 10, 20

	println("a is", a, "and b is", b)
	if a == b {
		println("a is equal to b")
	} else if a > b {
		println("a is greater than b")
	} else {
		println("b is greater than a")
	}

	println("Switch Example")

	var e = 100

	switch e {
	case 10:
		println("e is equal to 10")
	case 20:
		println("e is equal to 20")
	default:
		println("e is not equal to 10 or 20")
	}

}
