package common

import (
	"math/big"
)

type Uint64 uint64

// BlockWithTxHash represents a block that will serialize to the RPC representation of a block
// having ONLY transaction hashes in the Transaction & Staking transaction fields.
type BlockWithTxHash struct {
	Header
	Uncles          []string `json:"uncles"`
	Transactions    []string `json:"transactions"`
	EthTransactions []string `json:"transactionsInEthHash"`
	StakingTxs      []string `json:"stakingTransactions"`
	Signers         []string `json:"signers,omitempty"`
}

// BlockWithFullTx represents a block that will serialize to the RPC representation of a block
// having FULL transactions in the Transaction & Staking transaction fields.
type BlockWithFullTx struct {
	Header
	Uncles       []string              `json:"uncles"`
	Transactions []*Transaction        `json:"transactions"`
	StakingTxs   []*StakingTransaction `json:"stakingTransactions"`
	Signers      []string              `json:"signers,omitempty"`
}

type Header struct {
	Number           *big.Int `json:"number"`
	ViewID           *big.Int `json:"viewID"`
	Epoch            *big.Int `json:"epoch"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            uint64   `json:"nonce"`
	MixHash          string   `json:"mixHash"`
	LogsBloom        string   `json:"logsBloom"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       uint64   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	Size             uint64   `json:"size"`
	GasLimit         uint64   `json:"gasLimit"`
	GasUsed          uint64   `json:"gasUsed"`
	VRF              string   `json:"vrf"`
	VRFProof         string   `json:"vrfProof"`
	Timestamp        *big.Int `json:"timestamp"`
	TransactionsRoot string   `json:"transactionsRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
}

type T struct {
	BlockHash        string        `json:"blockHash"`
	BlockNumber      int           `json:"blockNumber"`
	CrossLinks       []interface{} `json:"crossLinks"`
	Epoch            int           `json:"epoch"`
	LastCommitBitmap string        `json:"lastCommitBitmap"`
	LastCommitSig    string        `json:"lastCommitSig"`
	Leader           string        `json:"leader"`
	ShardID          int           `json:"shardID"`
	Timestamp        string        `json:"timestamp"`
	Unixtime         int           `json:"unixtime"`
	ViewID           int           `json:"viewID"`
	Vrf              string        `json:"vrf"`
	VrfProof         string        `json:"vrfProof"`
}

type LatestHeader struct {
	Difficulty       string   `json:"difficulty"`
	Epoch            *big.Int `json:"epoch"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	ShardID          uint32   `json:"shardID"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TransactionsRoot string   `json:"transactionsRoot"`
	ViewID           *big.Int `json:"viewID"`
}
