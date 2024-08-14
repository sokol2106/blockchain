package model

type BlockHeader struct {
	Hash    string `json:"hash"`
	Merkley string `json:"merkley"`
	Key     string `json:"key"`
	Noce    string `json:"noce"`
}

type Block struct {
	Head BlockHeader `json:"head"`
	Data string      `json:"data"`
}
