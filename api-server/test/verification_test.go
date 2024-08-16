package test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
	"github.com/ivan/blockchain/api-server/internal/service"
	"github.com/ivan/blockchain/api-server/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVerification(t *testing.T) {

	//pgsql := storage.NewPostgresql("")
	stor := storage.NewPostgresql("host=localhost port=5432 user=pia password=12345678 dbname=yandex sslmode=disable")
	err := stor.Connect()
	require.NoError(t, err)

	//err = stor.Migrations("file://../migrations/postgresql")
	//require.NoError(t, err)

	vrf := service.NewVerification(stor)
	defer vrf.Close()

	t.Run("Test connect postgresql", func(t *testing.T) {

		key := uuid.New().String()
		block := model.Block{Data: "Hello word", Head: model.BlockHeader{
			Hash:    "11111111",
			Key:     key,
			Noce:    "000",
			Merkley: "99999999",
		}}

		err = stor.InsertBlock(block)
		require.NoError(t, err)

		queueId, err := vrf.AddData(key, "hello word")
		require.NoError(t, err)

		status := vrf.GetProcessStatus(queueId)
		require.Equal(t, model.StatusCreated, status)

		vrf.StartBlockSearchProcess()

		time.Sleep(1 * time.Second)

		status = vrf.GetProcessStatus(queueId)
		require.Equal(t, model.StatusProcessing, status)

		res, err := vrf.RetrieveProcessedData()
		require.NoError(t, err)

		blockVrf := model.VerificationData{}
		err = json.Unmarshal([]byte(res), &blockVrf)
		require.NoError(t, err)
		require.Equal(t, queueId, blockVrf.QueueId)

		time.Sleep(3 * time.Second)

		// Проверка что поток обработки данных не остановился

		queueId2, err := vrf.AddData(key, "hello word 2")
		require.NoError(t, err)

		endTime := time.Now().Add(10 * time.Second)

		for {
			if time.Now().After(endTime) {
				break
			}

			res, err = vrf.RetrieveProcessedData()
			if err == nil {
				break
			}
		}

		require.NoError(t, err)
		blockVrf = model.VerificationData{}
		err = json.Unmarshal([]byte(res), &blockVrf)
		require.NoError(t, err)
		require.Equal(t, queueId2, blockVrf.QueueId)

		status = vrf.GetProcessStatus(queueId2)
		require.Equal(t, model.StatusProcessing, status)
	})

}
