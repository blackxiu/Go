package hmy

type TxsCountArgs struct {
	Address string
	Type    string
}

// TxHistoryArgs is struct to include optional transaction formatting params.
type TxHistoryArgs struct {
	Address   string `json:"address"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	FullTx    bool   `json:"fullTx"`
	TxType    string `json:"txType"`
	Order     string `json:"order"`
}

// BlockArgs is struct to include optional block formatting params.
type BlockArgs struct {
	WithSigners bool     `json:"withSigners"`
	FullTx      bool     `json:"fullTx"`
	Signers     []string `json:"-"`
	InclStaking bool     `json:"inclStaking"`
}
