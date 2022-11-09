package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Threads    int             // threads
	CacheRedis cache.CacheConf // redis 缓存
	Endpoints  nodesBalancing  // endpoint 配置
}

type nodesBalancing struct {
	Nodes []nodeWithWeight
	Type  int
}

type nodeWithWeight struct {
	node
	Weight int `json:""`
}

type node struct {
	Name string `json:""`
	Url  string `json:""`
}
