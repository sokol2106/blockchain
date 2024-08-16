package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/ivan/blockchain/api-server/internal/model"
	"github.com/ivan/blockchain/api-server/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnectPostgresql(t *testing.T) {

	t.Run("Test connect postgresql", func(t *testing.T) {
		stor := storage.NewPostgresql("host=localhost port=5432 user=pia password=12345678 dbname=yandex sslmode=disable")
		err := stor.Connect()
		require.NoError(t, err)
		defer stor.Disconnect()

		//err = stor.Migrations("file://../migrations/postgresql")
		//require.NoError(t, err)

		key := uuid.New().String()

		block := model.Block{Data: generateRandomString(10000), Head: model.BlockHeader{
			Hash:    "11111111",
			Key:     key,
			Noce:    "000",
			Merkley: "99999999",
		}}

		err = stor.InsertBlock(block)
		require.NoError(t, err)

		block2, err := stor.SelectBlock(context.Background(), key)
		require.NoError(t, err)
		require.Equal(t, block.Data, block2.Data)
	})
}
