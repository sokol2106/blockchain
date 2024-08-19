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
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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
	stor := storage.NewPostgresql("host=localhost port=5432 user=pia password=12345678 dbname=yandex sslmode=disable")
	err := stor.Connect()
	suite.Nil(err)
	err = stor.PingContext()
	suite.Nil(err)
	srvBlockchain := service.NewBlockchain(stor)
	srvVerify := service.NewVerification(stor)
	srvBlockchain.StartBlockchainProcessing()
	//srvBlockchain.RunBlockchainDBLoad()

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

	resp, err = http.Get(suite.server.URL + "/api/blockchain/data")
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
	resp, err := http.Post(suite.server.URL+"/api/verify/7e5b4b50-57b4-4f8c-9c54-847c5fa2f4df",
		"text/plain", strings.NewReader("check data in blockchain"))
	suite.Nil(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusCreated, resp.StatusCode)

}

func (suite *ServerTestSuite) TestStatusProcessCheckData() {
	resp, err := http.Get(suite.server.URL + "/api/verify/status/1234567")
	suite.Nil(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestAddBlock() {

	var (
		wg1 sync.WaitGroup
	)

	blocks := make([]model.Block, 200)
	respBlocks := make([]model.Block, len(blocks))

	for i, _ := range blocks {
		in := i
		wg1.Add(1)
		go func(ind int) {
			defer wg1.Done()
			blocks[i] = NewBlock(generateRandomString(10000), strconv.Itoa(ind))
			jsonBlock, err := json.Marshal(blocks[ind])
			suite.Nil(err)
			resp, err := http.Post(suite.server.URL+"/api/blockchain/block", "application/json", strings.NewReader(string(jsonBlock)))
			suite.Nil(err)
			defer resp.Body.Close()
			suite.Equal(http.StatusCreated, resp.StatusCode)
		}(in)
	}

	wg1.Wait()

	time.Sleep(10 * time.Second)

	for i := 0; i < len(blocks); i++ {
		in := i
		wg1.Add(1)
		go func(ind int) {
			defer wg1.Done()
			resp, err := http.Get(suite.server.URL + "/api/blockchain/block")
			suite.Nil(err)
			defer resp.Body.Close()
			suite.Equal(http.StatusOK, resp.StatusCode)

			bodyGetData, err := io.ReadAll(resp.Body)
			suite.Nil(err)
			json.Unmarshal(bodyGetData, &respBlocks[ind])
			suite.Nil(err)

			//prevHead, err := json.Marshal(respBlocks[i].Head)
			//suite.Nil(err)

			//	prevHash := sha256.Sum256(prevHead)

			//	log.Printf("!!!! KEY  %s !!!! Previous  %s !!!! CURREN Previous  %s",
			//		respBlocks[i].Head.Key,
			//		respBlocks[i].Head.Hash,
			//		hex.EncodeToString(prevHash[:]),
			//	)

			//suite.Nil(err)
		}(in)
	}

	wg1.Wait()

	suite.Equal(http.StatusOK, http.StatusOK)
}

func (suite *ServerTestSuite) TestVerification() {
	resp, err := http.Post(suite.server.URL+"/api/verify/12345678",
		"text/plain", strings.NewReader("data for verification"))
	suite.Nil(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusCreated, resp.StatusCode)
	bodyPost, err := io.ReadAll(resp.Body)
	suite.Nil(err)

	type result struct {
		QueueId string `json:"queueId"`
	}

	queueId := result{}
	err = json.Unmarshal(bodyPost, &queueId)
	suite.Nil(err)

	resp3, err := http.Get(suite.server.URL + "/api/verify/status/" + queueId.QueueId)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp3.StatusCode)
	defer resp3.Body.Close()

	bodyGetData, err := io.ReadAll(resp3.Body)
	suite.Nil(err)
	suite.Equal(model.StatusCreated.String(), string(bodyGetData))

	type queueIdStatus struct {
		QueueId string       `json:"queueId"`
		Status  model.Status `json:"status"`
	}

	reqStatus := queueIdStatus{
		QueueId: queueId.QueueId,
		Status:  model.StatusNotFound,
	}

	bodyReq, err := json.Marshal(reqStatus)
	suite.Nil(err)

	resp2, err := http.Post(suite.server.URL+"/api/blockchain/block/verify",
		"application/json", strings.NewReader(string(bodyReq)))

	suite.Nil(err)
	resp2.Body.Close()
	suite.Equal(http.StatusOK, resp2.StatusCode)

	resp4, err := http.Get(suite.server.URL + "/api/verify/status/" + queueId.QueueId)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp4.StatusCode)
	defer resp4.Body.Close()

	bodyGetData, err = io.ReadAll(resp4.Body)
	suite.Nil(err)
	suite.Equal(model.StatusNotFound.String(), string(bodyGetData))

	resp5, err := http.Get(suite.server.URL + "/api/blockchain/block/verify")
	suite.Nil(err)
	suite.Equal(http.StatusNoContent, resp5.StatusCode)
	defer resp4.Body.Close()

}

func NewBlock(msg string, key string) model.Block {
	block := model.Block{}
	// увеличить объём
	block.Data = msg
	block.Head.Noce = "12345678"
	block.Head.Key = key
	hash := sha256.Sum256([]byte(block.Data))
	block.Head.Merkley = hex.EncodeToString(hash[:])
	block.Head.Hash = ""

	return block
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

// Benchmark

func BenchmarkAddBlock(b *testing.B) {
	suite := &ServerTestSuite{}
	suite.SetupSuite()
	defer suite.TearSuiteDownSuite()

	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			block := NewBlock(generateRandomString(10000), "ggggggg")
			jsonBlock, _ := json.Marshal(block)
			_, _ = http.Post(suite.server.URL+"/api/blockchain/block", "application/json", strings.NewReader(string(jsonBlock)))
			_, _ = http.Get(suite.server.URL + "/api/blockchain/block")
		}()
	}

	wg.Wait()
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
