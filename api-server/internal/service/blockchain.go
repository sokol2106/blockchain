//go:generate mockgen -destination=mock_storageblockchain.go -package=service . Service
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
	"log"
)

type StorageBlockchain interface {
	InsertBlock(model.Block) error
	SelectBlock(context.Context, string) (*model.Block, error)
	SelectLastBlock(context.Context) (*model.Block, error)
	Disconnect() error
}

type KeyData struct {
	Key  string `json:"key"`
	Data string `json:"data,omitempty"`
}

type Blockchain struct {
	dataQueue       chan KeyData      // Очередь для поступающих данных
	rawBlockQueue   chan model.Block  // Очередь для необработанных блоков (готовящихся к добавлению)
	blockchainQueue chan model.Block  // Очередь для подтвержденных блоков в блокчейне
	previousBlock   model.Block       // Ссылка на предыдущий блок в цепи
	storage         StorageBlockchain // Хранилище блокчейна
}

func NewBlockchain(stor StorageBlockchain) *Blockchain {
	prevBlock := &model.Block{
		Data: "start block",
		Head: model.BlockHeader{
			Hash:    "6c818bd1063cb91ebd803fc894c01a49fd5fecaa4a86693c7c02b8296b9d45ee",
			Merkley: "4rfvbgt56yhn",
			Key:     "12345678",
			Noce:    "12345678",
		},
	}

	if stor != nil {
		prevBlock2, err := stor.SelectLastBlock(context.Background())
		if err != nil {
			log.Printf("error load start block: %s", err)
		} else {
			prevBlock = prevBlock2
		}

	}

	return &Blockchain{
		dataQueue:       make(chan KeyData, 100000),
		rawBlockQueue:   make(chan model.Block, 100000),
		blockchainQueue: make(chan model.Block, 100000),
		previousBlock:   *prevBlock,
		storage:         stor,
	}
}

func (b *Blockchain) StoreData(data string) (string, error) {
	key := uuid.New()
	b.dataQueue <- KeyData{Key: key.String(), Data: data}
	jsonData, err := json.Marshal(KeyData{Key: key.String()})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (b *Blockchain) ReceiveData() (string, error) {
	keyData := <-b.dataQueue
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (b *Blockchain) AddNewBlock(blockStr string) error {
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
	case block := <-b.blockchainQueue:
		jsonData, err := json.Marshal(block)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	default:
		return "", errors.New("block queue is empty")
	}
}

func (b *Blockchain) StartBlockchainProcessing() {
	go func() {
		for rawBlock := range b.rawBlockQueue {
			rawHead, err := json.Marshal(rawBlock.Head)
			if err != nil {
				log.Printf("error marshal (Run Blockchain): %s", err)
				break
			}

			currentBlockHash := sha256.Sum256(rawHead)

			newHash := hex.EncodeToString(currentBlockHash[:]) + b.previousBlock.Head.Hash
			prevHash := sha256.Sum256([]byte(newHash))
			rawBlock.Head.Hash = hex.EncodeToString(prevHash[:])

			//log.Printf("ProcessBlockChain %s : DATA : %s", b.previousBlock.Head.Hash, b.previousBlock.Data)
			b.blockchainQueue <- b.previousBlock
			b.previousBlock = rawBlock
		}
	}()
}

func (b *Blockchain) StartDBSync() {
	go func() {
		for rawBlock := range b.blockchainQueue {
			//log.Printf("BlockchainDBLoad %s : DATA : %s", rawBlock.Head.Hash, rawBlock.Data)
			err := b.storage.InsertBlock(rawBlock)
			if err != nil {
				log.Printf("error Load DB : %s", err)
			}
		}
	}()
}

func (b *Blockchain) Close() error {
	close(b.dataQueue)
	close(b.rawBlockQueue)
	close(b.blockchainQueue)
	return b.storage.Disconnect()
}
