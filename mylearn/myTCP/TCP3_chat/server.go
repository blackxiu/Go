package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

//定义用户结构体
type Client struct {
	C    chan string //用户发送数据的管道
	Name string      //用户名
	Addr string      //网络地址
}

//保存在线用户  cliAddr===>Client
var onlineMap map[string]Client

//通信的管道
var message = make(chan string)

//广播用户在线的函数
func MakeMsg(cli Client, msg string) (buf string) {
	buf = "[" + cli.Addr + "]" + cli.Name + ": " + msg
	return
}

//处理用户链接
func HandleConn(conn net.Conn) {
	//获取客户端的网络地址
	cliAddr := conn.RemoteAddr().String()

	//创建一个结构体,默认用户名和网络地址一样
	cli := Client{make(chan string), cliAddr, cliAddr}

	//把结构体加入到map
	onlineMap[cliAddr] = cli

	//新开一个协程  专门给当前客户端发送信息
	go WriteMsgToClient(cli, conn)

	//广播某个客户在线
	//browser <- "["+cli.Addr+"]"+cli.Name+": login"
	message <- MakeMsg(cli, "login")

	//提示我是谁
	cli.C <- MakeMsg(cli, "I am here")

	isQuit := make(chan bool)  //判断客户是否自动退出
	hasData := make(chan bool) //判断客户是否有数据是否断连
	//新开一个协程  接收并广播用户发来的数据
	go func() {
		buf := make([]byte, 2048)

		for {
			n, err := conn.Read(buf)
			if n == 0 { //有两种情况   对方断开或者网络出问题
				isQuit <- true //判断客户是否自动退出
				fmt.Println("conn.Read err", err)
				return
			}

			msg := string(buf[:n-1]) //如果用nc测试,会多一个换行符  因此用n-1    否则用n

			if len(msg) == 3 && msg == "who" {
				//便利map  给当前用户发送所有成员
				conn.Write([]byte("user list:\n"))
				for _, tmp := range onlineMap {
					msg = tmp.Addr + ":" + tmp.Name + "\n"
					conn.Write([]byte(msg))
				}

			} else if len(msg) >= 8 && msg[:6] == "rename" {
				//重命名
				name := strings.Split(msg, "|")[1]
				cli.Name = name
				onlineMap[cliAddr] = cli
				conn.Write([]byte("rename ok\n"))
			} else { //转发此内容
				message <- MakeMsg(cli, msg)
			}
			hasData <- true //代表有数据 没有断连
		}
	}()

	for {
		//通过selec检测channel的流动
		select {
		case <-isQuit:
			delete(onlineMap, cliAddr)           //把当前用户从map移除
			message <- MakeMsg(cli, "login out") //广播谁下线了
		case <-hasData: //有数据没断连什么都不做
		case <-time.After(60 * time.Second): //60秒后超时断开连接
			delete(onlineMap, cliAddr)                //把当前用户从map移除
			message <- MakeMsg(cli, "time out leave") //广播谁下线了
			return
		}
	}
}

//转发消息
//只要有消息进来,遍历Map,给每个Map成员发送此消息
func Manageer() {

	onlineMap = make(map[string]Client) //给map分配空间

	for {
		msg := <-message //没有消息前,这里会阻塞死循环

		//一旦有消息进来
		//遍历Map,给每个Map成员发送此消息
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}

//给客户端发信息
func WriteMsgToClient(cli Client, conn net.Conn) {
	for msg := range cli.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func main() {

	//监听
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer listener.Close()

	//新开一个协程,用来转发消息
	//只要有消息进来,遍历Map,给每个Map成员发送此消息
	go Manageer()

	//主协程  死循环阻塞等待用户链接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err =", err)
			continue
		}

		go HandleConn(conn) //处理用户链接
	}

}
