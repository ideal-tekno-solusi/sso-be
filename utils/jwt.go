package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/spf13/viper"
)

func GenerateAuthToken(payloads map[string]any, expTime int) (*string, error) {
	privString := viper.GetString("secret.internal.private")
	if privString == "" {
		return nil, fmt.Errorf("private key not found")
	}

	block, _ := pem.Decode([]byte(privString))

	ecKey, _ := x509.ParseECPrivateKey(block.Bytes)

	key, err := jwk.Import(ecKey)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	exp := now.Add(time.Minute * time.Duration(expTime))

	token, err := jwt.NewBuilder().
		Issuer("sso").
		Expiration(exp).
		Build()
	if err != nil {
		return nil, err
	}

	payloads["exp"] = exp.Unix()
	payloads["iat"] = now.Unix()

	for k, v := range payloads {
		token.Set(k, v)
	}

	sign, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), key))
	if err != nil {
		return nil, err
	}

	result := string(sign)

	return &result, nil
}

func ValidateJwt(message string) (bool, error) {
	pubString := viper.GetString("secret.internal.public")
	if pubString == "" {
		return false, fmt.Errorf("public key not found")
	}

	block, _ := pem.Decode([]byte(pubString))

	ecKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	key, err := jwk.PublicRawKeyOf(ecKey)
	if err != nil {
		return false, err
	}

	_, err = jwt.Parse([]byte(message), jwt.WithKey(jwa.ES256(), key), jwt.WithValidate(true))
	if err != nil {
		return false, err
	}

	return true, nil
}

func DecryptJwt(message string) (*map[string]any, error) {
	pubString := viper.GetString("secret.internal.public")
	if pubString == "" {
		return nil, fmt.Errorf("public key not found")
	}

	block, _ := pem.Decode([]byte(pubString))

	ecKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	key, err := jwk.PublicRawKeyOf(ecKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse([]byte(message), jwt.WithKey(jwa.ES256(), key), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}

	keys := token.Keys()
	payloads := make(map[string]any)

	for _, v := range keys {
		var data any

		token.Get(v, &data)

		payloads[v] = data
	}

	return &payloads, nil
}
