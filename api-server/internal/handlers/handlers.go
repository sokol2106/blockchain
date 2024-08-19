package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/ivan/blockchain/api-server/internal/customerrors"
	"github.com/ivan/blockchain/api-server/internal/middleware"
	"github.com/ivan/blockchain/api-server/internal/model"
	"github.com/ivan/blockchain/api-server/internal/service"
	"io"
	"net/http"
)

type Handlers struct {
	srvBlockchain *service.Blockchain
	srvVerify     *service.Verification
}

func NewHandlers(blch *service.Blockchain, vrf *service.Verification) *Handlers {
	return &Handlers{
		srvBlockchain: blch,
		srvVerify:     vrf,
	}
}

func (h *Handlers) handlerError(err error) int {
	statusCode := http.StatusBadRequest
	if errors.Is(err, customerrors.ErrNoDataToVerification) {
		statusCode = http.StatusNoContent
	}

	return statusCode
}

func (h *Handlers) CreateBlockchainData(res http.ResponseWriter, req *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	result, err := h.srvBlockchain.StoreData(string(body))

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(handlerStatus)
	res.Write([]byte(result))
}

func (h *Handlers) GetDataForBlockCreation(res http.ResponseWriter, req *http.Request) {
	data, err := h.srvBlockchain.ReceiveData()
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(data))
}

func (h *Handlers) VerifyData(res http.ResponseWriter, req *http.Request) {
	_, cancel := context.WithCancel(req.Context())
	defer cancel()
	key := chi.URLParam(req, "key")

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	type result struct {
		QueueId string `json:"queueId"`
	}

	strResult := result{}
	strResult.QueueId, err = h.srvVerify.AddData(key, string(body))
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	bodyResult, err := json.Marshal(strResult)
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(bodyResult)
}

func (h *Handlers) GetVerificationStatus(res http.ResponseWriter, req *http.Request) {
	_, cancel := context.WithCancel(req.Context())
	defer cancel()

	queueID := chi.URLParam(req, "queue_id")
	status := h.srvVerify.GetProcessStatus(queueID)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(status.String()))
}

func (h *Handlers) AddBlockchainBlock(res http.ResponseWriter, req *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	err = h.srvBlockchain.AddNewBlock(string(body))

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	res.WriteHeader(handlerStatus)
}

func (h *Handlers) GetBlockchainBlock(res http.ResponseWriter, req *http.Request) {
	block, err := h.srvBlockchain.ReceiveBlock()
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(block))
}

func (h *Handlers) GetBlockVerificationData(res http.ResponseWriter, req *http.Request) {
	blockVrf, err := h.srvVerify.RetrieveProcessedData()
	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	if err != nil {
		res.WriteHeader(h.handlerError(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(blockVrf))
}

func (h *Handlers) UpdateVerificationStatus(res http.ResponseWriter, req *http.Request) {
	type queueIdStatus struct {
		QueueId string       `json:"queueId"`
		Status  model.Status `json:"status"`
	}

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	strResult := queueIdStatus{}
	err = json.Unmarshal(body, &strResult)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	h.srvVerify.UpdateStatus(strResult.QueueId, strResult.Status)
	res.WriteHeader(http.StatusOK)
}

func Router(handler *Handlers) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.СompressionResponseRequest)
	router.Use(middleware.LoggingResponseRequest)

	// router

	// приходит внешне
	router.Post("/api/data", http.HandlerFunc(handler.CreateBlockchainData))                     // новые данные для сохранения +
	router.Post("/api/verify/{key}", http.HandlerFunc(handler.VerifyData))                       // ключ + данные для проверки подлинности +
	router.Get("/api/verify/status/{queue_id}", http.HandlerFunc(handler.GetVerificationStatus)) // результат проверки подлинности +

	// приходит от второго сервиса
	router.Get("/api/blockchain/data", http.HandlerFunc(handler.GetDataForBlockCreation)) // данные для создания блока +

	router.Post("/api/blockchain/block", http.HandlerFunc(handler.AddBlockchainBlock)) // добавление сформированного блока +
	router.Get("/api/blockchain/block", http.HandlerFunc(handler.GetBlockchainBlock))  // запрос блока из цепи блокчейн +

	router.Get("/api/blockchain/block/verify", http.HandlerFunc(handler.GetBlockVerificationData))  // запрос данных для проверки подлинности +
	router.Post("/api/blockchain/block/verify", http.HandlerFunc(handler.UpdateVerificationStatus)) // отправка результата проверки +

	return router
}
