package utils

import (
	"time"

	"github.com/google/uuid"
)

type DefaultResponse struct {
	Guid      string `json:"guid"`
	Timestamp string `json:"timestamp"`
	IsSuccess bool   `json:"isSuccess"`
	Data      any    `json:"data"`
}

func GenerateResponseJson(isSuccess bool, data any) DefaultResponse {
	return DefaultResponse{
		Guid:      uuid.NewString(),
		Timestamp: time.Now().Format("2006-01-02 15:04:05 MST"),
		IsSuccess: isSuccess,
		Data:      data,
	}
}
