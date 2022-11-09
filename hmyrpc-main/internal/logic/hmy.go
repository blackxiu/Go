package logic

import (
	"context"
	"encoding/json"

	"github.com/shenzhendev/hmyrpc/common"
	"github.com/shenzhendev/hmyrpc/rpc"
	"github.com/shenzhendev/hmyrpc/rpc/hmy"
)

func (l *ApiLogic) setStakingTxHistoryCache(rpcResp *rpc.JSONRpcResp, args hmy.TxHistoryArgs) {
	// 初始化full tx list
	var stakingFullTxs []common.StakingTransaction
	// 依据是否 fulltx 写入redis
	if !args.FullTx {
		// 解析 hash
		var reply *rpc.StakingTxHistoryWithHash
		json.Unmarshal(*rpcResp.Result, &reply)
		hashList := reply.StakingTransactions
		// 获取fulltx
		stakingFullTxs = l.getStakingTxHistory(hashList)
	} else {
		var reply *rpc.StakingTxHistoryWithFullTx
		json.Unmarshal(*rpcResp.Result, &reply)
		stakingFullTxs = reply.StakingTransactions
	}
	// 写入cache
	l.svcCtx.Cache.SetStakingTransactionsHistory(args.Address, args.TxType, stakingFullTxs)

}

func (l *ApiLogic) getStakingTxHistory(hashList []string) []common.StakingTransaction {
	if len(hashList) == 0 {
		return nil
	}
	var stakingTxs []common.StakingTransaction
	for _, hash := range hashList {
		tx := l.getStakingTxByHash(hash)
		if tx != nil {
			stakingTxs = append(stakingTxs, *tx)
		}
	}
	return stakingTxs
}

func (l *ApiLogic) getStakingTxByHash(hash string) *common.StakingTransaction {
	resp, err := l.svcCtx.Endpoint.Get("getStakingTx").(*hmy.HmyRpcClient).GetStakingTransactionByHash(context.Background(), hash)
	if resp.Error != nil || err != nil {
		return nil
	}
	var tx *common.StakingTransaction
	json.Unmarshal(*resp.Result, &tx)
	return tx
}

func (l *ApiLogic) setBlocks(from, to rpc.BlockNumber) {
	// fixme: 解决此处暴力地获取数据
	// 获取 !fullTx
	var blocksWithTxHash []common.BlockWithTxHash
	resp, err := l.svcCtx.Endpoint.Get("getBlocks").(*hmy.HmyRpcClient).
		GetBlocks(context.Background(), from, to, &hmy.BlockArgs{
			WithSigners: true,
			FullTx:      false,
			InclStaking: true,
		})
	if resp.Error != nil || err != nil {
		return
	}
	json.Unmarshal(*resp.Result, &blocksWithTxHash)
	// 获取 fullTx
	var blocksWithFullTx []common.BlockWithFullTx
	resp, err = l.svcCtx.Endpoint.Get("getBlocks").(*hmy.HmyRpcClient).
		GetBlocks(context.Background(), from, to, &hmy.BlockArgs{
			WithSigners: false,
			FullTx:      true,
			InclStaking: true,
		})
	if resp.Error != nil || err != nil {
		return
	}
	json.Unmarshal(*resp.Result, &blocksWithFullTx)
	// set cache
	l.svcCtx.Cache.SetBlocks(blocksWithFullTx, blocksWithTxHash)
}

func (l *ApiLogic) setBlock(rpcResp *rpc.JSONRpcResp) {
	var header *common.Header
	json.Unmarshal(*rpcResp.Result, &header)

	//if block != nil {
	number := rpc.BlockNumber(header.Number.Int64())
	l.setBlocks(number, number)
	//}
}

func (l *ApiLogic) setSigners(rpcResp *rpc.JSONRpcResp, height rpc.BlockNumber) {
	var signers []string
	json.Unmarshal(*rpcResp.Result, &signers)
	l.svcCtx.Cache.SetBlockSigners(height, signers)
}

func (l *ApiLogic) setSignerKeys(rpcResp *rpc.JSONRpcResp, height rpc.BlockNumber) {
	var signerKeys []string
	json.Unmarshal(*rpcResp.Result, &signerKeys)
	l.svcCtx.Cache.SetBlockSignersKeys(height, signerKeys)
}

