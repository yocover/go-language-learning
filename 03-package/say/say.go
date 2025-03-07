package say

func init() {
	println("Say package init")
	SayHello()
}

func SayHello() string {
	println("Say package!")
	return "Say package!"
}
