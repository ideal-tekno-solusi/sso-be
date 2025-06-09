package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type DefaultResponse struct {
	Guid      string `json:"guid"`
	Timestamp string `json:"timestamp"`
	IsSuccess bool   `json:"isSuccess"`
	Data      any    `json:"data"`
}

func GenerateResponseJson(baseResponse *DefaultResponse, isSuccess bool, data any) DefaultResponse {
	if baseResponse == nil {
		return DefaultResponse{
			Guid:      uuid.NewString(),
			Timestamp: time.Now().Format("2006-01-02 15:04:05 MST"),
			IsSuccess: isSuccess,
			Data:      data,
		}
	}

	return DefaultResponse{
		Guid:      baseResponse.Guid,
		Timestamp: baseResponse.Timestamp,
		IsSuccess: baseResponse.IsSuccess,
		Data:      data,
	}
}

func BindResponse(data []byte) (*DefaultResponse, *[]byte, error) {
	res := DefaultResponse{}

	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, nil, err
	}

	body, err := defaultResponseMapToByte(res.Data)
	if err != nil {
		return nil, nil, err
	}

	return &res, body, nil
}

func defaultResponseMapToByte(data any) (*[]byte, error) {
	if _, ok := data.(map[string]any); !ok {
		return nil, fmt.Errorf("data type is %v, only accept map[string]any", reflect.TypeOf(data))
	}

	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
