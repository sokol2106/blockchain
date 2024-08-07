package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
)

type KeyData struct {
	Key  string `json:"key"`
	Data string `json:"data,omitempty"`
}

type KeyBlock struct {
	Key   string      `json:"key"`
	Block model.Block `json:"block,omitempty"`
}

type Blockchain struct {
	queueData  chan KeyData
	queueBlock chan KeyBlock
}

func NewBlockchain() *Blockchain {
	return &Blockchain{queueData: make(chan KeyData, 1000), queueBlock: make(chan KeyBlock, 1000)}
}

func (b *Blockchain) AddData(data string) (string, error) {
	key := uuid.New()
	b.queueData <- KeyData{Key: key.String(), Data: data}
	jsonData, err := json.Marshal(KeyData{Key: key.String()})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (b *Blockchain) ReceiveData() (string, error) {
	keyData := <-b.queueData
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (b *Blockchain) AddBlock(blockStr string) error {
	block := KeyBlock{}
	err := json.Unmarshal([]byte(blockStr), &block)
	if err != nil {
		return err
	}

	b.queueBlock <- block
	return nil
}

func (b *Blockchain) ReceiveBlock() (string, error) {
	block := <-b.queueBlock
	jsonData, err := json.Marshal(block)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (b *Blockchain) Close() {
	close(b.queueData)
}
