package service

import (
	"github.com/ivan/blockchain/api-server/internal/model"
)

type Verification struct {
	queueData      chan model.DataKey
	queueDataBlock chan model.CheckDataBlock
	storage        StorageBlockchain
}

func NewVerification() *Verification {
	return &Verification{
		queueData:      make(chan model.DataKey, 100000),
		queueDataBlock: make(chan model.CheckDataBlock, 100000),
	}
}

func (v *Verification) AddData(key string, data string) {
	v.queueData <- model.DataKey{Key: key, Data: data}
}

func (v *Verification) RunProcessSearchBlock() {
	go func() {
		for {
			// Получаем не обработанный блок
			data, ok := <-v.queueData
			if !ok {
				break
			}
			data.Data = ""

		}
	}()
}

func (v *Verification) AddDataBlock(msg model.CheckDataBlock) {

}
