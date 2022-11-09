package svc

import (
	"github.com/shenzhendev/hmyrpc/balancing"
	"github.com/shenzhendev/hmyrpc/cache"
	"github.com/shenzhendev/hmyrpc/internal/config"
)

type ServiceContext struct {
	Config   config.Config
	Cache    *cache.RedisClient
	Endpoint balancing.LoadBalance
}

func NewServiceContext(c config.Config, cache *cache.RedisClient, balance balancing.LoadBalance) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Cache:    cache,
		Endpoint: balance,
	}
}
