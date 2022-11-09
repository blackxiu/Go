package rpc

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/shenzhendev/hmyrpc/util"
	"github.com/valyala/fasthttp"
)

type RPCClient interface {
	Request(method string, params interface{}) (rpcResp *JSONRpcResp, err error)
	Url() string
	Rate() (sickRate int, successRate int)
}

type BaseClient struct {
	sync.RWMutex
	url         string
	sickRate    int
	successRate int
	client      *fasthttp.Client
	id          int
}

func (r *BaseClient) Url() string {
	return r.url
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func NewBaseClient(name, url, timeout string) *BaseClient {
	rpcClient := &BaseClient{url: url, id: 0}
	timeoutIntv := util.MustParseDuration(timeout)
	rpcClient.client = &fasthttp.Client{
		ReadTimeout:         timeoutIntv,
		WriteTimeout:        timeoutIntv,
		MaxIdleConnDuration: 5 * time.Minute,
		Name:                name,
	}
	return rpcClient
}

func DefaultClient(name, url, timeout string) RPCClient {
	return NewBaseClient(name, url, timeout)
}

func (r *BaseClient) rawRequest(method string, params interface{}) (result []byte, err error) {
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": r.id}
	fmt.Println("============ jsonreq", jsonReq)
	data, _ := json.Marshal(jsonReq)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(data)) //
	// use fasthttp request
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBody(data)
	req.SetRequestURI(r.url)
	// init response
	res := fasthttp.AcquireResponse()
	err = r.client.DoTimeout(req, res, 10*time.Second)
	if err != nil {
		r.markSick()
		return nil, err
	}
	// release response
	fasthttp.ReleaseRequest(req)
	//
	result = make([]byte, len(res.Body()))
	// deep copy
	copy(result, res.Body()) // 如果此处不copy，会引发goroutine下的panic（res被释放了）
	// release response
	fasthttp.ReleaseResponse(res)
	// lock
	r.Lock()
	r.id++
	r.Unlock()
	// mark
	r.markAlive()
	return result, err
}

func (r *BaseClient) Request(method string, params interface{}) (rpcResp *JSONRpcResp, err error) {
	raw, err := r.rawRequest(method, params)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(raw, &rpcResp)
	return rpcResp, err
}

func (r *BaseClient) Rate() (sickRate int, successRate int) {
	r.RLock()
	defer r.RUnlock()
	return r.sickRate, r.successRate
}

func (r *BaseClient) markSick() {
	r.Lock()
	r.sickRate++
	r.Unlock()
}

func (r *BaseClient) markAlive() {
	r.Lock()
	r.successRate++
	r.Unlock()
}
