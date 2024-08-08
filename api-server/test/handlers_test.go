package test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/ivan/blockchain/api-server/internal/handlers"
	"github.com/ivan/blockchain/api-server/internal/model"
	"github.com/ivan/blockchain/api-server/internal/service"
	"github.com/ivan/blockchain/api-server/internal/storage"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
)

import (
	"net/http"
	"net/http/httptest"
)

type ServerTestSuite struct {
	suite.Suite
	server *httptest.Server
	cookie *http.Cookie
}

func (suite *ServerTestSuite) SetupSuite() {
	stor := storage.NewPostgresql("")
	srvBlockchain := service.NewBlockchain(stor)
	srvVerify := service.NewVerification()

	suite.server = httptest.NewServer(handlers.Router(handlers.NewHandlers(srvBlockchain, srvVerify)))
}

func (suite *ServerTestSuite) TearSuiteDownSuite() {
	suite.server.Close()
}

func (suite *ServerTestSuite) TestAddDataBlockchain() {
	resp, err := http.Post(suite.server.URL+"/api/data",
		"text/plain", strings.NewReader("NEW data in blockchain"))
	suite.Nil(err)
	defer resp.Body.Close()
	bodySetData, err := io.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(suite.server.URL + "/api/data")
	suite.Nil(err)
	defer resp.Body.Close()
	bodyGetData, err := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var (
		jsSetData, jsGetData service.KeyData
	)

	err = json.Unmarshal(bodySetData, &jsSetData)
	suite.Nil(err)

	err = json.Unmarshal(bodyGetData, &jsGetData)
	suite.Nil(err)

	suite.Equal(jsGetData.Key, jsSetData.Key)
	suite.Equal("NEW data in blockchain", jsGetData.Data)
}

func (suite *ServerTestSuite) TestCheckData() {
	resp, err := http.Post(suite.server.URL+"/api/check/7e5b4b50-57b4-4f8c-9c54-847c5fa2f4df",
		"text/plain", strings.NewReader("check data in blockchain"))
	suite.Nil(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestStatusProcessCheckData() {
	resp, err := http.Get(suite.server.URL + "/api/check/1234567")
	suite.Nil(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestAddBlock() {

	var (
		wg1 sync.WaitGroup
	)

	blocks := make([]*model.Block, 100)
	respBlocks := make([]*model.Block, len(blocks))

	for i, _ := range blocks {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			blocks[i] = NewBlock("test "+strconv.Itoa(i), strconv.Itoa(i))
			jsonBlock, err := json.Marshal(blocks[i])
			suite.Nil(err)
			resp, err := http.Post(suite.server.URL+"/api/block", "application/json", strings.NewReader(string(jsonBlock)))
			suite.Nil(err)
			defer resp.Body.Close()
			suite.Equal(http.StatusCreated, resp.StatusCode)
		}()
	}

	wg1.Wait()

	for i := 0; i < len(blocks); i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			resp, err := http.Get(suite.server.URL + "/api/block")
			suite.Nil(err)
			defer resp.Body.Close()
			suite.Equal(http.StatusOK, resp.StatusCode)

			bodyGetData, err := io.ReadAll(resp.Body)
			suite.Nil(err)
			err = json.Unmarshal(bodyGetData, &respBlocks[i])
			suite.Nil(err)

			prevHead, err := json.Marshal(respBlocks[i].Head)
			suite.Nil(err)

			prevHash := sha256.Sum256(prevHead)

			log.Printf("!!!! KEY  %s !!!! Previous  %s !!!! CURREN Previous  %s",
				respBlocks[i].Head.Key,
				respBlocks[i].Head.Previous,
				hex.EncodeToString(prevHash[:]),
			)

			suite.Nil(err)
		}()
	}

	wg1.Wait()

	suite.Equal(http.StatusOK, http.StatusOK)
}

func NewBlock(msg string, key string) *model.Block {
	block := model.Block{}
	// увеличить объём
	block.Data = msg
	block.Head.Noce = "12345678"
	block.Head.Key = key
	hash := sha256.Sum256([]byte(block.Data))
	block.Head.Merkley = hex.EncodeToString(hash[:])

	return &block
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
