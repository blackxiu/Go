package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

//发送文件内容的函数
func SendFile(path string, conn net.Conn) {
	//以只读方式打开文件
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("os.Open err = ", err)
		return
	}
	defer f.Close()

	//读取文件内容,读多少就给receiver发多少
	buf := make([]byte, 1024*4)
	for {
		n, err := f.Read(buf) //读取文件内容到n
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件发送完毕")
			} else {
				fmt.Println("f.Read err = ", err)
			}

			return
		}
		//发送内容
		conn.Write(buf[:n]) //发送内容
	}
}

func main() {
	//提示输入文件
	fmt.Println("请输入文件:")
	var path string
	fmt.Scan(&path)

	//获取文件名info.Name()
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("os.Stat err = ", err)
		return
	}

	//主动连接服务器
	conn, err1 := net.Dial("tcp", "127.0.0.1:8000")
	if err1 != nil {
		fmt.Println("net.Dial err =", err1)
		return
	}

	//关闭连接
	defer conn.Close()

	//给接收方先发送文件名
	_, err = conn.Write([]byte(info.Name()))
	if err != nil {
		fmt.Println("conn.Write err =", err)
		return
	}

	//接受对方的回复
	//如果回复ok,说明对方准备就绪,可以发送文件
	var n int
	buf := make([]byte, 1024)

	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err =", err)
		return
	}

	if "ok" == string(buf[:n]) {
		//发送文件内容
		SendFile(path, conn)

	}
}
