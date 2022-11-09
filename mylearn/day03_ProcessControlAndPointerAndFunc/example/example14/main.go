package main

import "fmt"

type opFunc func(int, int) int

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func operator(op opFunc, a, b int) int {

	return op(a, b)
}

func main() {
	var a, b int
	add(a, b)

	var c opFunc
	c = add
	fmt.Println(add)
	fmt.Println(c)

	sum := operator(c, 100, 200)
	fmt.Println(sum)
}
