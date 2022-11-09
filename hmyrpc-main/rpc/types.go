package rpc

import (
	"bytes"
	"math/big"

	jsoniter "github.com/json-iterator/go"
	"github.com/shenzhendev/hmyrpc/common"
)

type BlockNumber int64

// StructuredResponse type of RPCs
type StructuredResponse = map[string]interface{}

// NewStructuredResponse creates a structured response from the given input
func NewStructuredResponse(input interface{}) (StructuredResponse, error) {
	var objMap StructuredResponse
	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	dat, err := jsonIter.Marshal(input)
	if err != nil {
		return nil, err
	}
	d := jsonIter.NewDecoder(bytes.NewReader(dat))
	d.UseNumber()
	err = d.Decode(&objMap)
	if err != nil {
		return nil, err
	}
	return objMap, nil
}

const (
	PendingBlockNumber  = BlockNumber(-2)
	LatestBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber = BlockNumber(0)
)

// HeaderInformation represents the latest consensus information
type HeaderInformation struct {
	BlockHash        string             `json:"blockHash"`
	BlockNumber      uint64             `json:"blockNumber"`
	ShardID          uint32             `json:"shardID"`
	Leader           string             `json:"leader"`
	ViewID           uint64             `json:"viewID"`
	Epoch            uint64             `json:"epoch"`
	Timestamp        string             `json:"timestamp"`
	UnixTime         uint64             `json:"unixtime"`
	LastCommitSig    string             `json:"lastCommitSig"`
	LastCommitBitmap string             `json:"lastCommitBitmap"`
	VRF              string             `json:"vrf"`
	VRFProof         string             `json:"vrfProof"`
	CrossLinks       *common.CrossLinks `json:"crossLinks,omitempty"`
}

// HeaderPair ..
type HeaderPair struct {
	BeaconHeader *common.LatestHeader `json:"beacon-chain-header"`
	ShardHeader  *common.LatestHeader `json:"shard-chain-header"`
}

type ShardingStructure struct {
	WS      string `json:"ws"`
	Http    string `json:"http"`
	ShardID int    `json:"shard"`
	Current bool   `json:"current"`
}

type ValidatorsFields struct {
	ShardID    int          `json:"shardID"`
	Validators []*Validator `json:"validators"`
}

type Validator struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     *string  `json:"from"`
	To       *string  `json:"to"`
	Gas      *uint64  `json:"gas"`
	GasPrice *big.Int `json:"gasPrice"`
	Value    *big.Int `json:"value"`
	Data     *string  `json:"data"`
}

// Delegation represents a particular delegation to a validator
type Delegation struct {
	ValidatorAddress string         `json:"validator_address"`
	DelegatorAddress string         `json:"delegator_address"`
	Amount           *big.Int       `json:"amount"`
	Reward           *big.Int       `json:"reward"`
	Undelegations    []Undelegation `json:"Undelegations"`
}

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount *big.Int
	Epoch  *big.Int
}

// AddressOrList represents an address or a list of addresses
type AddressOrList struct {
	Address     *string
	AddressList []string
}

// StakingNetworkInfo returns global staking info.
type StakingNetworkInfo struct {
	TotalSupply       string   `json:"total-supply"`
	CirculatingSupply string   `json:"circulating-supply"`
	EpochLastBlock    uint64   `json:"epoch-last-block"`
	TotalStaking      *big.Int `json:"total-staking"`
	MedianRawStake    string   `json:"median-raw-stake"`
}

type PendingPoolStats struct {
	ExeCount    int `json:"executable-count"`
	NonExeCount int `json:"non-executable-count"`
}

type TxHistoryWithHash struct {
	Transactions []string `json:"transactions"`
}

type TxHistoryWithFullTx struct {
	Transactions []common.Transaction `json:"transactions"`
}

type StakingTxHistoryWithHash struct {
	StakingTransactions []string `json:"staking_transactions"`
}

type StakingTxHistoryWithFullTx struct {
	StakingTransactions []common.StakingTransaction `json:"staking_transactions"`
}
