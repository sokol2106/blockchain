package model

type Status int

const (
	StatusCreated                 Status = iota // Создан
	StatusProcessing                            // Обрабатывается
	StatusMatched                               // Совпадает
	StatusFailedAuthenticityCheck               // Не прошёл проверку подлинности
	StatusNotFound                              // Не найдено
)

func (s Status) String() string {
	switch s {
	case StatusCreated:
		return "Created"
	case StatusProcessing:
		return "Processing"
	case StatusMatched:
		return "Matched"
	case StatusFailedAuthenticityCheck:
		return "Failed Authenticity Check"
	case StatusNotFound:
		return "Status Not Found"
	default:
		return "Unknown"
	}
}

type VerificationData struct {
	QueueId string `json:"queueID"`
	Status  Status `json:"status"`
	Key     string `json:"key"`
	Data    string `json:"data"`
	Block   Block  `json:"block"`
}
