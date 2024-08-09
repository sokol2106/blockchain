package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
)

type StorageBlockchain interface {
	AddBlock(model.Block) error
	GetBlock(string) (*model.Block, error)
	Close() error
}

type KeyData struct {
	Key  string `json:"key"`
	Data string `json:"data,omitempty"`
}

type Blockchain struct {
	queueData  chan KeyData
	lastBlock  model.Block
	queueBlock chan model.Block
	storage    StorageBlockchain
}

func NewBlockchain(stor StorageBlockchain) *Blockchain {
	return &Blockchain{queueData: make(chan KeyData, 100000), queueBlock: make(chan model.Block, 100000), storage: stor}
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

	// обернуть мьютексом
	// копасити
	// запустить отдельный поток для обработки канала с блоками

	block := model.Block{}
	err := json.Unmarshal([]byte(blockStr), &block)
	if err != nil {
		return err
	}

	prevHead, err := json.Marshal(b.lastBlock.Head)
	if err != nil {
		return err
	}

	prevHash := sha256.Sum256(prevHead)
	block.Head.Previous = hex.EncodeToString(prevHash[:])

	select {
	case b.queueBlock <- b.lastBlock:
		b.lastBlock = block
		return nil
	default:
		return errors.New("block queue is full")
	}
}

func (b *Blockchain) ReceiveBlock() (string, error) {
	select {
	case block := <-b.queueBlock:
		jsonData, err := json.Marshal(block)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	default:
		return "", errors.New("block queue is empty")
	}
}

func (b *Blockchain) Close() error {
	close(b.queueData)
	close(b.queueBlock)
	return b.storage.Close()
}
