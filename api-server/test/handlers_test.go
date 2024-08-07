package test

import (
	"encoding/json"
	"github.com/ivan/blockchain/api-server/internal/handlers"
	"github.com/ivan/blockchain/api-server/internal/service"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
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
	srvBlockchain := service.NewBlockchain()
	srvVerify := service.NewVerification()

	suite.server = httptest.NewServer(handlers.Router(handlers.NewHandlers(srvBlockchain, srvVerify)))
}

func (suite *ServerTestSuite) TearSuiteDownSuite() {
	suite.server.Close()
}

func (suite *ServerTestSuite) TestAddDataBlockchain() {
	resp, err := http.Post(suite.server.URL+"/api/setdata",
		"text/plain", strings.NewReader("NEW data in blockchain"))
	suite.Nil(err)
	defer resp.Body.Close()
	bodySetData, err := io.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(suite.server.URL + "/api/block")
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
	resp, err := http.Post(suite.server.URL+"/api/checkdata/7e5b4b50-57b4-4f8c-9c54-847c5fa2f4df",
		"text/plain", strings.NewReader("check data in blockchain"))
	suite.Nil(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestStatusProcessCheckData() {
	resp, err := http.Get(suite.server.URL + "/api/checkdata/1234567")
	suite.Nil(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
