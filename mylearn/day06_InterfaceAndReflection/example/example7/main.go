package main

import "fmt"

func main() {

	var link Link
	for i := 0; i < 10; i++ {
		//link.InsertHead(fmt.Sprintf("str %d", i))
		link.InsertTail(fmt.Sprintf("str %d", i))
	}

	link.Trans()
}
