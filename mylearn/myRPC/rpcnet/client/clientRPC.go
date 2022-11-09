package main

import (
	"fmt"
	"log"
	"net/rpc"
)

//算数运算请求结构体
type ArithRequest struct {
	A int
	B int
}

////算数运算响应结构体
type ArithResponse struct {
	Pro int //乘积
	Quo int //商
	Rem int //余数
}

func main() {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8090") //拨号启动RPC服务
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	req := ArithRequest{9, 2}
	var res ArithResponse

	err = conn.Call("Arith.Multiply", req, &res) //乘法运算
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

	err = conn.Call("Arith.Divide", req, &res) //除法运算
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d = %d …… %d \n", req.A, req.B, res.Quo, res.Rem)
}
