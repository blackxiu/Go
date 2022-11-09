package cache

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shenzhendev/hmyrpc/common"
	"github.com/shenzhendev/hmyrpc/rpc"
	"github.com/shenzhendev/hmyrpc/rpc/hmy"
	"golang.org/x/exp/errors/fmt"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Password string `json:"password"`
	Database int    `json:"database"`
	PoolSize int    `json:"poolSize"`
}

type RedisClient struct {
	client  *redis.Client
	prefix  string
	timeout time.Duration
}

func NewRedisClient(cfg *Config, prefix string, timeout time.Duration) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})
	return &RedisClient{client: client, prefix: prefix, timeout: timeout}
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Check() (string, error) {
	return r.client.Ping(context.Background()).Result()
}

func (r *RedisClient) BgSave() (string, error) {
	return r.client.BgSave(context.Background()).Result()
}

func (r *RedisClient) formatKey(args ...interface{}) string {
	return join(r.prefix, join(args...))
}

func (r *RedisClient) formatRound(height int64, nonce string) string {
	return r.formatKey("shares", "round"+strconv.FormatInt(height, 10), nonce)
}

func join(args ...interface{}) string {
	s := make([]string, len(args))
	for i, v := range args {
		switch v.(type) {
		case string:
			s[i] = v.(string)
		case int64:
			s[i] = strconv.FormatInt(v.(int64), 10)
		case uint64:
			s[i] = strconv.FormatUint(v.(uint64), 10)
		case float64:
			s[i] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		case bool:
			if v.(bool) {
				s[i] = "1"
			} else {
				s[i] = "0"
			}
		case *big.Int:
			n := v.(*big.Int)
			if n != nil {
				s[i] = n.String()
			} else {
				s[i] = "0"
			}
		default:
			panic("Invalid type specified for conversion")
		}
	}
	return strings.Join(s, ":")
}

