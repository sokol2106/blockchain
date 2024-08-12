//go:generate mockgen -destination=mock_storageblockchain.go -package=service . Service
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
	"log"
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
	queueData       chan KeyData
	rawBlockQueue   chan model.Block
	blockChainQueue chan model.Block
	prevBlock       model.Block
	storage         StorageBlockchain
}

func NewBlockchain(stor StorageBlockchain) *Blockchain {

	return &Blockchain{
		queueData:       make(chan KeyData, 100000),
		rawBlockQueue:   make(chan model.Block, 100000),
		blockChainQueue: make(chan model.Block, 100000),
		prevBlock: model.Block{
			Data: "load block",
			Head: model.BlockHeader{
				Hash:    "6c818bd1063cb91ebd803fc894c01a49fd5fecaa4a86693c7c02b8296b9d45ee",
				Merkley: "4rfvbgt56yhn",
				Key:     "12345678",
				Noce:    "12345678",
			},
		},
		storage: stor,
	}
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
	block := model.Block{}
	err := json.Unmarshal([]byte(blockStr), &block)
	if err != nil {
		return err
	}

	select {
	case b.rawBlockQueue <- block:
		return nil
	default:
		return errors.New("block queue is full")
	}
}

func (b *Blockchain) ReceiveBlock() (string, error) {
	select {
	case block := <-b.blockChainQueue:
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
	close(b.rawBlockQueue)
	close(b.blockChainQueue)
	return b.storage.Close()
}

func (b *Blockchain) Run() {
	go func() {
		for {
			// Получаем не обработанный блок
			rawBlock, ok := <-b.rawBlockQueue
			if !ok {
				break
			}

			rawHead, err := json.Marshal(rawBlock.Head)
			if err != nil {
				log.Printf("error marshal (Run Blockchain): %s", err)
				break
			}

			currentBlockHash := sha256.Sum256(rawHead)

			newHash := hex.EncodeToString(currentBlockHash[:]) + b.prevBlock.Head.Hash
			prevHash := sha256.Sum256([]byte(newHash))
			rawBlock.Head.Hash = hex.EncodeToString(prevHash[:])

			b.blockChainQueue <- b.prevBlock
			b.prevBlock = rawBlock
		}
	}()
}
