package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	pb "mylearn/myRPC/message/pb"
	"net/http"
)

func main() {
	//建立连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	//建立一个http服务
	r := gin.Default()
	r.GET("/rest/n/:name", func(c *gin.Context) {
		name := c.Param("name")

		//链接服务 并输出返回的响应
		req := &pb.HelloRequest{Name: name}
		res, err := client.SayHello(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res.Message),
		})
	})

	if err := r.Run(":8052"); err != nil {
		log.Fatalf("could not run server : %v", err)
	}
}

//浏览器网址访问 http://localhost:8052/rest/n/***   ***就是输入的参数name
