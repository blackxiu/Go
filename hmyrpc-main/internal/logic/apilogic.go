package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/INFURA/go-ethlibs/jsonrpc"
	"github.com/shenzhendev/hmyrpc/balancing"
	"github.com/shenzhendev/hmyrpc/internal/svc"
	"github.com/shenzhendev/hmyrpc/internal/types"
	"github.com/shenzhendev/hmyrpc/rpc"
	"github.com/shenzhendev/hmyrpc/rpc/hmy"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) ApiLogic {
	return ApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiLogic) HandleApiRequest(req types.Request) (resp *types.Response, err error) {
	// todo: get data from cache or chain rpc server
	// get data
	id := req.ID
	methods := strings.Split(req.Method, "_")
	params := req.Params
	// parse method
	if len(methods) != 2 {
		println(">>>>>>>!!!>>>>", req.Method)
		return l.InvalidResp("2.0", id), nil
	}
	prefix, method := methods[0], methods[1]
	_ = prefix

	// return fake data to http.ResponseWriter
	return l.handleApiMessage(id, method, params, l.svcCtx.Endpoint)
}

func (l *ApiLogic) handleApiMessage(id uint64, method string, params jsonrpc.Params, balance balancing.LoadBalance) (resp *types.Response, err error) {

	// todo: 降低圈复杂度 ！
	// req data
	switch method {
	case "getBalance":
		//var balance *big.Int
		var address string
		params.UnmarshalSingleParam(0, &address)
		if l.useCache() {
			balance, err := l.svcCtx.Cache.GetBalance(address)
			//fmt.Println(">>>>>>>>>>>>>>>> balance ", balance, err)
			if balance != nil {
				raw, _ := balance.MarshalJSON()
				return l.SuccessResponse(id, raw), err
			}
		}

		res, err := balance.Get("getBalance").(*hmy.HmyRpcClient).GetBalance(context.Background(), address)
		if err != nil {
			var req map[string]interface{}
			req["code"] = -32600
			req["message"] = err.Error()
			return l.ErrorResponse(id, req), err
		}
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		//fmt.Println("_++_+_+)_+)+)_+)+)", res.Result, res.Error)
		var r *big.Int
		if res.Result != nil {
			err = json.Unmarshal(*res.Result, &r)
		}
		// set to redis
		l.svcCtx.Cache.SetBalance(address, r)
		//log.Fatalf(">>>>>", err)

		return l.SuccessResponse(id, *res.Result), nil

	case "getBalanceByBlockNumber":
		var address string
		var num int64
		var test string
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &num)
		params.UnmarshalSingleParam(1, &test)
		// fixme: shit hmy ⚠️
		// 出现此代码的原因是：hmy混用了int和string当参数
		if num == 0 && test != "" {
			if test == "latest" {
				num = -1
			} else {
				res, _ := strconv.Atoi(test)
				num = int64(res)
			}
		}

		if l.useCache() {
			balance, err := l.svcCtx.Cache.GetBalanceByBlockNumber(address, num)
			fmt.Println(">>>>>>>>>>>>>>>> balance ", balance, err)
			if balance != nil {
				raw, _ := balance.MarshalJSON()
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getBalanceByBlockNumber").(*hmy.HmyRpcClient).GetBalanceByBlockNumber(context.Background(), address, num)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// set cache
		var r *big.Int
		if res.Result != nil {
			err = json.Unmarshal(*res.Result, &r)
		}
		err = l.svcCtx.Cache.SetBalanceByNumber(address, num, r)

		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingTransactionsCount":
		var address string
		var txType string
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &txType)

		if l.useCache() {
			count, err := l.svcCtx.Cache.GetStakingTransactionsCount(address, txType)
			fmt.Println(">>>>>>>>>>>>>>>> balance ", count, err)
			if count != nil {
				raw, _ := count.MarshalJSON()
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getStakingTransactionsCount").(*hmy.HmyRpcClient).GetStakingTransactionsCount(context.Background(), hmy.TxsCountArgs{Address: address, Type: txType})
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// set cache
		var r *big.Int
		if res.Result != nil {
			err = json.Unmarshal(*res.Result, &r)
		}
		err = l.svcCtx.Cache.SetStakingTransactionsCount(address, txType, r)

		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingTransactionsHistory":
		var args hmy.TxHistoryArgs
		params.UnmarshalSingleParam(0, &args)

		//if l.useCache() {
		//	stakingTxsHistory, err := l.svcCtx.Cache.GetStakingTransactionsHistory(args)
		//	fmt.Println(">>>>>>>>>>>>>>>> stakingTxsHistory ", stakingTxsHistory, err)
		//	if stakingTxsHistory != nil {
		//		fmt.Println(">>>>>>>>>>>>>>>> stakingTxsHistory not nil")
		//		raw, _ := json.Marshal(stakingTxsHistory)
		//		return l.SuccessResponse(id, raw), err
		//	}
		//}

		res, _ := balance.Get("getStakingTransactionsHistory").(*hmy.HmyRpcClient).GetStakingTransactionsHistory(context.Background(), args)
		//fmt.Println(res.Result, res.Error, err)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// todo： 此处的一个bug，https://hmyapis0.metamemo.one返回结果跟官方api结果不一致
		// 如果后续一致，可以去掉此处代码
		if res.Result == nil {
			return l.SuccessResponse(id, json.RawMessage(`{"staking_transactions": []}`)), nil
		}
		//
		if args.FullTx {
			var r rpc.StakingTxHistoryWithFullTx
			json.Unmarshal(*res.Result, &r)
			err = l.svcCtx.Cache.SetStakingTransactionsHistory(args.Address, args.TxType, r.StakingTransactions)
		}
		//go l.setStakingTxHistoryCache(res, args)

		return l.SuccessResponse(id, *res.Result), nil
	case "getTransactionsCount":
		var address string
		var txType string
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &txType)

		if l.useCache() {
			count, err := l.svcCtx.Cache.GetTransactionsCount(address, txType)
			fmt.Println(">>>>>>>>>>>>>>>> balance ", count, err)
			if count != nil {
				raw, _ := count.MarshalJSON()
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getTransactionsCount").(*hmy.HmyRpcClient).GetTransactionsCount(context.Background(), hmy.TxsCountArgs{Address: address, Type: txType})
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}

		// set cache
		var r *big.Int
		if res.Result != nil {
			err = json.Unmarshal(*res.Result, &r)
		}
		err = l.svcCtx.Cache.SetTransactionsCount(address, txType, r)

		return l.SuccessResponse(id, *res.Result), nil

	case "getTransactionsHistory":
		var args hmy.TxHistoryArgs
		params.UnmarshalSingleParam(0, &args)
		res, _ := balance.Get("getTransactionsHistory").(*hmy.HmyRpcClient).GetTransactionsHistory(context.Background(), args)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// 如果后续一致，可以去掉此处代码
		if res.Result == nil {
			return l.SuccessResponse(id, json.RawMessage(`{"transactions": []}`)), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlocks":
		var from, to rpc.BlockNumber
		var args hmy.BlockArgs
		params.UnmarshalSingleParam(0, &from)
		params.UnmarshalSingleParam(1, &to)
		params.UnmarshalSingleParam(2, &args)
		res, _ := balance.Get("getBlocks").(*hmy.HmyRpcClient).GetBlocks(context.Background(), from, to, &args)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// write cache
		go l.setBlocks(from, to)

		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockByNumber":
		var number rpc.BlockNumber
		var args hmy.BlockArgs
		params.UnmarshalSingleParam(0, &number)
		params.UnmarshalSingleParam(1, &args)

		if l.useCache() {
			blk, err := l.svcCtx.Cache.GetBlockByNumber(number, args)
			if blk != nil {
				raw, _ := json.Marshal(blk)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getBlockByNumber").(*hmy.HmyRpcClient).GetBlockByNumber(context.Background(), number, &args)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		if res.Result == nil {
			return l.SuccessResponse(id, json.RawMessage(``)), nil
		}
		// write cache
		go l.setBlocks(number, number)
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockByHash":
		var hash string
		var args hmy.BlockArgs
		params.UnmarshalSingleParam(0, &hash)
		params.UnmarshalSingleParam(1, &args)

		if l.useCache() {
			blk, err := l.svcCtx.Cache.GetBlockByHash(hash, args)
			//fmt.Println(">>>>>>>>>>>>>>>> blk ", blk, err)
			if blk != nil {
				raw, _ := json.Marshal(blk)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getBlockByHash").(*hmy.HmyRpcClient).GetBlockByHash(context.Background(), hash, &args)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		if res.Result == nil {
			return l.SuccessResponse(id, json.RawMessage(`null`)), nil
		}
		// write block cache
		go l.setBlock(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockSigners":
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &number)

		if l.useCache() {
			signers, err := l.svcCtx.Cache.GetBlockSigners(number)
			if signers != nil {
				raw, _ := json.Marshal(signers)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getBlockSigners").(*hmy.HmyRpcClient).GetBlockSigners(context.Background(), number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// write to cache
		go l.setSigners(res, number)
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockSignerKeys":
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &number)
		if l.useCache() {
			signerKeys, err := l.svcCtx.Cache.GetBlockSignersKeys(number)
			if signerKeys != nil {
				raw, _ := json.Marshal(signerKeys)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getBlockSignerKeys").(*hmy.HmyRpcClient).GetBlockSignersKeys(context.Background(), number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// write to cache
		go l.setSignerKeys(res, number)
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockTransactionCountByNumber":
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &number)
		if l.useCache() {
			count, err := l.svcCtx.Cache.GetBlockTxCountByNumber(number)
			// -1 == not found
			if count != -1 {
				raw, _ := json.Marshal(count)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getBlockTransactionCountByNumber").(*hmy.HmyRpcClient).GetBlockTransactionCountByNumber(context.Background(), number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setBlockTxCountByNumber(res, number)
		return l.SuccessResponse(id, *res.Result), nil

	case "getBlockTransactionCountByHash":
		var hash string
		params.UnmarshalSingleParam(0, &hash)
		if l.useCache() {
			count, err := l.svcCtx.Cache.GetBlockTxCountByHash(hash)
			// -1 == not found
			if count != -1 {
				raw, _ := json.Marshal(count)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getBlockTransactionCountByHash").(*hmy.HmyRpcClient).GetBlockTransactionCountByHash(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getHeaderByNumber":
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &number)
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetHeaderByNumber(number)
			if header != nil {
				raw, _ := json.Marshal(*header)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getHeaderByNumber").(*hmy.HmyRpcClient).GetHeaderByNumber(context.Background(), number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setHeaderInfo(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "getLatestChainHeaders":
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetLatestChainHeaders()
			if header != nil {
				raw, _ := json.Marshal(*header)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("getLatestChainHeaders").(*hmy.HmyRpcClient).GetLatestChainHeaders(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setLatestChainHeader(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "latestHeader":
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetLatestHeader()
			if header != nil {
				raw, _ := json.Marshal(*header)
				return l.SuccessResponse(id, raw), err
			}
		}

		res, _ := balance.Get("latestHeader").(*hmy.HmyRpcClient).LatestHeader(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setLatestHeader(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "blockNumber":
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetBlockNumber()
			if header != nil {
				raw, _ := json.Marshal(*header)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("blockNumber").(*hmy.HmyRpcClient).BlockNumber(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setBlockNumber(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "getCirculatingSupply":
		res, _ := balance.Get("getCirculatingSupply").(*hmy.HmyRpcClient).GetCirculatingSupply(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getEpoch":
		res, _ := balance.Get("getEpoch").(*hmy.HmyRpcClient).GetEpoch(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "epochLastBlock":
		var epoch uint64
		params.UnmarshalSingleParam(0, &epoch)
		res, _ := balance.Get("epochLastBlock").(*hmy.HmyRpcClient).GetEpochLastBlock(context.Background(), epoch)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getLastCrossLinks":
		res, _ := balance.Get("getLastCrossLinks").(*hmy.HmyRpcClient).GetLastCrossLinks(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getLeader":
		res, _ := balance.Get("getLeader").(*hmy.HmyRpcClient).GetLeader(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "gasPrice":
		res, _ := balance.Get("gasPrice").(*hmy.HmyRpcClient).GasPrice(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getShardingStructure":
		res, _ := balance.Get("getShardingStructure").(*hmy.HmyRpcClient).GetShardingStructure(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getTotalSupply":
		res, _ := balance.Get("getTotalSupply").(*hmy.HmyRpcClient).GetTotalSupply(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getValidators":
		var epoch uint64
		params.UnmarshalSingleParam(0, &epoch)
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetValidators(epoch)
			if header != nil {
				raw, _ := json.Marshal(*header)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getValidators").(*hmy.HmyRpcClient).GetValidators(context.Background(), epoch)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.SetValidators(res, epoch)
		return l.SuccessResponse(id, *res.Result), nil

	case "getValidatorKeys":
		var epoch uint64
		params.UnmarshalSingleParam(0, &epoch)
		if l.useCache() {
			header, err := l.svcCtx.Cache.GetValidatorKeys(epoch)
			if header != nil {
				raw, _ := json.Marshal(header)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getValidatorKeys").(*hmy.HmyRpcClient).GetValidatorKeys(context.Background(), epoch)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.SetValidatorKeys(res, epoch)
		return l.SuccessResponse(id, *res.Result), nil
	case "getCurrentBadBlocks":
		return l.SuccessResponse(id, json.RawMessage(`[]`)), nil

	case "getNodeMetadata":
		if l.useCache() {
			metadata, err := l.svcCtx.Cache.GetMetadata()
			if metadata != nil {
				raw, _ := json.Marshal(metadata)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getNodeMetadata").(*hmy.HmyRpcClient).GetNodeMetadata(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setMetaData(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "protocolVersion":
		res, _ := balance.Get("protocolVersion").(*hmy.HmyRpcClient).ProtocolVersion(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "peerCount":
		res, _ := balance.Get("peerCount").(*hmy.HmyRpcClient).PeerCount(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "call":
		var args rpc.CallArgs
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &args)
		params.UnmarshalSingleParam(1, &number)
		res, _ := balance.Get("call").(*hmy.HmyRpcClient).Call(context.Background(), args, number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "estimateGas":
		var args rpc.CallArgs
		params.UnmarshalSingleParam(0, &args)
		res, _ := balance.Get("estimateGas").(*hmy.HmyRpcClient).EstimateGas(context.Background(), args)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getCode":
		var address string
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &number)
		if l.useCache() {
			metadata, err := l.svcCtx.Cache.GetMetadata()
			if metadata != nil {
				raw, _ := json.Marshal(metadata)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getCode").(*hmy.HmyRpcClient).GetCode(context.Background(), address, number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getStorageAt":
		var address, key string
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &key)
		params.UnmarshalSingleParam(2, &number)
		res, _ := balance.Get("getStorageAt").(*hmy.HmyRpcClient).GetStorageAt(context.Background(), address, key, number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getDelegationsByDelegator":
		var address string
		params.UnmarshalSingleParam(0, &address)
		res, _ := balance.Get("getDelegationsByDelegator").(*hmy.HmyRpcClient).GetDelegationsByDelegator(context.Background(), address)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getDelegationsByDelegatorByBlockNumber":
		var address string
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &address)
		params.UnmarshalSingleParam(1, &number)
		res, _ := balance.Get("getDelegationsByDelegatorByBlockNumber").(*hmy.HmyRpcClient).GetDelegationsByDelegatorByBlockNumber(context.Background(), address, number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		// todo： 此处的一个bug，https://hmyapis0.metamemo.one返回结果跟官方api结果不一致
		// 如果后续一致，可以去掉此处代码
		if res.Result == nil {
			return l.SuccessResponse(id, json.RawMessage(nil)), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getDelegationsByValidator":
		var address string
		params.UnmarshalSingleParam(0, &address)
		res, _ := balance.Get("getDelegationsByValidator").(*hmy.HmyRpcClient).GetDelegationsByValidator(context.Background(), address)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getAllValidatorAddresses":
		res, _ := balance.Get("getAllValidatorAddresses").(*hmy.HmyRpcClient).GetAllValidatorAddresses(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getAllValidatorInformation":
		var page int
		params.UnmarshalSingleParam(0, &page)
		res, _ := balance.Get("getAllValidatorInformation").(*hmy.HmyRpcClient).GetAllValidatorInformation(context.Background(), page)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getAllValidatorInformationByBlockNumber":
		var page int
		var number rpc.BlockNumber
		params.UnmarshalSingleParam(0, &page)
		params.UnmarshalSingleParam(1, &number)
		res, _ := balance.Get("getAllValidatorInformationByBlockNumber").(*hmy.HmyRpcClient).GetAllValidatorInformationByBlockNumber(context.Background(), page, number)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}

		return l.SuccessResponse(id, *res.Result), nil

	case "getElectedValidatorAddresses":
		res, _ := balance.Get("getElectedValidatorAddresses").(*hmy.HmyRpcClient).GetElectedValidatorAddresses(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getValidatorInformation":
		var address string
		params.UnmarshalSingleParam(0, &address)
		res, _ := balance.Get("getValidatorInformation").(*hmy.HmyRpcClient).GetValidatorInformation(context.Background(), address)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getCurrentUtilityMetrics":
		res, _ := balance.Get("getCurrentUtilityMetrics").(*hmy.HmyRpcClient).GetCurrentUtilityMetrics(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getMedianRawStakeSnapshot":
		res, _ := balance.Get("getMedianRawStakeSnapshot").(*hmy.HmyRpcClient).GetMedianRawStakeSnapshot(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingNetworkInfo":
		res, _ := balance.Get("getStakingNetworkInfo").(*hmy.HmyRpcClient).GetStakingNetworkInfo(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getSuperCommittees":
		res, _ := balance.Get("getSuperCommittees").(*hmy.HmyRpcClient).GetSuperCommittees(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getCXReceiptByHash":
		var hash string
		res, _ := balance.Get("getCXReceiptByHash").(*hmy.HmyRpcClient).GetCXReceiptByHash(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getPendingCXReceipts":
		res, _ := balance.Get("getPendingCXReceipts").(*hmy.HmyRpcClient).GetPendingCXReceipts(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "resendCx":
		var hash string
		res, _ := balance.Get("resendCx").(*hmy.HmyRpcClient).ResendCx(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getPoolStats":
		res, _ := balance.Get("getPoolStats").(*hmy.HmyRpcClient).GetPoolStats(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "pendingStakingTransactions":
		res, _ := balance.Get("pendingStakingTransactions").(*hmy.HmyRpcClient).PendingStakingTransactions(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "pendingTransactions":
		res, _ := balance.Get("pendingTransactions").(*hmy.HmyRpcClient).PendingTransactions(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getCurrentStakingErrorSink":
		res, _ := balance.Get("getCurrentStakingErrorSink").(*hmy.HmyRpcClient).GetCurrentStakingErrorSink(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingTransactionByBlockNumberAndIndex":
		var number rpc.BlockNumber
		var index uint
		params.UnmarshalSingleParam(0, &number)
		params.UnmarshalSingleParam(1, &index)
		res, _ := balance.Get("getStakingTransactionByBlockNumberAndIndex").(*hmy.HmyRpcClient).
			GetStakingTransactionByBlockNumberAndIndex(context.Background(), number, index)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingTransactionByBlockHashAndIndex":
		var hash string
		var index uint
		params.UnmarshalSingleParam(0, &hash)
		params.UnmarshalSingleParam(1, &index)
		res, _ := balance.Get("getStakingTransactionByBlockNumberAndIndex").(*hmy.HmyRpcClient).
			GetStakingTransactionByBlockHashAndIndex(context.Background(), hash, index)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getStakingTransactionByHash":
		var hash string
		params.UnmarshalSingleParam(0, &hash)
		if l.useCache() {
			tx, err := l.svcCtx.Cache.GetStakingTransactionByHash(hash)
			if tx != nil {
				raw, _ := json.Marshal(tx)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getStakingTransactionByHash").(*hmy.HmyRpcClient).
			GetStakingTransactionByHash(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setStakingTx(res, hash)
		return l.SuccessResponse(id, *res.Result), nil

	case "sendRawStakingTransaction":
		var rawTx string
		params.UnmarshalSingleParam(0, &rawTx)
		res, _ := balance.Get("sendRawStakingTransaction").(*hmy.HmyRpcClient).
			SendRawStakingTransaction(context.Background(), rawTx)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getCurrentTransactionErrorSink":
		res, _ := balance.Get("getCurrentTransactionErrorSink").(*hmy.HmyRpcClient).GetCurrentTransactionErrorSink(context.Background())
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	case "getTransactionByBlockHashAndIndex":
		var hash string
		var index uint
		params.UnmarshalSingleParam(0, &hash)
		params.UnmarshalSingleParam(1, &index)
		if l.useCache() {
			tx, err := l.svcCtx.Cache.GetTransactionByBlockHashAndIndex(hash, int64(index))
			if tx != nil {
				raw, _ := json.Marshal(tx)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getTransactionByBlockHashAndIndex").(*hmy.HmyRpcClient).
			GetTransactionByBlockHashAndIndex(context.Background(), hash, index)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setTxByBlkHashAndIndex(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "getTransactionByBlockNumberAndIndex":
		var number rpc.BlockNumber
		var index uint
		params.UnmarshalSingleParam(0, &number)
		params.UnmarshalSingleParam(1, &index)
		if l.useCache() {
			tx, err := l.svcCtx.Cache.GetTransactionByBlockNumberAndIndex(number, int64(index))
			if tx != nil {
				raw, _ := json.Marshal(tx)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getTransactionByBlockNumberAndIndex").(*hmy.HmyRpcClient).
			GetTransactionByBlockNumberAndIndex(context.Background(), number, index)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setTxByBlkAndIndex(res)
		return l.SuccessResponse(id, *res.Result), nil

	case "getTransactionByHash":
		var hash string
		params.UnmarshalSingleParam(0, &hash)
		if l.useCache() {
			tx, err := l.svcCtx.Cache.GetTransactionByHash(hash)
			if tx != nil {
				raw, _ := json.Marshal(tx)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getTransactionByHash").(*hmy.HmyRpcClient).GetTransactionByHash(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setTx(res, hash)
		return l.SuccessResponse(id, *res.Result), nil

	case "getTransactionReceipt":
		var hash string
		params.UnmarshalSingleParam(0, &hash)
		if l.useCache() {
			receipt, err := l.svcCtx.Cache.GetTransactionReceipt(hash)
			if receipt != nil {
				raw, _ := json.Marshal(receipt)
				return l.SuccessResponse(id, raw), err
			}
		}
		res, _ := balance.Get("getTransactionReceipt").(*hmy.HmyRpcClient).GetTransactionReceipt(context.Background(), hash)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		go l.setTxReceipt(res, hash)
		return l.SuccessResponse(id, *res.Result), nil

	case "sendRawTransaction":
		var rawTx string
		params.UnmarshalSingleParam(0, &rawTx)
		res, _ := balance.Get("sendRawTransaction").(*hmy.HmyRpcClient).
			SendRawTransaction(context.Background(), rawTx)
		if res.Error != nil {
			return l.ErrorResponse(id, res.Error), nil
		}
		return l.SuccessResponse(id, *res.Result), nil

	default:
		return l.InvalidResp("2.0", id), nil
	}
	return nil, err
}

//
func (l *ApiLogic) InvalidResp(jsonrpc string, id uint64) *types.Response {
	return &types.Response{
		Jsonrpc: jsonrpc,
		Id:      id,
		Error: &types.JsonError{
			Code:    -32600,
			Message: "invalid request",
			Data:    nil,
		},
		Result: nil,
	}
}

func (l *ApiLogic) SuccessResponse(id uint64, message json.RawMessage) *types.Response {
	return &types.Response{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  message,
	}
}

func (l *ApiLogic) ErrorResponse(id uint64, err map[string]interface{}) *types.Response {
	return &types.Response{
		Jsonrpc: "2.0",
		Id:      id,
		Error: &types.JsonError{
			Code:    int(err["code"].(float64)),
			Message: err["message"].(string),
			Data:    nil,
		},
	}
}

func (l *ApiLogic) useCache() bool {
	// todo 补充使用cache代码
	return true
}

func (l *ApiLogic) getFromCache(req *types.Request) (resp interface{}, err error) {

	return nil, err
}

func (l *ApiLogic) storeToCache(req *types.Request, resp *types.Response) (err error) {

	return nil
}

func (l *ApiLogic) getFromRPC(req *types.Request, url string) (resp interface{}, err error) {

	return nil, err
}
