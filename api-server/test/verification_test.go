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

		err = stor.AddBlock(block)
		require.NoError(t, err)

		queueId, err := vrf.AddData(key, "hello word")
		require.NoError(t, err)

		status := vrf.StatusProcess(queueId)
		require.Equal(t, model.StatusCreated, status)

		vrf.RunProcessSearchBlock()

		time.Sleep(1 * time.Second)

		status = vrf.StatusProcess(queueId)
		require.Equal(t, model.StatusProcessing, status)

		res, err := vrf.ReceiveDataHandelr()
		require.NoError(t, err)

		blockVrf := model.VerificationData{}
		err = json.Unmarshal([]byte(res), &blockVrf)
		require.NoError(t, err)
		require.Equal(t, queueId, blockVrf.QueueId)

		time.Sleep(3 * time.Second)

		// Проверка что поток обработки данных не остановился

		queueId, err = vrf.AddData(key, "hello word 2")
		require.NoError(t, err)

		status = vrf.StatusProcess(queueId)
		require.Equal(t, model.StatusProcessing, status)

		res, err = vrf.ReceiveDataHandelr()
		require.NoError(t, err)

		blockVrf = model.VerificationData{}
		err = json.Unmarshal([]byte(res), &blockVrf)
		require.NoError(t, err)
		require.Equal(t, queueId, blockVrf.QueueId)

	})

}
