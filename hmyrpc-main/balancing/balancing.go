package balancing

import "github.com/shenzhendev/hmyrpc/rpc"

//LBType is
type LBType int

const (
	LBWeight         LBType = 1
	LBConsistentHash LBType = 2
)

type LoadBalance interface {
	Add(endpoint rpc.RPCClient, params ...string) error
	Get(key string) rpc.RPCClient
	Mark(endpoint rpc.RPCClient, up bool) error
}

func LoadBalanceFactory(lbType LBType) LoadBalance {
	switch lbType {
	case LBWeight:
		return NewWeightBalance()
	case LBConsistentHash:
		return NewConsistentHashBalance(15, nil)
	default:
		return NewWeightBalance()
	}
}
