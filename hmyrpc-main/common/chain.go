package common

import "math/big"

// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	// ChainId identifies the current chain and is used for replay protection
	ChainID *big.Int `json:"chain-id"`

	// EthCompatibleChainID identifies the chain id used for ethereum compatible transactions
	EthCompatibleChainID *big.Int `json:"eth-compatible-chain-id"`

	// EthCompatibleShard0ChainID identifies the shard 0 chain id used for ethereum compatible transactions
	EthCompatibleShard0ChainID *big.Int `json:"eth-compatible-shard-0-chain-id"`

	// EthCompatibleEpoch is the epoch where ethereum-compatible transaction starts being processed.
	EthCompatibleEpoch *big.Int `json:"eth-compatible-epoch,omitempty"`

	// CrossTxEpoch is the epoch where cross-shard transaction starts being processed.
	CrossTxEpoch *big.Int `json:"cross-tx-epoch,omitempty"`

	// CrossLinkEpoch is the epoch where beaconchain starts containing cross-shard links.
	CrossLinkEpoch *big.Int `json:"cross-link-epoch,omitempty"`

	// AggregatedRewardEpoch is the epoch when block rewards are distributed every 64 blocks
	AggregatedRewardEpoch *big.Int `json:"aggregated-reward-epoch,omitempty"`

	// StakingEpoch is the epoch when shard assign takes staking into account
	StakingEpoch *big.Int `json:"staking-epoch,omitempty"`

	// PreStakingEpoch is the epoch we allow staking transactions
	PreStakingEpoch *big.Int `json:"prestaking-epoch,omitempty"`

	// QuickUnlockEpoch is the epoch when undelegation will be unlocked at the current epoch
	QuickUnlockEpoch *big.Int `json:"quick-unlock-epoch,omitempty"`

	// FiveSecondsEpoch is the epoch when block time is reduced to 5 seconds and block rewards adjusted to 17.5 ONE/block
	FiveSecondsEpoch *big.Int `json:"five-seconds-epoch,omitempty"`

	// TwoSecondsEpoch is the epoch when block time is reduced to 2 seconds and block rewards adjusted to 7 ONE/block
	TwoSecondsEpoch *big.Int `json:"two-seconds-epoch,omitempty"`

	// SixtyPercentEpoch is the epoch when internal voting power reduced from 68% to 60%
	SixtyPercentEpoch *big.Int `json:"sixty-percent-epoch,omitempty"`

	// RedelegationEpoch is the epoch when redelegation is supported and undelegation locking time is restored to 7 epoch
	RedelegationEpoch *big.Int `json:"redelegation-epoch,omitempty"`

	// NoEarlyUnlockEpoch is the epoch when the early unlock of undelegated token from validators who were elected for
	// more than 7 epochs is disabled
	NoEarlyUnlockEpoch *big.Int `json:"no-early-unlock-epoch,omitempty"`

	// VRFEpoch is the epoch when VRF randomness is enabled
	VRFEpoch *big.Int `json:"vrf-epoch,omitempty"`

	// PrevVRFEpoch is the epoch when previous VRF randomness can be fetched
	PrevVRFEpoch *big.Int `json:"prev-vrf-epoch,omitempty"`

	// MinDelegation100Epoch is the epoch when min delegation is reduced from 1000 ONE to 100 ONE
	MinDelegation100Epoch *big.Int `json:"min-delegation-100-epoch,omitempty"`

	// MinCommissionRateEpoch is the epoch when policy for minimum comission rate of 5% is started
	MinCommissionRateEpoch *big.Int `json:"min-commission-rate-epoch,omitempty"`

	// MinCommissionPromoPeriod is the number of epochs when newly elected validators can have 0% commission
	MinCommissionPromoPeriod *big.Int `json:"commission-promo-period,omitempty"`

	// EPoSBound35Epoch is the epoch when the EPoS bound parameter c is changed from 15% to 35%
	EPoSBound35Epoch *big.Int `json:"epos-bound-35-epoch,omitempty"`

	// EIP155 hard fork epoch (include EIP158 too)
	EIP155Epoch *big.Int `json:"eip155-epoch,omitempty"`

	// S3 epoch is the first epoch containing S3 mainnet and all ethereum update up to Constantinople
	S3Epoch *big.Int `json:"s3-epoch,omitempty"`

	// DataCopyFix epoch is the first epoch containing fix for evm datacopy bug.
	DataCopyFixEpoch *big.Int `json:"data-copy-fix-epoch,omitempty"`

	// Istanbul epoch
	IstanbulEpoch *big.Int `json:"istanbul-epoch,omitempty"`

	// ReceiptLogEpoch is the first epoch support receiptlog
	ReceiptLogEpoch *big.Int `json:"receipt-log-epoch,omitempty"`

	// IsSHA3Epoch is the first epoch in supporting SHA3 FIPS-202 standard
	SHA3Epoch *big.Int `json:"sha3-epoch,omitempty"`

	// IsHIP6And8Epoch is the first epoch to support HIP-6 and HIP-8
	HIP6And8Epoch *big.Int `json:"hip6_8-epoch,omitempty"`

	// StakingPrecompileEpoch is the first epoch to support the staking precompiles
	StakingPrecompileEpoch *big.Int `json:"staking-precompile-epoch,omitempty"`
}
