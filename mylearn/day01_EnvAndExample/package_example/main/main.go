package main

import(
	"fmt"
	"mylearn/day01_环境搭建&实例演练/package_example/calc"
)

func main() {
	sum := calc.Add(100, 300)
	sub := calc.Sub(100, 300)

	fmt.Println("sum=",sum)
	fmt.Println("sub=", sub)
}