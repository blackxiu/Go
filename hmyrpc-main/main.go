package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shenzhendev/hmyrpc/balancing"
	"github.com/shenzhendev/hmyrpc/cache"
	"github.com/shenzhendev/hmyrpc/internal/config"
	"github.com/shenzhendev/hmyrpc/internal/handler"
	"github.com/shenzhendev/hmyrpc/internal/svc"
	"github.com/shenzhendev/hmyrpc/rpc"
	"github.com/shenzhendev/hmyrpc/rpc/hmy"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "./etc/api.yaml", "the config file")

type Speeder struct {
	config    config.Config
	endpoints []rpc.RPCClient
	balancing balancing.LoadBalance
	cache     *cache.RedisClient
	http      *rest.Server
	lock      sync.Mutex
	quit      chan struct{}
}

func NewSpeeder(config config.Config) *Speeder {
	speeder := &Speeder{
		quit: make(chan struct{}),
	}
	// config
	speeder.config = config
	// for redis cache client 缓存配置
	speeder.cache = cache.NewRedisClient(&cache.Config{
		Endpoint: config.CacheRedis[0].Host,
		Password: config.CacheRedis[0].Pass,
		Database: 0,
		PoolSize: 20,
	}, "speeder", 2*time.Second)
	pong, err := speeder.cache.Check()
	if err != nil {
		logx.Errorf("Can't establish connection to backend: %v", err)
	} else {
		logx.Infof("Cache check reply: %v", pong)
	}

	// load balancing负载均衡
	speeder.balancing = balancing.LoadBalanceFactory(balancing.LBType(config.Endpoints.Type))

	// for endpoints 端点配置
	for _, endpoint := range config.Endpoints.Nodes {
		// new rpc client
		rpcClient := hmy.NewHarmonyRPCClient(endpoint.Name, endpoint.Url, "10m")
		// add to endpoints加入到端点中
		speeder.endpoints = append(speeder.endpoints, rpcClient)
		// add to load balancing加入到负载均衡中
		err := speeder.balancing.Add(rpcClient, strconv.Itoa(endpoint.Weight))
		if err != nil {
			log.Fatalf("error add [ %s ] to load balancing, error: [ %s ]", "rpcClient", err)
		}
	}

	// http server
	server := rest.MustNewServer(config.RestConf)
	speeder.http = server
	// return
	return speeder
}

func (s *Speeder) Wait() {
	defer func() {
		_ = s.cache.Close()
		logx.Infof("%s: %s, %s", "cache", "stopped", "leaving")
	}()
	<-s.quit
}

func (s *Speeder) StartHttpServer(balance balancing.LoadBalance) {
	defer s.http.Stop()
	// new context
	ctx := svc.NewServiceContext(s.config, s.cache, balance)
	// register handler
	handler.RegisterHandlers(s.http, ctx)
	//
	logx.Infof("Starting server at %s:%d...\n", s.config.Host, s.config.Port)
	// start
	s.http.Start()
}

func (s *Speeder) Next(key string) (endpoint rpc.RPCClient) {
	if endpoint = s.balancing.Get(key); endpoint == nil {
		logx.Errorf("Get error: %s ", "no endpoint picked")
	}
	return endpoint
}

func (s *Speeder) MarkUp(endpoint rpc.RPCClient) error {
	return s.balancing.Mark(endpoint, true)
}

func (s *Speeder) MarkDown(endpoint rpc.RPCClient) error {
	return s.balancing.Mark(endpoint, false)
}

func main() {
	// parse config解析配置
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// runtime运行时的预处理
	if c.Threads > 0 { //配置线程和cpu数量  gomaxprocs功能是设置可执行的最大cpu数
		runtime.GOMAXPROCS(c.Threads) //默认为conf中的值
		logx.Infof("Running with %v threads", c.Threads)
	}

	// new speedster instance新建speeder
	speeder := NewSpeeder(c)

	go speeder.StartHttpServer(speeder.balancing) // run api server

	// wait for os terminationw
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Kill, os.Interrupt)
	go func(signal chan os.Signal) {
		<-signal
		speeder.quit <- struct{}{}
	}(cancel)

	//go func() {
	//	for i := 0; i < 50; i++ {
	//		fmt.Println(
	//			speeder.balancing.Get("hmyv2_getStakingTransactionsCount").Url(),
	//		)
	//	}
	//
	//}()

	// wait termination
	speeder.Wait()
}
