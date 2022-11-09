package main

import (
	"fmt"
	"math/rand"
	"mylearn/day07_FileAndjsonAndErr/example/example1/balance"
	"os"
	"time"
)

func main() {

	var insts []*balance.Instance //主机列表
	for i := 0; i < 16; i++ {
		host := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
		one := balance.NewInstance(host, 8080)
		insts = append(insts, one)
	}

	//var balanceName = "roundrobin"    //轮询
	var balanceName = "random" //随机
	//var balanceName = "hash" //一致性哈希
	if len(os.Args) > 1 {
		balanceName = os.Args[1]
	}

	for {
		inst, err := balance.DoBalance(balanceName, insts)
		if err != nil {
			fmt.Fprintf(os.Stdout, "do balance err\n")
			continue
		}
		fmt.Println(inst)
		time.Sleep(time.Second)
	}

}
