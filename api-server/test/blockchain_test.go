package test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/ivan/blockchain/api-server/internal/model"
	"github.com/ivan/blockchain/api-server/internal/service"
	"github.com/ivan/blockchain/api-server/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRunBlockchain(t *testing.T) {
	stor := storage.NewPostgresql("host=localhost port=5432 user=pia password=12345678 dbname=yandex sslmode=disable")
	count := 100

	t.Run("Test Run blockchain", func(t *testing.T) {
		err := stor.Connect()
		require.NoError(t, err)
		err = stor.PingContext()
		require.NoError(t, err)

		blch := service.NewBlockchain(stor)
		defer blch.Close()
		blch.StartBlockchainProcessing()

		for range count {
			go func() {
				rawBlock := NewBlock(generateRandomString(1000000), generateRandomString(20))
				tBlock, err := json.Marshal(rawBlock)
				require.NoError(t, err)
				err = blch.AddNewBlock(string(tBlock))
				require.NoError(t, err)
			}()
		}

		time.Sleep(5 * time.Second)

		prev, err := blch.ReceiveBlock()
		require.NoError(t, err)
		prevBlock := model.Block{}
		json.Unmarshal([]byte(prev), &prevBlock)

		for range count - 1 {
			block, err := blch.ReceiveBlock()
			require.NoError(t, err)
			currentBlock := model.Block{}
			json.Unmarshal([]byte(block), &currentBlock)

			resultHash := currentBlock.Head.Hash
			currentBlock.Head.Hash = ""

			rawHead, err := json.Marshal(currentBlock.Head)
			require.NoError(t, err)
			newCurrentBlockHash := sha256.Sum256(rawHead)

			newDataHash := hex.EncodeToString(newCurrentBlockHash[:]) + prevBlock.Head.Hash
			prevHash := sha256.Sum256([]byte(newDataHash))

			assert.Equal(t, resultHash, hex.EncodeToString(prevHash[:]))
			//log.Printf(resultHash + " ----- " + hex.EncodeToString(prevHash[:]))

			currentBlock.Head.Hash = resultHash
			prevBlock = currentBlock
		}

	})

}