func (l *ApiLogic) setBlockTxCountByNumber(rpcResp *rpc.JSONRpcResp, height rpc.BlockNumber) {
	l.setBlocks(height, height)
}

func (l *ApiLogic) setHeaderInfo(rpcResp *rpc.JSONRpcResp) {
	var header *rpc.HeaderInformation
	json.Unmarshal(*rpcResp.Result, &header)

	//if block != nil {
	number := rpc.BlockNumber(header.BlockNumber)
	go l.setBlocks(number, number)
	l.svcCtx.Cache.SetHeaderInfo(number, header)
	//}
}

func (l *ApiLogic) setLatestChainHeader(rpcResp *rpc.JSONRpcResp) {
	var latestHeader *rpc.HeaderPair
	json.Unmarshal(*rpcResp.Result, &latestHeader)
	l.svcCtx.Cache.SetLatestChainHeaders(latestHeader)
}

func (l *ApiLogic) setLatestHeader(rpcResp *rpc.JSONRpcResp) {
	var header *rpc.HeaderInformation
	json.Unmarshal(*rpcResp.Result, &header)
	l.svcCtx.Cache.SetLatestHeader(header)
}

func (l *ApiLogic) setBlockNumber(rpcResp *rpc.JSONRpcResp) {
	var height *rpc.BlockNumber
	json.Unmarshal(*rpcResp.Result, &height)
	l.svcCtx.Cache.SetBlockNumber(height)
}

func (l *ApiLogic) SetValidators(rpcResp *rpc.JSONRpcResp, epoch uint64) {
	var validators *rpc.ValidatorsFields
	json.Unmarshal(*rpcResp.Result, &validators)
	l.svcCtx.Cache.SetValidators(epoch, -1, validators)
}

func (l *ApiLogic) SetValidatorKeys(rpcResp *rpc.JSONRpcResp, epoch uint64) {
	var validatorKeys []string
	json.Unmarshal(*rpcResp.Result, &validatorKeys)
	l.svcCtx.Cache.SetValidatorKeys(epoch, -1, validatorKeys)
}

func (l *ApiLogic) setMetaData(rpcResp *rpc.JSONRpcResp) {
	if rpcResp.Result != nil {
		var metadata *common.NodeMetadata
		json.Unmarshal(*rpcResp.Result, &metadata)
		l.svcCtx.Cache.SetMetadata(metadata)
	}
}

func (l *ApiLogic) setCode(rpcResp *rpc.JSONRpcResp, address string) {
	if rpcResp.Result == nil {
		return
	}
	var code string
	json.Unmarshal(*rpcResp.Result, &code)
	l.svcCtx.Cache.SetCode(address, code)
}

func (l *ApiLogic) setStakingTx(rpcResp *rpc.JSONRpcResp, hash string) {
	if rpcResp.Result == nil {
		return
	}
	var tx *common.StakingTransaction
	json.Unmarshal(*rpcResp.Result, &tx)
	l.svcCtx.Cache.SetStakingTransaction(tx)
}

func (l *ApiLogic) setTxByBlkHashAndIndex(rpcResp *rpc.JSONRpcResp) {
	if rpcResp.Result == nil {
		return
	}
	var tx *common.Transaction
	json.Unmarshal(*rpcResp.Result, &tx)
	l.svcCtx.Cache.SetTransactionByBlockHashAndIndex(tx)
}

func (l *ApiLogic) setTxByBlkAndIndex(rpcResp *rpc.JSONRpcResp) {
	if rpcResp.Result == nil {
		return
	}
	var tx *common.Transaction
	json.Unmarshal(*rpcResp.Result, &tx)
	l.svcCtx.Cache.SetTransactionByBlockNumberAndIndex(tx)
}

func (l *ApiLogic) setTx(rpcResp *rpc.JSONRpcResp, hash string) {
	if rpcResp.Result == nil {
		return
	}
	var tx *common.Transaction
	json.Unmarshal(*rpcResp.Result, &tx)
	l.svcCtx.Cache.SetTransaction(tx)
}

func (l *ApiLogic) setTxReceipt(rpcResp *rpc.JSONRpcResp, hash string) {
	if rpcResp.Result == nil {
		return
	}
	var receipt *common.TxReceipt
	json.Unmarshal(*rpcResp.Result, &receipt)
	l.svcCtx.Cache.SetTransactionReceipt(hash, receipt)
}
