package common

import (
	"math/big"
)

// ValidatorRPCEnhanced contains extra information for RPC consumer
type ValidatorRPCEnhanced struct {
	//Wrapper              ValidatorWrapper                 `json:"validator"`
	//Performance          *staking.CurrentEpochPerformance `json:"current-epoch-performance"` // todo: 重建结构体
	//ComputedMetrics      *staking.ValidatorStats          `json:"metrics"`                   // todo: 重建结构体
	//TotalDelegated       *big.Int                         `json:"total-delegation"`
	//CurrentlyInCommittee bool                             `json:"currently-in-committee"`
	//EPoSStatus           string                           `json:"epos-status"`
	//EPoSWinningStake     string                           `json:"epos-winning-stake"` // todo: 重建结构体
	//BootedStatus         string                           `json:"booted-status"`
	//ActiveStatus         string                           `json:"active-status"`
	//Lifetime             *staking.AccumulatedOverLifetime `json:"lifetime"` // todo: 重建结构体
}

// ValidatorWrapper contains validator,
// its delegation information
type ValidatorWrapper struct {
	//staking.Validator // todo: Validator结构体重建
	//Delegations       Delegations
	////
	//Counters counters `json:"-"`
	//// All the rewarded accumulated so far
	//BlockReward *big.Int `json:"-"`
}

//
//type Validator struct {
//	Address     string   `json:"address"`
//	SlotPubKeys []string `json:"bls-public-keys"`
//	// The number of the last epoch this validator is
//	// selected in committee (0 means never selected)
//	LastEpochInCommittee *big.Int `json:"last-epoch-in-committee"`
//	// validator's self declared minimum self delegation
//	MinSelfDelegation *big.Int `json:"min-self-delegation"`
//	// maximum total delegation allowed
//	MaxTotalDelegation *big.Int `json:"max-total-delegation"`
//	// Is the validator active in participating
//	// committee selection process or not
//	Status effective.Eligibility `json:"-"`
//	// commission parameters
//	Commission
//	// description for the validator
//	Description
//	// CreationHeight is the height of creation
//	CreationHeight *big.Int `json:"creation-height"`
//}

type counters struct {
	// The number of blocks the validator
	// should've signed when in active mode (selected in committee)
	NumBlocksToSign *big.Int `json:"to-sign",rlp:"nil"`
	// The number of blocks the validator actually signed
	NumBlocksSigned *big.Int `json:"signed",rlp:"nil"`
}

type Delegations []Delegation

type Delegation struct {
	DelegatorAddress string        `json:"delegator-address"`
	Amount           *big.Int      `json:"amount"`
	Reward           *big.Int      `json:"reward"`
	Undelegations    Undelegations `json:"undelegations"`
}

type Undelegations []Undelegation

type Undelegation struct {
	Amount *big.Int `json:"amount"`
	Epoch  *big.Int `json:"epoch"`
}

// UtilityMetric ..
type UtilityMetric struct {
	AccumulatorSnapshot     *big.Int `json:"AccumulatorSnapshot"`
	CurrentStakedPercentage string   `json:"CurrentStakedPercentage"`
	Deviation               string   `json:"Deviation"`
	Adjustment              string   `json:"Adjustment"`
}

// CompletedEPoSRound ..
type CompletedEPoSRound struct {
	MedianStake         string            `json:"epos-median-stake"`
	MaximumExternalSlot int               `json:"max-external-slots"`
	AuctionWinners      []SlotPurchase    `json:"epos-slot-winners"`
	AuctionCandidates   []*CandidateOrder `json:"epos-slot-candidates"`
}

// SlotPurchase ..
type SlotPurchase struct {
	Addr      string `json:"slot-owner"`
	Key       string `json:"bls-public-key"`
	RawStake  string `json:"raw-stake"`
	EPoSStake string `json:"eposed-stake"`
}

type CandidateOrder struct {
	*SlotOrder
	StakePerKey *big.Int `json:"stake-per-key"`
	Validator   string   `json:"validator"`
}

// SlotOrder ..
type SlotOrder struct {
	Stake       *big.Int `json:"stake"`
	SpreadAmong []string `json:"keys-at-auction"`
	Percentage  string   `json:"percentage-of-total-auction-stake"`
}
