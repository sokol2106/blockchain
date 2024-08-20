package test

import (
	"github.com/ivan/blockchain/block-client/internal/model"
	"github.com/ivan/blockchain/block-client/internal/service"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRunMining(t *testing.T) {
	srv := service.NewBlockMiner("", "0000")
	data := model.MiningData{
		Data: generateRandomString(1000000),
		Key:  "666666",
	}

	t.Run("Test Run Mining", func(t *testing.T) {
		srv.SetData(data)
		srv.MineData()
		block := srv.GetBlock()
		assert.Equal(t, data.Data, block.Data)
	})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
