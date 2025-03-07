package model

type BlockHeader struct {
	Hash    string `json:"hash"`
	Merkley string `json:"merkley"`
	Key     string `json:"key"`
}

type Block struct {
	Head BlockHeader `json:"status"`
	Data string      `json:"queueID"`
}

type CheckDataBlock struct {
	Block Block  `json:"block"`
	Data  string `json:"data"`
}
