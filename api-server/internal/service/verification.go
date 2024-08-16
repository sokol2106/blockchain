package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/customerrors"
	"github.com/ivan/blockchain/api-server/internal/model"
	"log"
	"sync"
)

type Verification struct {
	rawDataQueue       chan *model.VerificationData // Очередь необработанных данных для проверки
	dataBlockQueue     chan *model.VerificationData // Очередь данных, привязанных к блокам для проверки
	activeProcessesMap sync.Map                     // Карта для отслеживания текущих процессов проверки
	storage            StorageBlockchain            // Хранилище блокчейна
}

func NewVerification(stor StorageBlockchain) *Verification {
	return &Verification{
		rawDataQueue:   make(chan *model.VerificationData, 100000),
		dataBlockQueue: make(chan *model.VerificationData, 100000),
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

	v.rawDataQueue <- objVrf
	v.activeProcessesMap.Store(queueId, objVrf)
	return queueId, nil
}

func (v *Verification) StartBlockSearchProcess() {
	rwmu := sync.RWMutex{}
	go func() {
		for objRawVrf := range v.rawDataQueue {
			rwmu.Lock()
			objRawVrf.Status = model.StatusProcessing
			rwmu.Unlock()

			foundBlock, err := v.storage.SelectBlock(context.Background(), objRawVrf.Key)
			if err != nil {
				log.Println("error getting block ", err)
				break
			}

			rwmu.Lock()
			objRawVrf.Block = *foundBlock
			v.dataBlockQueue <- objRawVrf
			rwmu.Unlock()
		}
	}()
}

func (v *Verification) GetProcessStatus(queueID string) model.Status {
	value, exist := v.activeProcessesMap.Load(queueID)
	if exist {
		modelVer := value.(*model.VerificationData)
		return modelVer.Status
	} else {
		return model.StatusNotFound
	}

}

func (v *Verification) RetrieveProcessedData() (string, error) {
	select {
	case data, ok := <-v.dataBlockQueue:
		if !ok {
			return "", customerrors.ErrNoDataToVerification
		}
		jsonData, err := json.Marshal(*data)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	default:
		return "", customerrors.ErrNoDataToVerification
	}
}

func (v *Verification) UpdateStatus(queueId string, status model.Status) {
	objVref, exist := v.activeProcessesMap.Load(queueId)
	if exist {
		modelVer := objVref.(*model.VerificationData)
		modelVer.Status = status
	}
}

func (v *Verification) Close() error {
	close(v.rawDataQueue)
	close(v.dataBlockQueue)
	return v.storage.Disconnect()
}
