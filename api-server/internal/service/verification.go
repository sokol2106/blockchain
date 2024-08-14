package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
	"log"
	"sync"
)

type Verification struct {
	queueRawData      chan *model.VerificationData
	queueDataBlock    chan *model.VerificationData
	processingDataMap sync.Map
	storage           StorageBlockchain
}

func NewVerification(stor StorageBlockchain) *Verification {
	return &Verification{
		queueRawData:   make(chan *model.VerificationData, 100000),
		queueDataBlock: make(chan *model.VerificationData, 100000),
		storage:        stor,
	}
}

func (v *Verification) AddData(key string, data string) (string, error) {
	queueId := uuid.New().String()
	objVrf := &model.VerificationData{
		Key:     key,
		Data:    data,
		QueueId: queueId,
		Status:  model.StatusCreated,
	}

	v.queueRawData <- objVrf
	v.processingDataMap.Store(queueId, objVrf)
	return queueId, nil
}

func (v *Verification) RunProcessSearchBlock() {
	rwmu := sync.RWMutex{}
	go func() {
		for {
			objRawVrf, ok := <-v.queueRawData
			if !ok {
				break
			}

			rwmu.Lock()
			objRawVrf.Status = model.StatusProcessing
			rwmu.Unlock()

			foundBlock, err := v.storage.GetBlock(context.Background(), objRawVrf.Key)
			if err != nil {
				log.Println("error getting block ", err)
				break
			}

			rwmu.Lock()
			objRawVrf.Block = *foundBlock
			v.queueDataBlock <- objRawVrf
			rwmu.Unlock()
		}
	}()
}

func (v *Verification) StatusProcess(queueID string) model.Status {
	value, exist := v.processingDataMap.Load(queueID)
	if exist {
		modelVer := value.(*model.VerificationData)
		return modelVer.Status
	} else {
		return model.StatusNotFound
	}

}

func (v *Verification) ReceiveDataHandelr() (string, error) {
	data := <-v.queueDataBlock
	jsonData, err := json.Marshal(*data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (v *Verification) SetStatus(queueId string, status model.Status) {
	objVref, exist := v.processingDataMap.Load(queueId)
	if exist {
		modelVer := objVref.(*model.VerificationData)
		modelVer.Status = status
	}
}

func (v *Verification) Close() error {
	close(v.queueRawData)
	close(v.queueDataBlock)
	return v.storage.Close()
}
