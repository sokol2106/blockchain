package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/ivan/blockchain/api-server/internal/middleware"
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

func (h *Handlers) AddDataBlockchain(res http.ResponseWriter, req *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.srvBlockchain.AddData(string(body))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(handlerStatus)
	res.Write([]byte(result))
}

func (h *Handlers) GetDataBlockchain(res http.ResponseWriter, req *http.Request) {
	data, err := h.srvBlockchain.ReceiveData()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(data))
}

func (h *Handlers) CheckData(res http.ResponseWriter, req *http.Request) {
	_, cancel := context.WithCancel(req.Context())
	defer cancel()

	key := chi.URLParam(req, "key")
	res.Header().Set("Location", key)
	res.WriteHeader(http.StatusOK)

}

func (h *Handlers) StatusProcessCheckData(res http.ResponseWriter, req *http.Request) {
	_, cancel := context.WithCancel(req.Context())
	defer cancel()

	type StatusResponse struct {
		Status  string `json:"status"`
		QueueID string `json:"queueID"`
	}

	sts := StatusResponse{
		Status: "INVALID",
	}

	sts.QueueID = chi.URLParam(req, "queue_id")

	jsonData, err := json.Marshal(sts)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

func (h *Handlers) AddBlock(res http.ResponseWriter, req *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.srvBlockchain.AddBlock(string(body))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(handlerStatus)
}

func (h *Handlers) GetBlock(res http.ResponseWriter, req *http.Request) {
	block, err := h.srvBlockchain.ReceiveBlock()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(block))
}

func (h *Handlers) GetCheckDataBlock(res http.ResponseWriter, req *http.Request) {

}

func (h *Handlers) SetStatusProcessCheckData(res http.ResponseWriter, req *http.Request) {

}

func Router(handler *Handlers) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.СompressionResponseRequest)
	router.Use(middleware.LoggingResponseRequest)

	// router

	// приходит внешне
	router.Post("/api/data", http.HandlerFunc(handler.AddDataBlockchain))                 // новые данные для сохранения
	router.Post("/api/check/{key}", http.HandlerFunc(handler.CheckData))                  // ключ + данные для проверки подлинности
	router.Get("/api/check/{queue_id}", http.HandlerFunc(handler.StatusProcessCheckData)) // результат проверки подлинности

	// приходит от второго сервиса
	router.Get("/api/data", http.HandlerFunc(handler.GetDataBlockchain)) // данные для создания блока

	router.Post("/api/block", http.HandlerFunc(handler.AddBlock)) // добавление сформированного блока
	router.Get("/api/block", http.HandlerFunc(handler.GetBlock))  // запрос блока из цепи блокчейн

	router.Get("/api/block/check", http.HandlerFunc(handler.GetCheckDataBlock))          // запрос данных для проверки подлинности
	router.Post("/api/block/check", http.HandlerFunc(handler.SetStatusProcessCheckData)) // отправка результата проверки

	return router
}