func (r *RedisClient) SetBalance(address string, balance *big.Int) error {
	//tx := r.client.Multi()
	//defer tx.Close()
	//_, err := tx.Exec(func() error {
	//	tx.HMSet(r.formatKey("balance"), address, balance.String())
	//	tx.Expire(r.formatKey("balance"), r.timeout)
	//	return nil
	//})

	pipe := r.client.TxPipeline()
	defer pipe.Close()
	pipe.HMSet(context.Background(), r.formatKey("balanceLatest", address), map[string]interface{}{
		"latest": balance.String(),
	})
	// todo: 重新考虑过期时间
	pipe.Expire(context.Background(), r.formatKey("balanceLatest", address), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds

	return err
}

func (r *RedisClient) GetBalance(address string) (*big.Int, error) {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMGet(r.formatKey("balance"), address)
	//if cmd.Err() != nil {
	//	return nil, cmd.Err()
	//}

	//
	cmd := r.client.HMGet(context.Background(), r.formatKey("balanceLatest", address), "latest")
	if cmd.Val()[0] == nil {
		return nil, nil
	}
	result := new(big.Int)
	_, ok := result.SetString(cmd.Val()[0].(string), 10)
	if !ok {
		return nil, fmt.Errorf("[ %s ] : [ %s ]", "redis", "set big.Int failed")
	}
	return result, nil
}

// SetBalanceByNumber todo：此key暂时不设计过期时间，如有必要，可以考虑设置过期时间
// 这里写入前要判断当前传入的 number 是否高于现有
func (r *RedisClient) SetBalanceByNumber(address string, number int64, balance *big.Int) error {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMSet(r.formatKey("balance", join(number)), address, balance.String())
	//return cmd.Err()
	if number == int64(rpc.LatestBlockNumber) {
		return r.SetBalance(address, balance)
	}

	pipe := r.client.TxPipeline()
	defer pipe.Close()
	//
	pipe.HMSet(context.Background(), r.formatKey("balance", address), map[string]interface{}{
		join(number): balance.String(),
	})
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetBalanceByBlockNumber(address string, number int64) (*big.Int, error) {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMGet(r.formatKey("balance", join(number)), address)
	//if cmd.Err() != nil {
	//	return nil, cmd.Err()
	//}
	// latest number
	if number == int64(rpc.LatestBlockNumber) {
		return r.GetBalance(address)
	}

	pipe := r.client.TxPipeline()
	defer pipe.Close()

	//
	cmd := pipe.HMGet(context.Background(), r.formatKey("balance", address), join(number))
	pipe.Exec(context.Background())

	if cmd.Val()[0] == nil {
		return nil, nil
	}
	result := new(big.Int)
	_, ok := result.SetString(cmd.Val()[0].(string), 10)
	if !ok {
		return nil, fmt.Errorf("[ %s ] : [ %s ]", "redis", "set big.Int failed")
	}
	return result, nil
}

func (r *RedisClient) SetStakingTransactionsCount(address string, txType string, count *big.Int) error {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMSet(r.formatKey("stakingTxCount", txType), address, strconv.FormatInt(int64(count), 10))
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	pipe.HMSet(context.Background(), r.formatKey("stakingTxCount", address), map[string]interface{}{
		txType: count.String(),
	})
	pipe.Expire(context.Background(), r.formatKey("stakingTxCount", address), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetStakingTransactionsCount(address, txType string) (*big.Int, error) {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMGet(r.formatKey("stakingTxCount", txType), address)
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	cmd := pipe.HMGet(context.Background(), r.formatKey("stakingTxCount", address), txType)
	pipe.Exec(context.Background())
	fmt.Println("REDIS >>>>>>>> ", cmd.Val())
	// len
	if cmd.Val()[0] == nil {
		return nil, nil
	}
	//
	result := new(big.Int)
	_, ok := result.SetString(cmd.Val()[0].(string), 10)
	if !ok {
		return nil, fmt.Errorf("[ %s ] : [ %s ]", "redis", "set big.Int failed")
	}
	return result, nil
}

// SetStakingTransactionsHistory
// Deprecated. fixme: 此接口不好处理，有好办法了再处理
func (r *RedisClient) SetStakingTransactionsHistory(address string, txType string, txs []common.StakingTransaction) error {
	// init
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	// get count
	count, err := r.client.ZCard(context.Background(),
		r.formatKey("stakingTxHistoryWithHash", address, txType)).Result()
	fmt.Println(">>>>>>>>>>>>>>>> SetStakingTransactionsHistory ", count, err)
	// if error
	if err != nil {
		return err
	}
	//
	for index, tx := range txs {
		// check tx exist
		//cmd := pipe.HMGet(context.Background(), r.formatKey("stakingTxHistoryWithHash", address), tx.Hash)
		val, err := r.client.ZAdd(context.Background(), r.formatKey("stakingTxHistoryWithHash", address, txType), &redis.Z{
			Score:  float64(count + int64(index)),
			Member: tx.Hash,
		}).Result()
		fmt.Println(">>>>>>>>>>>>>>>> SetStakingTransactionsHistory: stakingTxHistoryWithHash ", val, err)
		// is exist or write tx error
		if val == 0 || err != nil {
			continue
		}
		// write full tx
		txBytes, _ := json.Marshal(tx)
		cmd := pipe.HMSet(context.Background(), r.formatKey("stakingTxs"), map[string]interface{}{
			tx.Hash: txBytes,
		})
		fmt.Println(">>>>>>>>>>>>>>>> SetStakingTransactionsHistory: stakingTx ", cmd.Err(), cmd.Val())

		//if val == 0 {
		//	//pipe.HMSet(context.Background(), r.formatKey("stakingTxHistoryWithHash", address), map[uint64]interface{}{
		//	//	// just write tx hash
		//	//	tx.Nonce: tx.Hash,
		//	//})
		//	pipe.ZAdd(context.Background(), r.formatKey("stakingTxHistoryWithHash", address), redis.Z{
		//		Score:  float64(tx.Nonce),
		//		Member: tx.Hash,
		//	})
		//
		//}
	}

	_, err = pipe.Exec(context.Background())
	return err
}

// GetStakingTransactionsHistory
// Deprecated. fixme: 此接口不好处理，有好办法了再处理
func (r *RedisClient) GetStakingTransactionsHistory(args hmy.TxHistoryArgs) (map[string]interface{}, error) {
	//pipe := r.client.TxPipeline()
	//defer pipe.Close()

	start := int64(args.PageIndex * args.PageSize)
	stop := int64((args.PageIndex + 1) * args.PageSize)

	count, err := r.client.ZCount(context.Background(),
		r.formatKey("stakingTxHistoryWithHash", args.Address, args.TxType),
		strconv.FormatInt(start, 10),
		strconv.FormatInt(stop, 10)).Result()
	if count != int64(args.PageSize) {
		return nil, nil
	}
	//
	var txsHash []string
	if args.Order == "ASC" {
		// get all txs
		txsHash, err = r.client.ZRange(context.Background(),
			r.formatKey("stakingTxHistoryWithHash", args.Address, args.TxType),
			start,
			stop).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// get all txs
		txsHash, err = r.client.ZRevRange(context.Background(),
			r.formatKey("stakingTxHistoryWithHash", args.Address, args.TxType),
			start,
			stop).Result()
		if err != nil {
			return nil, err
		}
	}

	// if not full txs
	if !args.FullTx {
		return map[string]interface{}{
			"staking_transactions": txsHash,
		}, nil
	}

	// result
	txsWithFullTx := []common.StakingTransaction{}
	for _, hash := range txsHash {
		var tx common.StakingTransaction
		r.client.HGetAll(context.Background(), r.formatKey("stakingTx", hash)).Scan(&tx)
		//if cmd.Val()[0] == nil {
		//	return nil, nil
		//}
		//json.Unmarshal(res, &tx)

		//r.client.
		//	tx = cmd.Val()[0].(common.StakingTransaction)

		txsWithFullTx = append(txsWithFullTx, tx)

	}

	return map[string]interface{}{
		"staking_transactions": txsWithFullTx,
	}, nil

	panic("not implemented")
}

func (r *RedisClient) SetTransactionsCount(address string, txType string, count *big.Int) error {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMSet(r.formatKey("txCount", txType), address, strconv.FormatInt(int64(count), 10))
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	pipe.HMSet(context.Background(), r.formatKey("txCount", address), map[string]interface{}{
		txType: count.String(),
	})
	pipe.Expire(context.Background(), r.formatKey("txCount", address), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetTransactionsCount(address, txType string) (*big.Int, error) {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMGet(r.formatKey("txCount", txType), address)
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	cmd := pipe.HMGet(context.Background(), r.formatKey("txCount", address), txType)
	pipe.Exec(context.Background())
	fmt.Println("REDIS >>>>>>>> ", cmd.Val())
	// len
	if cmd.Val()[0] == nil {
		return nil, nil
	}
	//
	result := new(big.Int)
	_, ok := result.SetString(cmd.Val()[0].(string), 10)
	if !ok {
		return nil, fmt.Errorf("[ %s ] : [ %s ]", "redis", "set big.Int failed")
	}
	return result, nil
}

// SetTransactionsHistory
// Deprecated. fixme
func (r *RedisClient) SetTransactionsHistory(address string, txs []common.Transaction) error {
	// init
	panic("not implemented")
}

// GetTransactionsHistory
// Deprecated. fixme
func (r *RedisClient) GetTransactionsHistory(args hmy.TxHistoryArgs) ([]common.Transaction, error) {
	panic("not implemented")
}

func (r *RedisClient) SetBlocks(blocks []common.BlockWithFullTx, blocksWithTxHash []common.BlockWithTxHash) error {
	// 获取 height
	// 记录 height => block
	// 记录 height => hash
	// 遍历
	// init
	//tx := r.client.Multi()
	//defer tx.Close()
	pipe := r.client.TxPipeline()
	defer pipe.Close()
	// 分开存放 block的数据
	for _, blk := range blocksWithTxHash {
		height := blk.Number.String()
		headerBytes, _ := json.Marshal(blk.Header)
		pipe.HMSet(context.Background(), r.formatKey("blockHeader"), map[string]interface{}{
			height: string(headerBytes),
		})
		// write uncles
		uncleBytes, _ := json.Marshal(blk.Uncles)
		pipe.HMSet(context.Background(), r.formatKey("blockUncles"), map[string]interface{}{
			height: string(uncleBytes),
		})
		// write signers
		signersByte, _ := json.Marshal(blk.Signers)
		pipe.HMSet(context.Background(), r.formatKey("blockSigners"), map[string]interface{}{
			height: string(signersByte),
		})
		// write txs hash
		txsHashBytes, _ := json.Marshal(blk.Transactions)
		pipe.HMSet(context.Background(), r.formatKey("blockTxsHash"), map[string]interface{}{
			height: string(txsHashBytes),
		})
		// write staking txs hash
		stakingTxsHashBytes, _ := json.Marshal(blk.StakingTxs)
		pipe.HMSet(context.Background(), r.formatKey("blockStakingTxsHash"), map[string]interface{}{
			height: string(stakingTxsHashBytes),
		})
		// write eth txs
		txInEthHashBytes, _ := json.Marshal(blk.EthTransactions)
		pipe.HMSet(context.Background(), r.formatKey("blockTxInEthHash"), map[string]interface{}{
			height: string(txInEthHashBytes),
		})
		// write hash to height
		pipe.HMSet(context.Background(), r.formatKey("hashToHeight"), map[string]interface{}{
			blk.Hash: height,
		})
	}

	for _, block := range blocks {
		// write txs
		for _, tx := range block.Transactions {
			txBytes, _ := json.Marshal(*tx)
			pipe.HMSet(context.Background(), r.formatKey("txs"), map[string]interface{}{
				tx.Hash: string(txBytes),
			})
		}
		// write stakingTx
		for _, stakingTx := range block.StakingTxs {
			stakingTxBytes, _ := json.Marshal(*stakingTx)
			pipe.HMSet(context.Background(), r.formatKey("stakingTxs"), map[string]interface{}{
				stakingTx.Hash: string(stakingTxBytes),
			})
		}
	}

	//for _, block := range blocks {
	//	// get height
	//	height := block.Number.Int64()
	//	// 拆分 txs
	//	var txsWithHash []string
	//	var stakingTxsWithHash []string
	//	var blk common.BlockWithTxHash
	//	// write txs
	//
	//	txs := block.Transactions
	//	for _, tx := range txs {
	//		pipe.HMSet(context.Background(), r.formatKey("txs", tx.Hash), *tx)
	//		txsWithHash = append(txsWithHash, tx.Hash)
	//	}
	//	for _, stakingTx := range block.StakingTxs {
	//		pipe.HMSet(context.Background(), r.formatKey("stakingTxs", stakingTx.Hash), *stakingTx)
	//		stakingTxsWithHash = append(stakingTxsWithHash, stakingTx.Hash)
	//	}
	//	// replace block.transactions with txsWithHash
	//	blk = common.BlockWithTxHash{
	//		Header:          block.Header,
	//		Uncles:          block.Uncles,
	//		Transactions:    txsWithHash,
	//		EthTransactions: nil,
	//		StakingTxs:      stakingTxsWithHash,
	//		Signers:         nil,
	//	}
	//	// write block
	//	pipe.HMSet(context.Background(), r.formatKey("blockNumber", strconv.FormatInt(height, 10)), block.Header)

	// write signers

	// write staking Tx

	cmds, err := pipe.Exec(context.Background())

	// for
	//cmds, err := tx.Exec(func() error {
	//	for _, block := range blocks {
	//		// get height
	//		height := block.Number.Int64()
	//		// write block
	//		r.client.HMSet(context.Background(), r.formatKey("blockNumber", strconv.FormatInt(height, 10)), block)
	//	}
	//})
	_ = cmds
	return err
}

func (r *RedisClient) hashToHeight(hash string) (height *rpc.BlockNumber) {
	//tx := r.client.Multi()
	//defer tx.Close()
	//cmd := tx.HMGet(r.formatKey("hashToHeight"), hash)
	// get
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("hashToHeight"), hash).Result()
	if cmd[0] == nil || err != nil {
		return nil
	}
	//fmt.Println(">??>><><>>>??>>><>>>>", cmd, height, err)
	json.Unmarshal([]byte(cmd[0].(string)), &height)

	return height
}

func (r *RedisClient) GetBlockByNumber(height rpc.BlockNumber, args hmy.BlockArgs) (interface{}, error) {

	heightString := strconv.FormatInt(int64(height), 10)
	// get header
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockHeader"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	//fmt.Println(">??>><><>>>??>>><>>>>", cmd, height, heightString, err)
	var header common.Header
	json.Unmarshal([]byte(cmd[0].(string)), &header)

	// get uncles
	cmd, err = r.client.HMGet(context.Background(), r.formatKey("blockUncles"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var uncles []string
	json.Unmarshal([]byte(cmd[0].(string)), &uncles)
	// fixme get signers
	cmd, err = r.client.HMGet(context.Background(), r.formatKey("blockSigners"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var signers []string
	json.Unmarshal([]byte(cmd[0].(string)), &signers)
	// fixme get blockTxInEthHash
	cmd, err = r.client.HMGet(context.Background(), r.formatKey("blockTxInEthHash"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var blockTxInEthHash []string
	json.Unmarshal([]byte(cmd[0].(string)), &blockTxInEthHash)
	// fixme get txsWithHash
	cmd, err = r.client.HMGet(context.Background(), r.formatKey("blockTxsHash"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var txsWithHash []string
	json.Unmarshal([]byte(cmd[0].(string)), &txsWithHash)
	// fixme get stakingTxsWithHash
	cmd, err = r.client.HMGet(context.Background(), r.formatKey("blockStakingTxsHash"), heightString).Result()

	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var stakingTxsHash []string
	json.Unmarshal([]byte(cmd[0].(string)), &stakingTxsHash)

	//
	if args.FullTx {
		var reply common.BlockWithFullTx
		var txs []*common.Transaction
		var stakingTxs []*common.StakingTransaction

		reply.Header = header
		reply.Uncles = uncles
		reply.Signers = []string{}
		reply.Transactions = []*common.Transaction{}
		reply.StakingTxs = []*common.StakingTransaction{}

		// get full tx
		for _, hash := range txsWithHash {
			cmd, err = r.client.HMGet(context.Background(), r.formatKey("txs"), hash).Result()
			if cmd[0] == nil || err != nil {
				return nil, nil
			}
			var tx common.Transaction
			json.Unmarshal([]byte(cmd[0].(string)), &tx)
			//
			txs = append(txs, &tx)
		}

		// get full staking Tx
		for _, hash := range stakingTxsHash {
			cmd, err = r.client.HMGet(context.Background(), r.formatKey("stakingTxs"), hash).Result()
			if cmd[0] == nil || err != nil {
				return nil, nil
			}
			// fixme staking tx 的msg字段是interface，msg结构体字段里有Amount字段，Amount应为 big.Int类型，
			// 但interface{}类型的msg在Marshal时会变为科学计数法，表现出的问题就是，Amount输出类似 "5e+23"，这块需要确认是否会有影响
			// hmy此处msg会有 CreateValidatorMsg、EditValidatorMsg、CollectRewardsMsg、DelegateMsg、UndelegateMsg，
			// 依据type改变msg的解析
			var stakingTx common.StakingTransaction
			json.Unmarshal([]byte(cmd[0].(string)), &stakingTx)
			//
			stakingTxs = append(stakingTxs, &stakingTx)
		}

		// if with staking txs
		if args.InclStaking {
			reply.StakingTxs = stakingTxs
		}

		// if with signers
		if args.WithSigners {
			reply.Signers = signers
		}

		return reply, nil
	} else {
		var reply common.BlockWithTxHash

		reply.Header = header
		reply.Uncles = uncles
		reply.Transactions = txsWithHash
		reply.EthTransactions = blockTxInEthHash

		reply.StakingTxs = []string{}
		// if with staking txs
		if args.InclStaking {
			reply.StakingTxs = stakingTxsHash
		}

		reply.Signers = []string{}
		// if with signers
		if args.WithSigners {
			reply.Signers = signers
		}
		// return block
		return reply, nil
	}

}

func (r *RedisClient) GetBlockByHash(hash string, args hmy.BlockArgs) (interface{}, error) {
	height := r.hashToHeight(hash)
	if height == nil {
		return nil, nil
	}
	return r.GetBlockByNumber(*height, args)
}

// todo 每个epoch都有相同的signers，整合这部分数据，不然数据量太大
func (r *RedisClient) SetBlockSigners(height rpc.BlockNumber, signers []string) error {
	// write signers
	heightString := strconv.FormatInt(int64(height), 10)
	signersByte, _ := json.Marshal(signers)
	r.client.HMSet(context.Background(), r.formatKey("blockSigners"), map[string]interface{}{
		heightString: string(signersByte),
	})
	return nil
}

// todo 每个epoch都有相同的signers，整合这部分数据，不然数据量太大
// fixme 设计一个查找表，找到当前高度对应的初始signers列表
func (r *RedisClient) GetBlockSigners(height rpc.BlockNumber) ([]string, error) {
	heightString := strconv.FormatInt(int64(height), 10)
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockSigners"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var signers []string
	err = json.Unmarshal([]byte(cmd[0].(string)), &signers)
	return signers, err
}

// todo 每个epoch都有相同的signer keys，整合这部分数据，不然数据量太大
func (r *RedisClient) SetBlockSignersKeys(height rpc.BlockNumber, signerKeys []string) error {
	// write signers
	heightString := strconv.FormatInt(int64(height), 10)
	signerKeysByte, _ := json.Marshal(signerKeys)
	r.client.HMSet(context.Background(), r.formatKey("blockSignerKeys"), map[string]interface{}{
		heightString: string(signerKeysByte),
	})
	return nil
}

// todo 每个epoch都有相同的signer keys，整合这部分数据，不然数据量太大
// fixme 设计一个查找表，找到当前高度对应的初始signer keys列表
func (r *RedisClient) GetBlockSignersKeys(height rpc.BlockNumber) ([]string, error) {
	heightString := strconv.FormatInt(int64(height), 10)
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockSignerKeys"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var signerKeysByte []string
	err = json.Unmarshal([]byte(cmd[0].(string)), &signerKeysByte)
	return signerKeysByte, err
}

func (r *RedisClient) GetBlockTxCountByNumber(height rpc.BlockNumber) (int, error) {
	heightString := strconv.FormatInt(int64(height), 10)
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockTxsHash"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return -1, nil
	}
	var txsWithHash []string
	err = json.Unmarshal([]byte(cmd[0].(string)), &txsWithHash)
	return len(txsWithHash), nil
}

func (r *RedisClient) GetBlockTxCountByHash(hash string) (int, error) {
	height := r.hashToHeight(hash)
	heightString := strconv.FormatInt(int64(*height), 10)
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockTxsHash"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return -1, nil
	}
	var txsWithHash []string
	err = json.Unmarshal([]byte(cmd[0].(string)), &txsWithHash)
	return len(txsWithHash), nil
}

func (r *RedisClient) SetHeaderInfo(height rpc.BlockNumber, header *rpc.HeaderInformation) error {
	heightString := strconv.FormatInt(int64(height), 10)
	headerBytes, err := json.Marshal(*header)
	if err != nil {
		return err
	}
	err = r.client.HMSet(context.Background(), r.formatKey("blockHeaderInfo"), map[string]interface{}{
		heightString: headerBytes,
	}).Err()
	return err
}

func (r *RedisClient) GetHeaderByNumber(height rpc.BlockNumber) (*rpc.HeaderInformation, error) {
	heightString := strconv.FormatInt(int64(height), 10)
	// get header
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("blockHeaderInfo"), heightString).Result()
	if cmd[0] == nil || err != nil {
		return nil, nil
	}
	var header rpc.HeaderInformation
	err = json.Unmarshal([]byte(cmd[0].(string)), &header)
	return &header, err
}

func (r *RedisClient) GetHeaderByHash(hash string) (*rpc.HeaderInformation, error) {
	height := r.hashToHeight(hash)
	return r.GetHeaderByNumber(*height)
}

func (r *RedisClient) GetLatestChainHeaders() (*rpc.HeaderPair, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("latestChainHeaders"), "latest").Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var headerPair *rpc.HeaderPair
	json.Unmarshal([]byte(cmds[0].(string)), &headerPair)
	return headerPair, nil
}

func (r *RedisClient) SetLatestChainHeaders(headerPair *rpc.HeaderPair) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	headerPairBytes, err := json.Marshal(*headerPair)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("latestChainHeaders"), map[string]interface{}{
		"latest": headerPairBytes,
	})
	pipe.Expire(context.Background(), r.formatKey("latestChainHeaders"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetLatestHeader() (*rpc.HeaderInformation, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("blockLatestHeader"), "latest").Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var headerPair *rpc.HeaderInformation
	json.Unmarshal([]byte(cmds[0].(string)), &headerPair)
	return headerPair, nil
}

func (r *RedisClient) SetLatestHeader(headerPair *rpc.HeaderInformation) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	headerInfoBytes, err := json.Marshal(*headerPair)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("blockLatestHeader"), map[string]interface{}{
		"latest": headerInfoBytes,
	})
	pipe.Expire(context.Background(), r.formatKey("blockLatestHeader"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetBlockNumber() (*rpc.BlockNumber, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("blockNumber"), "latest").Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var height *rpc.BlockNumber
	json.Unmarshal([]byte(cmds[0].(string)), &height)
	return height, nil
}

func (r *RedisClient) SetBlockNumber(height *rpc.BlockNumber) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	heightBytes, err := json.Marshal(*height)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("blockNumber"), map[string]interface{}{
		"latest": heightBytes,
	})
	pipe.Expire(context.Background(), r.formatKey("blockNumber"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetCirculatingSupply() { panic("not implemented") }
func (r *RedisClient) SetCirculatingSupply() { panic("not implemented") }

func (r *RedisClient) GetEpoch() { panic("not implemented") }
func (r *RedisClient) SetEpoch() { panic("not implemented") }

func (r *RedisClient) GetLastCrossLinks() { panic("not implemented") }
func (r *RedisClient) SetLastCrossLinks() { panic("not implemented") }

func (r *RedisClient) GetLeader() { panic("not implemented") }
func (r *RedisClient) SetLeader() { panic("not implemented") }

func (r *RedisClient) GetShardingStructure() { panic("not implemented") }
func (r *RedisClient) SetShardingStructure() { panic("not implemented") }

func (r *RedisClient) GetTotalSupply() { panic("not implemented") }
func (r *RedisClient) SetTotalSupply() { panic("not implemented") }

func (r *RedisClient) GetValidators(epoch uint64) (*rpc.ValidatorsFields, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("validators"), strconv.FormatUint(epoch, 10)).Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var validators *rpc.ValidatorsFields
	json.Unmarshal([]byte(cmds[0].(string)), &validators)
	return validators, nil
}

func (r *RedisClient) SetValidators(epoch uint64, height rpc.BlockNumber, validators *rpc.ValidatorsFields) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	validatorsBytes, err := json.Marshal(*validators)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("validators"), map[string]interface{}{
		strconv.FormatUint(epoch, 10): validatorsBytes,
	})
	// todo 过期时间计算方法：当前block height 计算出下一个epoch的时间，作为timeout
	pipe.Expire(context.Background(), r.formatKey("validators"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetValidatorKeys(epoch uint64) ([]string, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("validatorKeys"), strconv.FormatUint(epoch, 10)).Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var validatorKeys []string
	json.Unmarshal([]byte(cmds[0].(string)), &validatorKeys)
	return validatorKeys, nil
}

func (r *RedisClient) SetValidatorKeys(epoch uint64, height rpc.BlockNumber, validators []string) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	validatorsBytes, err := json.Marshal(validators)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("validatorKeys"), map[string]interface{}{
		strconv.FormatUint(epoch, 10): validatorsBytes,
	})
	// todo 过期时间计算方法：当前block height 计算出下一个epoch的时间，作为timeout
	pipe.Expire(context.Background(), r.formatKey("validatorKeys"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetMetadata() (*common.NodeMetadata, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("nodeMetadata"), "latest").Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var metadata *common.NodeMetadata
	json.Unmarshal([]byte(cmds[0].(string)), &metadata)
	return metadata, nil
}

func (r *RedisClient) SetMetadata(metadata *common.NodeMetadata) error {
	pipe := r.client.TxPipeline()
	defer pipe.Close()

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	pipe.HMSet(context.Background(), r.formatKey("nodeMetadata"), map[string]interface{}{
		"latest": metadataBytes,
	})
	// todo 过期时间计算方法：当前block height 计算出下一个epoch的时间，作为timeout
	pipe.Expire(context.Background(), r.formatKey("nodeMetadata"), r.timeout)
	cmds, err := pipe.Exec(context.Background())
	_ = cmds
	return err
}

func (r *RedisClient) GetProtocolVersion() { panic("not implemented") }
func (r *RedisClient) SetProtocolVersion() { panic("not implemented") }

func (r *RedisClient) GetPeerCount() { panic("not implemented") }
func (r *RedisClient) SetPeerCount() { panic("not implemented") }

func (r *RedisClient) GetCode(address string) (*string, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("contractCode"), address).Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var code *string
	json.Unmarshal([]byte(cmds[0].(string)), &code)
	return code, nil
}

func (r *RedisClient) SetCode(address, code string) error {
	codeBytes, err := json.Marshal(code)
	if err != nil {
		return err
	}
	err = r.client.HMSet(context.Background(), r.formatKey("contractCode"), map[string]interface{}{
		address: codeBytes,
	}).Err()
	return err
}

func (r *RedisClient) GetDelegationsByDelegator() { panic("not implemented") }
func (r *RedisClient) SetDelegationsByDelegator() { panic("not implemented") }

func (r *RedisClient) GetDelegationsByDelegatorByBlockNumber() { panic("not implemented") }
func (r *RedisClient) SetDelegationsByDelegatorByBlockNumber() { panic("not implemented") }

func (r *RedisClient) GetDelegationsByValidator() { panic("not implemented") }
func (r *RedisClient) SetDelegationsByValidator() { panic("not implemented") }

func (r *RedisClient) GetAllValidatorAddresses() { panic("not implemented") }
func (r *RedisClient) SetAllValidatorAddresses() { panic("not implemented") }

func (r *RedisClient) GetAllValidatorInformation() { panic("not implemented") }
func (r *RedisClient) SetAllValidatorInformation() { panic("not implemented") }

func (r *RedisClient) GetAllValidatorInformationByBlockNumber() { panic("not implemented") }
func (r *RedisClient) SetAllValidatorInformationByBlockNumber() { panic("not implemented") }

func (r *RedisClient) GetDetElectedValidatorAddresses() { panic("not implemented") }
func (r *RedisClient) SetElectedValidatorAddresses()    { panic("not implemented") }

func (r *RedisClient) GetValidatorInformation() { panic("not implemented") }
func (r *RedisClient) SetValidatorInformation() { panic("not implemented") }

func (r *RedisClient) GetCurrentUtilityMetrics() { panic("not implemented") }
func (r *RedisClient) SetCurrentUtilityMetrics() { panic("not implemented") }

func (r *RedisClient) GetMedianRawStakeSnapshot() { panic("not implemented") }
func (r *RedisClient) SetMedianRawStakeSnapshot() { panic("not implemented") }

func (r *RedisClient) GetStakingNetworkInfo() { panic("not implemented") }
func (r *RedisClient) SetStakingNetworkInfo() { panic("not implemented") }

func (r *RedisClient) GetSuperCommittees() { panic("not implemented") }
func (r *RedisClient) SetSuperCommittees() { panic("not implemented") }

func (r *RedisClient) GetCXReceiptByHash() { panic("not implemented") }
func (r *RedisClient) SetCXReceiptByHash() { panic("not implemented") }

// can not be implemented
func (r *RedisClient) GetPendingCXReceipts() { panic("not implemented") }
func (r *RedisClient) SetPendingCXReceipts() { panic("not implemented") }

// can not be implemented
func (r *RedisClient) GetPoolStats() { panic("not implemented") }
func (r *RedisClient) SetPoolStats() { panic("not implemented") }

// can not be implemented
func (r *RedisClient) GetPendingStakingTransactions() { panic("not implemented") }
func (r *RedisClient) SetPendingStakingTransactions() { panic("not implemented") }

func (r *RedisClient) GetCurrentStakingErrorSink() { panic("not implemented") }
func (r *RedisClient) SetCurrentStakingErrorSink() { panic("not implemented") }

func (r *RedisClient) GetStakingTransactionByBlockNumberAndIndex() { panic("not implemented") }
func (r *RedisClient) SetStakingTransactionByBlockNumberAndIndex() { panic("not implemented") }

func (r *RedisClient) GetStakingTransactionByBlockHashAndIndex() { panic("not implemented") }
func (r *RedisClient) SetStakingTransactionByBlockHashAndIndex() { panic("not implemented") }

func (r *RedisClient) GetStakingTransactionByHash(hash string) (*common.StakingTransaction, error) {
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("stakingTxs"), hash).Result()
	if cmd[0] == nil || err != nil {
		return nil, err
	}
	var tx *common.StakingTransaction
	json.Unmarshal([]byte(cmd[0].(string)), &tx)
	return tx, nil
}
func (r *RedisClient) SetStakingTransaction(tx *common.StakingTransaction) error {
	txBytes, _ := json.Marshal(*tx)
	err := r.client.HMSet(context.Background(), r.formatKey("stakingTxs"), map[string]interface{}{
		tx.Hash: string(txBytes),
	}).Err()
	return err
}

func (r *RedisClient) GetCurrentTransactionErrorSink() { panic("not implemented") }
func (r *RedisClient) SetCurrentTransactionErrorSink() { panic("not implemented") }

func (r *RedisClient) GetTransactionByBlockHashAndIndex(hash string, index int64) (*common.Transaction, error) {
	return r.GetTransactionByBlockNumberAndIndex(*r.hashToHeight(hash), index)
}

func (r *RedisClient) SetTransactionByBlockHashAndIndex(tx *common.Transaction) error {
	// write hash to height
	r.client.HMSet(context.Background(), r.formatKey("hashToHeight"), map[string]interface{}{
		tx.BlockHash: tx.BlockNumber.String(),
	})

	return r.SetTransactionByBlockNumberAndIndex(tx)
}

func (r *RedisClient) GetTransactionByBlockNumberAndIndex(height rpc.BlockNumber, index int64) (*common.Transaction, error) {
	var hash []string
	hash, err := r.client.ZRangeByScore(context.Background(),
		r.formatKey("txHashByNumberAndIndex", strconv.FormatInt(int64(height), 10)),
		&redis.ZRangeBy{
			Min:    strconv.FormatInt(index, 10),
			Max:    strconv.FormatInt(index, 10),
			Offset: 0,
			Count:  1,
		}).Result()
	if len(hash) == 0 || err != nil {
		return nil, err
	}
	return r.GetTransactionByHash(hash[0])
}

func (r *RedisClient) SetTransactionByBlockNumberAndIndex(tx *common.Transaction) error {
	_, err := r.client.ZAdd(context.Background(), r.formatKey("txHashByNumberAndIndex", tx.BlockNumber.String()), &redis.Z{
		Score:  float64(tx.TransactionIndex),
		Member: tx.Hash,
	}).Result()
	return err
}

func (r *RedisClient) GetTransactionByHash(hash string) (*common.Transaction, error) {
	cmd, err := r.client.HMGet(context.Background(), r.formatKey("txs"), hash).Result()
	if cmd[0] == nil || err != nil {
		return nil, err
	}
	var tx *common.Transaction
	json.Unmarshal([]byte(cmd[0].(string)), &tx)
	return tx, nil
}

func (r *RedisClient) SetTransaction(tx *common.Transaction) error {
	txBytes, _ := json.Marshal(*tx)
	err := r.client.HMSet(context.Background(), r.formatKey("txs"), map[string]interface{}{
		tx.Hash: string(txBytes),
	}).Err()
	return err
}

func (r *RedisClient) GetTransactionReceipt(hash string) (*common.TxReceipt, error) {
	cmds, err := r.client.HMGet(context.Background(), r.formatKey("txReceipts"), hash).Result()
	if cmds[0] == nil || err != nil {
		return nil, err
	}
	var receipt *common.TxReceipt
	json.Unmarshal([]byte(cmds[0].(string)), &receipt)
	return receipt, nil
}

func (r *RedisClient) SetTransactionReceipt(hash string, receipt *common.TxReceipt) error {
	receiptBytes, err := json.Marshal(*receipt)
	if err != nil {
		return err
	}
	err = r.client.HMSet(context.Background(), r.formatKey("txReceipts"), map[string]interface{}{
		hash: receiptBytes,
	}).Err()
	return err
}

// other implementations ...
