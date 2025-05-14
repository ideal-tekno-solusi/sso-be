package utils

import "golang.org/x/crypto/bcrypt"

func HashBcrypt(message string) (*string, error) {
	var result string

	hash, err := bcrypt.GenerateFromPassword([]byte(message), 15)
	if err != nil {
		return nil, err
	}

	result = string(hash)

	return &result, nil
}

func ValidateHash(hash, message string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(message))
	if err != nil {
		return false
	}

	return true
}
