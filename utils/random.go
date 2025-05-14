package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) (*string, error) {
	var result string

	randByte := make([]byte, length)
	_, err := rand.Read(randByte)
	if err != nil {
		return nil, err
	}

	result = base64.URLEncoding.EncodeToString(randByte)

	return &result, nil
}
