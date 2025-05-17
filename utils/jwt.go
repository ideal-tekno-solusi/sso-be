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

func GenerateAuthToken(message, username string, expTime int) (*string, error) {
	privString := viper.GetString("secret.internal.private")
	if privString == "" {
		return nil, fmt.Errorf("private key note found")
	}

	block, _ := pem.Decode([]byte(privString))

	ecKey, _ := x509.ParseECPrivateKey(block.Bytes)

	key, err := jwk.Import(ecKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.NewBuilder().
		Issuer("sso").
		Expiration(time.Now().Add(time.Minute * time.Duration(expTime))).
		Build()
	if err != nil {
		return nil, err
	}

	token.Set("username", username)

	sign, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), key))
	if err != nil {
		return nil, err
	}

	result := string(sign)

	return &result, nil
}
