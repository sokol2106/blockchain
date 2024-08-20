package test

import (
	"github.com/ivan/blockchain/block-client/internal/model"
	"github.com/ivan/blockchain/block-client/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunVerify(t *testing.T) {

	srv := service.NewBlockMiner("", "0000")
	data := model.MiningData{
		Data: generateRandomString(1000000),
		Key:  "666666",
	}

	srv.SetData(data)
	srv.MineData()
	block := srv.GetBlock()

	srvVrf := service.NewBlockVerification("", "0000")
	blockVrf := model.VerificationBlock{
		QueueId: "1",
		Status:  model.StatusProcessing,
		Key:     block.Head.Key,
		Data:    "AAAAAAAA",
		Block:   block,
	}

	t.Run("Test Run Verify", func(t *testing.T) {
		srvVrf.SetBlock(blockVrf)
		srvVrf.VerifyData()

		resBlock := srvVrf.GetBlock()
		assert.Equal(t, model.StatusFailedAuthenticityCheck, resBlock.Status)

		blockVrf.Data = block.Data

		srvVrf.SetBlock(blockVrf)
		srvVrf.VerifyData()

		resBlock = srvVrf.GetBlock()
		assert.Equal(t, model.StatusMatched, resBlock.Status)
	})

}
