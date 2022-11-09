package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"mylearn/myRPC/first/pb"
	"net"
)

//监听端口
const (
	port = ":50051"
)

//服务对象
type server struct {
	pb.UnimplementedGreeterServer
}

//SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello！  " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//新起一个服务
	s := grpc.NewServer()

	//注册反射服务  这个服务是CLI使用的 跟服务本身没有关系
	pb.RegisterGreeterServer(s, &server{})

	fmt.Printf("开始提供服务")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
