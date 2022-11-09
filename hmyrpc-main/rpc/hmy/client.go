package hmy

import (
	"context"
	"fmt"
	"strings"

	"github.com/shenzhendev/hmyrpc/rpc"
)

type HmyRpcClient struct {
	client *rpc.BaseClient
	prefix map[string]string
}

func (h *HmyRpcClient) Url() string {
	return h.client.Url()
}

func NewHarmonyRPCClient(name, url, timeout string) rpc.RPCClient {
	client := &HmyRpcClient{
		client: rpc.NewBaseClient(name, url, timeout),
		prefix: make(map[string]string),
	}
	client.prefix["eth"] = "eth_"
	client.prefix["v1"] = "hmy_"
	client.prefix["v2"] = "hmyv2_"
	return client
}

func (h *HmyRpcClient) Request(method string, params interface{}) (rpcResp *rpc.JSONRpcResp, err error) {
	if strings.HasPrefix(method, "net_") {
		return h.client.Request(method, params)
	}
	return h.client.Request(h.prefix["v2"]+method, params)
}

func (h *HmyRpcClient) Rate() (int, int) {
	h.client.RLock()
	defer h.client.RUnlock()
	return h.client.Rate()
}

func (h *HmyRpcClient) GetBalance(ctx context.Context, address string) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getBalance", []string{address})
	if err != nil {
		return nil, err
	}
	//
	//var reply uint256.Int
	//if rpcResp.Result != nil {
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	if err != nil {
	//		return reply.ToBig(), err
	//	}
	//}
	return rpcResp, err
	//return util.String2Big(reply), err
}

func (h *HmyRpcClient) GetBalanceByBlockNumber(ctx context.Context, address string, number int64) (*rpc.JSONRpcResp, error) {
	params := []interface{}{address, number}
	rpcResp, err := h.Request("getBalanceByBlockNumber", params)
	if err != nil {
		return nil, err
	}
	//
	return rpcResp, err
	//var reply string
	//if rpcResp.Result != nil {
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//return util.String2Big(reply), err
}

func (h *HmyRpcClient) GetStakingTransactionsCount(ctx context.Context, args TxsCountArgs) (*rpc.JSONRpcResp, error) {
	// txType: type of staking transaction (SENT, RECEIVED, ALL)
	rpcResp, err := h.Request("getStakingTransactionsCount", []string{args.Address, args.Type})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	var count int64
	//	err = json.Unmarshal(*rpcResp.Result, &count)
	//	return count, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetStakingTransactionsHistory(ctx context.Context, args TxHistoryArgs) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getStakingTransactionsHistory", []interface{}{args})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	if !args.FullTx {
	//		var reply *rpc.StakingTxHistoryWithHash
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return rpc.StructuredResponse{"staking_transactions": reply.StakingTransactions}, err
	//	} else {
	//		var reply *rpc.StakingTxHistoryWithFullTx
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return rpc.StructuredResponse{"staking_transactions": reply.StakingTransactions}, err
	//	}
	//}
	//return nil, err	//return map[string]interface{}
}

func (h *HmyRpcClient) GetTransactionsCount(ctx context.Context, args TxsCountArgs) (*rpc.JSONRpcResp, error) {
	// txType: type of staking transaction (SENT, RECEIVED, ALL)
	rpcResp, err := h.Request("getTransactionsCount", []string{args.Address, args.Type})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	var count int64
	//	err = json.Unmarshal(*rpcResp.Result, &count)
	//	return count, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetTransactionsHistory(ctx context.Context, args TxHistoryArgs) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getTransactionsHistory", []interface{}{args})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	if !args.FullTx {
	//		var reply *rpc.TxHistoryWithHash
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return rpc.StructuredResponse{"transactions": reply.Transactions}, err
	//	} else {
	//		var reply *rpc.TxHistoryWithFullTx
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return rpc.StructuredResponse{"transactions": reply.Transactions}, err
	//	}
	//}
	//return nil, err
}

func (h *HmyRpcClient) GetBlocks(ctx context.Context, from, to rpc.BlockNumber, blockArgs *BlockArgs) (*rpc.JSONRpcResp, error) {
	params := []interface{}{from, to, blockArgs}
	rpcResp, err := h.Request("getBlocks", params)
	if err != nil {
		return nil, err
	}
	return rpcResp, err
	//if rpcResp.Result != nil {
	//	var reply []interface{}
	//	err := json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, blockArgs *BlockArgs) (*rpc.JSONRpcResp, error) {
	params := []interface{}{fmt.Sprintf("0x%x", number), blockArgs}
	return h.getBlockBy("getBlockByNumber", params, blockArgs.FullTx)
}

func (h *HmyRpcClient) GetBlockByHash(ctx context.Context, hash string, blockArgs *BlockArgs) (*rpc.JSONRpcResp, error) {
	params := []interface{}{hash, blockArgs}
	return h.getBlockBy("getBlockByHash", params, blockArgs.FullTx)
}

