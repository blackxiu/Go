package common

import (
	"encoding/json"
	"math/big"
)

type CrossLink struct {
	Hash        string   `json:"hash"`
	BlockNumber *big.Int `json:"block-number"`
	ViewID      *big.Int `json:"view-id"`
	Signature   string   `json:"signature"`
	Bitmap      string   `json:"signature-bitmap"`
	ShardID     uint32   `json:"shard-id"`
	EpochNumber *big.Int `json:"epoch-number"`
}

type CrossLinks []CrossLink

// MarshalJSON ..
func (cl *CrossLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Hash        string   `json:"hash"`
		BlockNumber *big.Int `json:"block-number"`
		ViewID      *big.Int `json:"view-id"`
		Signature   string   `json:"signature"`
		Bitmap      string   `json:"signature-bitmap"`
		ShardID     uint32   `json:"shard-id"`
		EpochNumber *big.Int `json:"epoch-number"`
	}{
		cl.Hash,
		cl.BlockNumber,
		cl.ViewID,
		cl.Signature[:],
		cl.Bitmap,
		cl.ShardID,
		cl.EpochNumber,
	})
}
