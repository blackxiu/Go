package common

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
	//
	PublicKeySizeInBytes = 48
	//
	BLSSignatureSizeInBytes = 96
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

//
type Address [AddressLength]byte

//
type Directive byte

// SerializedPublicKey defines the serialized bls public key
type SerializedPublicKey [PublicKeySizeInBytes]byte

// SerializedSignature defines the bls signature
type SerializedSignature [BLSSignatureSizeInBytes]byte

// Hex converts a hash to a hex string.
//func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

type C struct {
	TotalKnownPeers int `json:"total-known-peers"`
	Connected       int `json:"connected"`
	NotConnected    int `json:"not-connected"`
}

// ConsensusInternal captures consensus internal data
type ConsensusInternal struct {
	ViewID        uint64 `json:"viewId"`
	ViewChangeID  uint64 `json:"viewChangeId"`
	Mode          string `json:"mode"`
	Phase         string `json:"phase"`
	BlockNum      uint64 `json:"blocknum"`
	ConsensusTime int64  `json:"finality"`
}

// NodeMetadata captures select metadata of the RPC answering node
type NodeMetadata struct {
	BLSPublicKey    []string          `json:"blskey"`
	Version         string            `json:"version"`
	NetworkType     string            `json:"network"`
	ChainConfig     ChainConfig       `json:"chain-config"`
	IsLeader        bool              `json:"is-leader"`
	ShardID         uint32            `json:"shard-id"`
	CurrentBlockNum uint64            `json:"current-block-number"`
	CurrentEpoch    uint64            `json:"current-epoch"`
	BlocksPerEpoch  *uint64           `json:"blocks-per-epoch,omitempty"`
	Role            string            `json:"role"`
	DNSZone         string            `json:"dns-zone"`
	Archival        bool              `json:"is-archival"`
	IsBackup        bool              `json:"is-backup"`
	NodeBootTime    int64             `json:"node-unix-start-time"`
	PeerID          string            `json:"peerid"`
	Consensus       ConsensusInternal `json:"consensus"`
	C               C                 `json:"p2p-connectivity"`
	SyncPeers       map[string]int    `json:"sync-peers,omitempty"`
}

// Transition is for the interface of "hmyv2_getSuperCommittees"
type Transition struct {
	Previous Registry `json:"previous"`
	Current  Registry `json:"current"`
}

// Registry ..
type Registry struct {
	Deciders      map[string]Decider `json:"quorum-deciders"`
	ExternalCount int                `json:"external-slot-count"`
	MedianStake   string             `json:"epos-median-stake"`
	Epoch         int                `json:"epoch"`
}

type Decider struct {
	Policy              string        `json:"policy"`
	Count               int           `json:"count"`
	Externals           int           `json:"external-validator-slot-count"`
	Participants        []Participant `json:"committee-members"`
	HmyVotingPower      string        `json:"hmy-voting-power"`
	StakedVotingPower   string        `json:"staked-voting-power"`
	TotalRawStake       string        `json:"total-raw-stake"`
	TotalEffectiveStake string        `json:"total-effective-stake"`
}

type Participant struct {
	IsHarmony      bool   `json:"is-harmony-slot"`
	EarningAccount string `json:"earning-account"`
	Identity       string `json:"bls-public-key"`
	RawPercent     string `json:"voting-power-unnormalized"`
	VotingPower    string `json:"voting-power-%"`
	EffectiveStake string `json:"effective-stake,omitempty"`
	RawStake       string `json:"raw-stake,omitempty"`
}
