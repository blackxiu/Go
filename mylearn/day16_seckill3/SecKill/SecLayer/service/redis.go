package service

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"

	"crypto/md5"
	"math/rand"
)

func initRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.RedisAddr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}
	return
}

func initRedis(conf *SecLayerConf) (err error) {

	secLayerContext.proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init proxy2layer redis pool failed, err:%v", err)
		return
	}

	secLayerContext.layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis pool failed, err:%v", err)
		return
	}

	return
}

func RunProcess() (err error) {

	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleReader()
	}

	for i := 0; i < secLayerContext.secLayerConf.WriteGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleWrite()
	}

	for i := 0; i < secLayerContext.secLayerConf.HandleUserGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleUser()
	}

	logs.Debug("all process goroutine started")
	secLayerContext.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func HandleReader() {

	logs.Debug("read goroutine running")
	for {
		conn := secLayerContext.proxy2LayerRedisPool.Get()
		for {
			data, err := redis.String(conn.Do("blpop", secLayerContext.secLayerConf.Proxy2LayerRedis.RedisQueueName, 0))
			if err != nil {
				logs.Error("pop from queue failed, err:%v", err)
				break
			}

			logs.Debug("pop from queue, data:%s", data)

			var req SecRequest
			err = json.Unmarshal([]byte(data), &req)
			if err != nil {
				logs.Error("unmarshal to secrequest failed, err:%v", err)
				continue
			}

			now := time.Now().Unix()
			if now-req.AccessTime.Unix() >= int64(secLayerContext.secLayerConf.MaxRequestWaitTimeout) {
				logs.Warn("req[%v] is expire", req)
				continue
			}

			timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToHandleChanTimeout))
			select {
			case secLayerContext.Read2HandleChan <- &req:
			case <-timer.C:
				logs.Warn("send to handle chan timeout, req:%v", req)
				break
			}
		}

		conn.Close()
	}
}

func HandleWrite() {
	logs.Debug("handle write running")

	for res := range secLayerContext.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send to redis, err:%v, res:%v", err, res)
			continue
		}
	}
}

func sendToRedis(res *SecResponse) (err error) {

	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("marshal failed, err:%v", err)
		return
	}

	conn := secLayerContext.layer2ProxyRedisPool.Get()
	_, err = redis.String(conn.Do("rpush", secLayerContext.secLayerConf.Layer2ProxyRedis.RedisQueueName, string(data)))
	if err != nil {
		logs.Warn("rpush to redis failed, err:%v", err)
		return
	}

	return
}

func HandleUser() {

	logs.Debug("handle user running")
	for req := range secLayerContext.Read2HandleChan {
		logs.Debug("begin process request:%v", req)
		res, err := HandleSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed, err:%v", err)
			res = &SecResponse{
				Code: ErrServiceBusy,
			}
		}

		timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToWriteChanTimeout))
		select {
		case secLayerContext.Handle2WriteChan <- res:
		case <-timer.C:
			logs.Warn("send to response chan timeout, res:%v", res)
			break
		}

	}
	return
}

func HandleSecKill(req *SecRequest) (res *SecResponse, err error) {

	secLayerContext.RWSecProductLock.RLock()
	defer secLayerContext.RWSecProductLock.Unlock()

	res = &SecResponse{}
	product, ok := secLayerContext.secLayerConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("not found product:%v", req.ProductId)
		res.Code = ErrNotFoundProduct
		return
	}

	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldout
		return
	}

	now := time.Now().Unix()
	alreadySoldCount := product.secLimit.Check(now)
	if alreadySoldCount >= product.SoldMaxLimit {
		res.Code = ErrRetry
		return
	}

	secLayerContext.HistoryMapLock.Lock()
	userHistory, ok := secLayerContext.HistoryMap[req.UserId]
	if !ok {
		userHistory = &UserBuyHistory{
			history: make(map[int]int, 16),
		}

		secLayerContext.HistoryMap[req.UserId] = userHistory
	}

	histryCount := userHistory.GetProductBuyCount(req.ProductId)
	secLayerContext.HistoryMapLock.Unlock()

	if histryCount >= product.OnePersonBuyLimit {
		res.Code = ErrAlreadyBuy
		return
	}

	curSoldCount := secLayerContext.productCountMgr.Count(req.ProductId)
	if curSoldCount >= product.Total {
		res.Code = ErrSoldout
		product.Status = ProductStatusSoldout
		return
	}

	curRate := rand.Float64()
	if curRate > product.BuyRate {
		res.Code = ErrRetry
		return
	}

	userHistory.Add(req.ProductId, 1)
	secLayerContext.productCountMgr.Add(req.ProductId, 1)

	//??????id&??????id&????????????&??????
	res.Code = ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, now, secLayerContext.secLayerConf.TokenPasswd)

	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now

	return
}