func (h *HmyRpcClient) getBlockBy(method string, params []interface{}, fullTx bool) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request(method, params)
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	if fullTx {
	//		var reply *common.BlockWithFullTx
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return reply, err
	//	} else {
	//		var reply *common.BlockWithTxHash
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return reply, err
	//	}
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetBlockSigners(ctx context.Context, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getBlockSigners", []interface{}{number})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	var reply []string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return []string{}, nil
}

func (h *HmyRpcClient) GetBlockSignersKeys(ctx context.Context, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getBlockSignerKeys", []interface{}{number})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	var reply []string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return []string{}, nil
}

func (h *HmyRpcClient) GetBlockTransactionCountByNumber(ctx context.Context, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getBlockTransactionCountByNumber", []interface{}{number})
	if err != nil {
		return nil, err
	}
	return rpcResp, nil
	//if rpcResp.Result != nil {
	//	var count int64
	//	err = json.Unmarshal(*rpcResp.Result, &count)
	//	return count, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetBlockTransactionCountByHash(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getBlockTransactionCountByHash", []interface{}{hash})
	//if err != nil {
	//	return nil, err
	//}
	return rpcResp, err
	//if rpcResp.Result != nil {
	//	var count int64
	//	err = json.Unmarshal(*rpcResp.Result, &count)
	//	return count, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getHeaderByNumber", []interface{}{number})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *rpc.HeaderInformation
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetLatestChainHeaders(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getLatestChainHeaders", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *rpc.HeaderPair
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) LatestHeader(ctx context.Context) (*rpc.JSONRpcResp, error) {
	// one
	//rpcResp, err := h.Request("latestHeader", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *HeaderInformation
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
	// other solution
	return h.GetHeaderByNumber(ctx, rpc.BlockNumber(rpc.LatestBlockNumber))

}

func (h *HmyRpcClient) BlockNumber(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("blockNumber", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return 0, err
	//}
	//if rpcResp.Result != nil {
	//	var count int64
	//	err = json.Unmarshal(*rpcResp.Result, &count)
	//	return count, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetCirculatingSupply(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getCirculatingSupply", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GetEpoch(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getEpoch", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return 0, err
	//}
	//if rpcResp.Result != nil {
	//	var reply uint64
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetEpochLastBlock(ctx context.Context, epoch uint64) (*rpc.JSONRpcResp, error) {
	return h.Request("epochLastBlock", []interface{}{epoch})
}

func (h *HmyRpcClient) GetLastCrossLinks(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getLastCrossLinks", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.CrossLink
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil //[]*common.CrossLink
}

func (h *HmyRpcClient) GetLeader(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getLeader", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GasPrice(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("gasPrice", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return 0, err
	//}
	//if rpcResp.Result != nil {
	//	var reply uint64
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetShardingStructure(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getShardingStructure", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*rpc.ShardingStructure
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetTotalSupply(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getTotalSupply", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GetValidators(ctx context.Context, epoch uint64) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getValidators", []interface{}{epoch})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//var reply *rpc.ValidatorsFields
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetValidatorKeys(ctx context.Context, epoch uint64) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getValidatorKeys", []interface{}{epoch})
	return rpcResp, err
	//if err != nil {
	//	return []string{}, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return []string{}, nil
}

func (h *HmyRpcClient) GetNodeMetadata(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getNodeMetadata", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//
	//if rpcResp.Result != nil {
	//	var reply *common.NodeMetadata
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) ProtocolVersion(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("protocolVersion", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return 0, err
	//}
	//if rpcResp.Result != nil {
	//	var reply int64
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) PeerCount(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("net_peerCount", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) Call(ctx context.Context, args rpc.CallArgs, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("call", []interface{}{args, number})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply hexutil.Bytes
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	replyStr := reply.String()
	//	return replyStr, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) EstimateGas(ctx context.Context, args rpc.CallArgs) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("estimateGas", []interface{}{args})
	return rpcResp, err
	//if err != nil {
	//	return 0, err
	//}
	//if rpcResp.Result != nil {
	//	var reply uint64
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return 0, nil
}

func (h *HmyRpcClient) GetCode(ctx context.Context, address string, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getCode", []interface{}{address, number})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply hexutil.Bytes
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	replyStr := reply.String()
	//	return replyStr, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GetStorageAt(ctx context.Context, address, key string, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getStorageAt", []interface{}{address, key, number})
	return rpcResp, err
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply hexutil.Bytes
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	replyStr := reply.String()
	//	return replyStr, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GetDelegationsByDelegator(ctx context.Context, address string) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getDelegationsByDelegator", []interface{}{address})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*rpc.Delegation
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetDelegationsByDelegatorByBlockNumber(ctx context.Context, address string, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getDelegationsByDelegatorByBlockNumber", []interface{}{address, number})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	if address.Address != nil {
	//		var reply []*rpc.Delegation
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return reply, err
	//	} else {
	//		var reply [][]*rpc.Delegation
	//		err = json.Unmarshal(*rpcResp.Result, &reply)
	//		return reply, err
	//	}
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetDelegationsByValidator(ctx context.Context, address string) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getDelegationsByValidator", []interface{}{address})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*rpc.Delegation
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetAllValidatorAddresses(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getAllValidatorAddresses", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return []string{}, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return []string{}, nil
}

func (h *HmyRpcClient) GetAllValidatorInformation(ctx context.Context, page int) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getAllValidatorInformation", []interface{}{page})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.ValidatorRPCEnhanced
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetAllValidatorInformationByBlockNumber(ctx context.Context, page int, number rpc.BlockNumber) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getAllValidatorInformationByBlockNumber", []interface{}{page, number})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.ValidatorRPCEnhanced
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetElectedValidatorAddresses(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getElectedValidatorAddresses", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return []string{}, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return []string{}, nil
}

func (h *HmyRpcClient) GetValidatorInformation(ctx context.Context, address string) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getValidatorInformation", []interface{}{address})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.ValidatorRPCEnhanced
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetCurrentUtilityMetrics(ctx context.Context) (*rpc.JSONRpcResp, error) {
	rpcResp, err := h.Request("getCurrentUtilityMetrics", []interface{}{})
	return rpcResp, err
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.UtilityMetric
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetMedianRawStakeSnapshot(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getMedianRawStakeSnapshot", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.CompletedEPoSRound
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetStakingNetworkInfo(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getStakingNetworkInfo", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *rpc.StakingNetworkInfo
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetSuperCommittees(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getSuperCommittees", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.Transition
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetCXReceiptByHash(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	return h.Request("getCXReceiptByHash", []interface{}{hash})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.CxReceipt
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetPendingCXReceipts(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getPendingCXReceipts", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.CXReceiptsProof
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) ResendCx(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	return h.Request("resendCx", []interface{}{hash})
	//if err != nil {
	//	return false, err
	//}
	//if rpcResp.Result != nil {
	//	var reply bool
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return false, nil
}

func (h *HmyRpcClient) GetPoolStats(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getPoolStats", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *rpc.PendingPoolStats
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) PendingStakingTransactions(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("pendingStakingTransactions", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.StakingTransaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) PendingTransactions(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("pendingTransactions", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.Transaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetCurrentStakingErrorSink(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getCurrentStakingErrorSink", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.TransactionErrorReport
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetStakingTransactionByBlockNumberAndIndex(ctx context.Context, number rpc.BlockNumber, index uint) (*rpc.JSONRpcResp, error) {
	return h.Request("getStakingTransactionByBlockNumberAndIndex", []interface{}{number, index})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.StakingTransaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetStakingTransactionByBlockHashAndIndex(ctx context.Context, hash string, index uint) (*rpc.JSONRpcResp, error) {
	return h.Request("getStakingTransactionByBlockHashAndIndex", []interface{}{hash, index})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.StakingTransaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetStakingTransactionByHash(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	return h.Request("getStakingTransactionByHash", []interface{}{hash})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.StakingTransaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) SendRawStakingTransaction(ctx context.Context, rawTx string) (*rpc.JSONRpcResp, error) {
	return h.Request("sendRawStakingTransaction", []interface{}{rawTx})
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}

func (h *HmyRpcClient) GetCurrentTransactionErrorSink(ctx context.Context) (*rpc.JSONRpcResp, error) {
	return h.Request("getCurrentTransactionErrorSink", []interface{}{})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply []*common.TransactionErrorReport
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetTransactionByBlockHashAndIndex(ctx context.Context, hash string, index uint) (*rpc.JSONRpcResp, error) {
	return h.Request("getTransactionByBlockHashAndIndex", []interface{}{hash, index})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.Transaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetTransactionByBlockNumberAndIndex(ctx context.Context, number rpc.BlockNumber, index uint) (*rpc.JSONRpcResp, error) {
	return h.Request("getTransactionByBlockNumberAndIndex", []interface{}{number, index})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.Transaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetTransactionByHash(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	return h.Request("getTransactionByHash", []interface{}{hash})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.Transaction
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) GetTransactionReceipt(ctx context.Context, hash string) (*rpc.JSONRpcResp, error) {
	return h.Request("getTransactionReceipt", []interface{}{hash})
	//if err != nil {
	//	return nil, err
	//}
	//if rpcResp.Result != nil {
	//	var reply *common.TxReceipt
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return nil, nil
}

func (h *HmyRpcClient) SendRawTransaction(ctx context.Context, rawTx string) (*rpc.JSONRpcResp, error) {
	return h.Request("sendRawTransaction", []interface{}{rawTx})
	//if err != nil {
	//	return "", err
	//}
	//if rpcResp.Result != nil {
	//	var reply string
	//	err = json.Unmarshal(*rpcResp.Result, &reply)
	//	return reply, err
	//}
	//return "", nil
}
