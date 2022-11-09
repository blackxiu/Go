package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// Arith 算数运算结构体
type Arith struct {
}

// ArithRequest 算数运算请求结构体
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

//乘法运算方法
func (this *Arith) Multiply(req ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

//除法运算方法
func (this *Arith) Divide(req ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by 0")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func main() {
	rpc.Register(new(Arith)) //注册RPC服务
	rpc.HandleHTTP()         //采用HTTP协议作为RPC载体

	lis, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}

	fmt.Fprintf(os.Stdout, "%s", "start connection")

	for {
		conn, err := lis.Accept() //接受客户端连接请求
		if err != nil {
			log.Fatalln("fatal error: ", err)
		}

		go func(conn net.Conn) { //并发处理客户端请求
			fmt.Fprintf(os.Stdout, "%s", "new client in coming\n")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
