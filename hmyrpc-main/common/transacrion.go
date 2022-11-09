package common

import (
	"math/big"

	"github.com/shenzhendev/hmyrpc/util"
)

// Transaction represents a transaction that will serialize to the RPC representation of a transaction
type Transaction struct {
	BlockHash        string   `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	From             string   `json:"from"`
	Timestamp        uint64   `json:"timestamp"`
	Gas              uint64   `json:"gas"`
	GasPrice         *big.Int `json:"gasPrice"`
	Hash             string   `json:"hash"`
	EthHash          string   `json:"ethHash"`
	Input            string   `json:"input"`
	Nonce            uint64   `json:"nonce"`
	To               string   `json:"to"`
	TransactionIndex uint64   `json:"transactionIndex"`
	Value            *big.Int `json:"value"`
	ShardID          uint32   `json:"shardID"`
	ToShardID        uint32   `json:"toShardID"`
	V                string   `json:"v"`
	R                string   `json:"r"`
	S                string   `json:"s"`
}

// StakingTransaction represents a transaction that will serialize to the RPC representation of a staking transaction
type StakingTransaction struct {
	BlockHash        string      `json:"blockHash"`
	BlockNumber      *big.Int    `json:"blockNumber"`
	From             string      `json:"from"`
	Timestamp        uint64      `json:"timestamp"`
	Gas              uint64      `json:"gas"`
	GasPrice         *big.Int    `json:"gasPrice"`
	Hash             string      `json:"hash"`
	Nonce            uint64      `json:"nonce"`
	TransactionIndex uint64      `json:"transactionIndex"`
	V                string      `json:"v"`
	R                string      `json:"r"`
	S                string      `json:"s"`
	Type             string      `json:"type"`
	Msg              interface{} `json:"msg"`
}

// TxReceipt represents a transaction receipt that will serialize to the RPC representation.
type TxReceipt struct {
	BlockHash         string `json:"blockHash"`
	TransactionHash   string `json:"transactionHash"`
	BlockNumber       uint64 `json:"blockNumber"`
	TransactionIndex  uint64 `json:"transactionIndex"`
	GasUsed           uint64 `json:"gasUsed"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []*Log `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	ShardID           uint32 `json:"shardID"`
	From              string `json:"from"`
	To                string `json:"to"`
	Root              []byte `json:"root"`
	Status            uint   `json:"status"`
}

// StakingTxReceipt represents a staking transaction receipt that will serialize to the RPC representation.
type StakingTxReceipt struct {
	BlockHash         string `json:"blockHash"`
	TransactionHash   string `json:"transactionHash"`
	BlockNumber       uint64 `json:"blockNumber"`
	TransactionIndex  uint64 `json:"transactionIndex"`
	GasUsed           uint64 `json:"gasUsed"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []*Log `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Sender            string `json:"sender"`
	Type              string `json:"type"`
	Root              []byte `json:"root"`
	Status            uint   `json:"status"`
}

// CxReceipt represents a CxReceipt that will serialize to the RPC representation of a CxReceipt
type CxReceipt struct {
	BlockHash   string   `json:"blockHash"`
	BlockNumber *big.Int `json:"blockNumber"`
	TxHash      string   `json:"hash"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	ShardID     uint32   `json:"shardID"`
	ToShardID   uint32   `json:"toShardID"`
	Amount      *big.Int `json:"value"`
}

const ReceiptStatusFailed = uint(0)
const ReceiptStatusSuccessful = uint(1)

func (r *TxReceipt) Confirmed() bool {
	return !util.IsZeroHash(r.BlockHash)
}

func (r *TxReceipt) Successful() bool {
	return r.Status == ReceiptStatusSuccessful
}

// CXReceiptsProof carrys the cross shard receipts and merkle proof
type CXReceiptsProof struct {
	Receipts    []*CXReceipt   `json:"receipts"     gencodec:"required"`
	MerkleProof *CXMerkleProof `json:"merkleProof"  gencodec:"required"`
	//Header       *block.Header  `json:"header"       gencoded:"required"`
	CommitSig    []byte `json:"commitSig"`
	CommitBitmap []byte `json:"commitBitmap"`
}

// CXReceipt represents a receipt for cross-shard transaction
type CXReceipt struct {
	TxHash    string   `json:"hash"` // hash of the cross shard transaction in source shard
	From      string   `json:"from"`
	To        string   `json:"to"`
	ShardID   uint32   `json:"shardID"`
	ToShardID uint32   `json:"toShardID"`
	Amount    *big.Int `json:"value"`
}

// CXMerkleProof represents the merkle proof of a collection of ordered cross shard transactions
type CXMerkleProof struct {
	BlockNum      *big.Int `json:"blockNum"`
	BlockHash     string   `json:"blockHash"   gencodec:"required"`
	ShardID       uint32   `json:"shardID"`
	CXReceiptHash string   `json:"receiptHash" gencodec:"required"`
	ShardIDs      []uint32 `json:"shardIDs"`
	CXShardHashes []string `json:"shardHashes" gencodec:"required"`
}

// TransactionErrorReport ..
type TransactionErrorReport struct {
	TxHashID             string `json:"tx-hash-id"`
	StakingDirective     string `json:"directive-kind,omitempty"`
	TimestampOfRejection int64  `json:"time-at-rejection"`
	ErrMessage           string `json:"error-message"`
}
