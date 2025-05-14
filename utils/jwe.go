package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/viper"
)

func EncryptJwe(message string, clientId string) (*string, error) {
	pubKey := viper.GetString(fmt.Sprintf("secret.%v.public", clientId))
	if pubKey == "" {
		return nil, fmt.Errorf("public key of %v not found", clientId)
	}

	block, _ := pem.Decode([]byte(pubKey))

	ecKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	key, err := jwk.PublicRawKeyOf(ecKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key with error: %v", err)
	}

	encrypted, err := jwe.Encrypt([]byte(message), jwe.WithKey(jwa.ECDH_ES(), key))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message with error: %v", err)
	}

	ciphertext := string(encrypted)

	return &ciphertext, nil
}
