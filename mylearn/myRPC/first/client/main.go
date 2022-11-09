package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"mylearn/myRPC/first/pb"
	"os"
	"time"
)

//使用端口   键盘不输入名字时 默认打印的是Hello World
const (
	address     = "localhost:50051"
	defaultName = "World"
)

func main() {
	//建立连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	//链接服务 并输出返回的响应
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	//1s的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet : %v", err)
	}
	log.Printf("Greeting : %s", r.Message)
}
