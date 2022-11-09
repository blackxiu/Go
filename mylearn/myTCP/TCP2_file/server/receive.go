package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func RecvFile(fileName string, conn net.Conn) {

	//新建文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Print("os.Create err = ", err)
		return
	}

	//接受多少就向新文件里写多少
	buf := make([]byte, 1024*4)
	for {
		n, err := conn.Read(buf) //接受文件的内容
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件接收完毕")
			} else {
				fmt.Println("conn.Read err = ", err)
			}
			return
		}

		if n == 0 {
			fmt.Println("n==0 文件接收完毕")
			break
		}

		//接收内容
		f.Write(buf[:n]) //往新文件里写内容
	}
}

func main() {
	//监听来自send方的请求
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Print("net.Listen err = ", err)
		return
	}
	defer listener.Close()

	//阻塞等待用户链接
	conn, err1 := listener.Accept()
	if err1 != nil {
		fmt.Print("listener.Accept err = ", err1)
		return
	}

	//接收到send方发来的第一次链接请求
	//并读取文件名
	buf := make([]byte, 1024)
	var n int
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Print("conn.Read err = ", err)
		return
	}

	fileName := string(buf[:n]) //保存文件名

	//给send方回复ok
	conn.Write([]byte("ok"))

	//准备接收文件内容
	RecvFile(fileName, conn)

}
