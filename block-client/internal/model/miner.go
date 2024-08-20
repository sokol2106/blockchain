package model

type Header struct {
	Hash    string `json:"hash"`
	Merkley string `json:"merkley"`
	Key     string `json:"key"`
	Nonce   string `json:"nonce"`
}

type Block struct {
	Head Header `json:"head"`
	Data string `json:"data"`
}

type MiningData struct {
	Key  string `json:"key"`
	Data string `json:"data,omitempty"`
}
