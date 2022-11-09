package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/codeskyblue/groupcache"
)


//缓存管理    groupcache包
//创建一个group（一个group是一个存储模块，类似命名空间，可以创建多个）“thumbNails”。
////NewGroup参数分别是group名字、group大小byte、 getter函数
//当获取不到key对应的缓存的时候，getter函数处理如何获取相应数据，并设置给dest，然后thumbNails_cache便会缓存key对应的数据
var thumbNails = groupcache.NewGroup("thumbnail", 512<<20, groupcache.GetterFunc( //512<<20表示512 2^20B也就是512MB
	func(ctx groupcache.Context, key string, dest groupcache.Sink) error {//缓存未命中时从源获取数据
		fileName := key
		bytes, err := generateThumbnail(fileName)
		if err != nil {
			return err
		}
		dest.SetBytes(bytes)
		return nil
	}))

//提供了HTTP客户端和服务端的实现
func generateThumbnail(key string) ([]byte, error) {
	//Parse函数解析rawurl为一个URL结构体，把url结构体返回给u
	u, _ := url.Parse(*mirror)
	u.Path = key

	//程序在使用完回复后必须关闭回复的主体。
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// FileHandler 文件处理逻辑
func FileHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path

	state.addActiveDownload(1)
	defer state.addActiveDownload(-1)

	//如果upstream为空,重定向
	if *upstream == "" { // Master
		//如果节点地址为空,重新查询节点地址,重定向
		if slaveAddr, err := slaveMap.PeekSlave(); err == nil {
			u, _ := url.Parse(slaveAddr)
			u.Path = r.URL.Path
			u.RawQuery = r.URL.RawQuery

			//func Redirect(w ResponseWriter, r *Request, urlStr string, code int)
			//Redirect回复请求一个重定向地址urlStr和状态码code。该重定向地址可以是相对于请求r的相对地址。
			http.Redirect(w, r, u.String(), 302)
			return
		}
	}
	fmt.Println("KEY:", key)
	var data []byte
	var ctx groupcache.Context
	//这里还没看懂
	err := thumbNails.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var modTime time.Time = time.Now()

	rd := bytes.NewReader(data)
	http.ServeContent(w, r, filepath.Base(key), modTime, rd)
}

//String用指定的名称、默认值、使用信息注册一个string类型flag。返回一个保存了该flag的值的指针。
var (
	mirror   = flag.String("mirror", "", "Mirror Web Base URL")
	logfile  = flag.String("log", "-", "Set log file, default STDOUT")
	upstream = flag.String("upstream", "", "Server base URL, conflict with -mirror")
	address  = flag.String("addr", ":5000", "Listen address")
	token    = flag.String("token", "1234567890ABCDEFG", "slave and master token should be same")
)

// InitSignal 初始化信号
func InitSignal() {
	//建立管道sig  用于接收发送os.Signal类型的管道  缓存容量为2个int
	sig := make(chan os.Signal, 2)

	//func Notify(c chan<- os.Signal, sig ...os.Signal)
	//Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {   //开一个协程
		for {
			s := <-sig  //传输控制段的指令  比如停止等命令
			fmt.Println("Got signal:", s)
			if state.Closed { //冷断开
				fmt.Println("Cold close !!!")
				os.Exit(1)
			}
			fmt.Println("Warm close, waiting ...") //热断开
			go func() {  //新开一个协程 确认断开的状态
				state.Close()
				os.Exit(0)
			}()
		}
	}()
}

func main() {
	//从os.Args[1:]中解析注册的flag。
	flag.Parse()

	//设置错误信息
	if *mirror != "" && *upstream != "" {
		log.Fatal("Can't set both -mirror and -upstream")
		//Fatal等价于{l.Print(v...); os.Exit(1)} 输出然后退出
	}
	if *mirror == "" && *upstream == "" {
		log.Fatal("Must set one of -mirror and -upstream")
	}
	if *upstream != "" {
		//先赋值err,然后判断err是否为空,如果err不为空,输出.
		if err := InitSlave(); err != nil {
			log.Fatal(err)
		}
	}
	if *mirror != "" {
		if _, err := url.Parse(*mirror); err != nil {
			log.Fatal(err)
		}
		if err := InitMaster(); err != nil {
			log.Fatal(err)
		}
	}


	InitSignal()  //接收信息并打印出收到信号
	fmt.Println("Hello CDN")

	//HandleFunc注册一个处理器函数handler和对应的模式pattern
	// 这里设置路由，groupcache会自动解析该路由为groupname
	//文件处理逻辑
	http.HandleFunc("/", FileHandler)

	//打印出监听的端口
	log.Printf("Listening on %s", *address)

	//建立tcp连接并阻塞
	log.Fatal(http.ListenAndServe(*address, nil))
}
