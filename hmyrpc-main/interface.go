// interface.go provides access to the back-end rpc interface.

package main

import (
	"context"
	"math/big"

	"github.com/shenzhendev/hmyrpc/common"
	"github.com/shenzhendev/hmyrpc/rpc"
	"github.com/shenzhendev/hmyrpc/rpc/hmy"
)

// AccountReader is the account reader interface
type AccountReader interface {
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	GetBalanceByBlockNumber(ctx context.Context, address string, number int64) (*big.Int, error)
	GetStakingTransactionsCount(ctx context.Context, address string, txType string) (int64, error)
	GetStakingTransactionsHistory(ctx context.Context, args *hmy.TxHistoryArgs) (map[string]interface{}, error)
	GetTransactionsCount(ctx context.Context, address string, txType string) (int64, error)
	GetTransactionsHistory(ctx context.Context, args *hmy.TxHistoryArgs) (map[string]interface{}, error)
}

// BlocksReader is the blocks/ node/ network reader interface
type BlocksReader interface {
	GetBlocks(ctx context.Context, from, to rpc.BlockNumber, blockArgs *hmy.BlockArgs) ([]interface{}, error)
	GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, blockArgs *hmy.BlockArgs) (interface{}, error)
	GetBlockByHash(ctx context.Context, hash string, blockArgs *hmy.BlockArgs) (interface{}, error)
	GetBlockSigners(ctx context.Context, number rpc.BlockNumber) ([]string, error)
	GetBlockSignersKeys(ctx context.Context, number rpc.BlockNumber) ([]string, error)
	GetBlockTransactionCountByNumber(ctx context.Context, number rpc.BlockNumber) (int64, error)
	GetBlockTransactionCountByHash(ctx context.Context, hash string) (int64, error)
	GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*rpc.HeaderInformation, error)
	GetLatestChainHeaders(ctx context.Context) (*rpc.HeaderPair, error)
	LatestHeader(ctx context.Context) (*rpc.HeaderInformation, error)

	BlockNumber(ctx context.Context) (int64, error)
	GetCirculatingSupply(ctx context.Context) (string, error)
	GetEpoch(ctx context.Context) (uint64, error)
	GetLastCrossLinks(ctx context.Context) ([]*common.CrossLink, error)
	GetLeader(ctx context.Context) (string, error)
	GasPrice(ctx context.Context) (uint64, error)
	GetShardingStructure(ctx context.Context) ([]*rpc.ShardingStructure, error)
	GetTotalSupply(ctx context.Context) (string, error)
	GetValidators(ctx context.Context, epoch uint64) ([]*rpc.ValidatorsFields, error)
	GetValidatorKeys(ctx context.Context, epoch uint64) ([]string, error)

	GetNodeMetadata(ctx context.Context) (*common.NodeMetadata, error)
	ProtocolVersion(ctx context.Context) (int64, error)

	PeerCount(ctx context.Context) (string, error)
}

type ContractReader interface {
	Call(ctx context.Context, args rpc.CallArgs, number rpc.BlockNumber) (string, error)
	EstimateGas(ctx context.Context, args rpc.CallArgs) (uint64, error)
	GetCode(ctx context.Context, address string, number rpc.BlockNumber) (string, error)
	GetStorageAt(ctx context.Context, address string, key string, number rpc.BlockNumber) (string, error)
}

// StakingReader is the staking relative interface
type StakingReader interface {
	GetDelegationsByDelegator(ctx context.Context, address string) ([]*rpc.Delegation, error)
	GetDelegationsByDelegatorByBlockNumber(ctx context.Context, address rpc.AddressOrList, number rpc.BlockNumber) (interface{}, error)
	GetDelegationsByValidator(ctx context.Context, address string) ([]*rpc.Delegation, error)

	GetAllValidatorAddresses(ctx context.Context) ([]string, error)
	GetAllValidatorInformation(ctx context.Context, page int) ([]*common.ValidatorRPCEnhanced, error)
	GetAllValidatorInformationByBlockNumber(ctx context.Context, page int, number rpc.BlockNumber) ([]common.ValidatorRPCEnhanced, error)
	// GetElectedValidatorAddresses returns the address of elected validators for current epoch
	GetElectedValidatorAddresses(ctx context.Context) ([]string, error)
	// GetValidatorInformation returns the information of validator
	GetValidatorInformation(ctx context.Context, address string) (*common.ValidatorRPCEnhanced, error)

	// GetCurrentUtilityMetrics returns the ...
	GetCurrentUtilityMetrics(ctx context.Context) (*common.UtilityMetric, error)
	// GetMedianRawStakeSnapshot returns the ...
	GetMedianRawStakeSnapshot(ctx context.Context) (*common.CompletedEPoSRound, error)
	// GetStakingNetworkInfo returns the ...
	GetStakingNetworkInfo(ctx context.Context) (*rpc.StakingNetworkInfo, error)
	// GetSuperCommittees returns the ...
	GetSuperCommittees(ctx context.Context) (*common.Transition, error)
}

type TransactionReader interface {
	// GetCXReceiptByHash returns the ...
	GetCXReceiptByHash(ctx context.Context, hash string) (*common.CxReceipt, error)
	// GetPendingCXReceipts returns the ...
	GetPendingCXReceipts(ctx context.Context) ([]*common.CXReceiptsProof, error)
	// ResendCx returns the ...
	ResendCx(ctx context.Context, hash string) (bool, error)

	// GetPoolStats returns the ...
	GetPoolStats(ctx context.Context) (*rpc.PendingPoolStats, error)
	// PendingStakingTransactions returns the ...
	PendingStakingTransactions(ctx context.Context) ([]*common.StakingTransaction, error)
	// PendingTransactions returns the ...
	PendingTransactions(ctx context.Context) ([]*common.Transaction, error)

	// GetCurrentStakingErrorSink returns the ...
	GetCurrentStakingErrorSink(ctx context.Context) ([]*common.TransactionErrorReport, error)
	// GetStakingTransactionByBlockNumberAndIndex returns the ...
	GetStakingTransactionByBlockNumberAndIndex(ctx context.Context, number rpc.BlockNumber, index uint) (*common.StakingTransaction, error)
	// GetStakingTransactionByBlockHashAndIndex returns the ...
	GetStakingTransactionByBlockHashAndIndex(ctx context.Context, hash string, index uint) (*common.StakingTransaction, error)
	// GetStakingTransactionByHash return the ...
	GetStakingTransactionByHash(ctx context.Context, hash string) (*common.StakingTransaction, error)
	// SendRawStakingTransaction returns the ...
	SendRawStakingTransaction(ctx context.Context, rawTx string) (string, error)

	// GetCurrentTransactionErrorSink returns the ...
	GetCurrentTransactionErrorSink(ctx context.Context) ([]*common.TransactionErrorReport, error)
	// GetTransactionByBlockHashAndIndex returns the ...
	GetTransactionByBlockHashAndIndex(ctx context.Context, hash string, index uint) (*common.Transaction, error)
	// GetTransactionByBlockNumberAndIndex returns the ...
	GetTransactionByBlockNumberAndIndex(ctx context.Context, number rpc.BlockNumber, index uint) (*common.Transaction, error)
	// GetTransactionByHash returns the ...
	GetTransactionByHash(ctx context.Context, hash string) (*common.Transaction, error)
	// GetTransactionReceipt returns the ...
	GetTransactionReceipt(ctx context.Context, hash string) (*common.TxReceipt, error)
	// SendRawTransaction returns the ...
	SendRawTransaction(ctx context.Context, rawTx string) (string, error)
}
