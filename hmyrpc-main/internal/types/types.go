package types

import (
	rpc "github.com/INFURA/go-ethlibs/jsonrpc"
)

type Request struct {
	Chain   string     `path:"chain"`
	Version string     `path:"version"`
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	ID      uint64     `json:"id"`
	Params  rpc.Params `json:"params"`
	//jsonrpc.Request
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Error   *JsonError  `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
