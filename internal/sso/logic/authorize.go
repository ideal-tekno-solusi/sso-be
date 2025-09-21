package logic

import (
	database "app/database/main"

	"github.com/spf13/viper"
)

type Authorize interface {
	ValidateRedirectUris(urls []database.FetchClientRedirectsRow, url string) bool
}

type AuthorizeService struct {
	Authorize
}

func AuthorizeLogic(authorize Authorize) *AuthorizeService {
	return &AuthorizeService{
		Authorize: authorize,
	}
}

func (l *Logic) ValidateRedirectUris(urls []database.FetchClientRedirectsRow, url string) (valid bool) {
	debug := viper.GetBool("config.debug")

	for _, v := range urls {
		if debug {
			valid = true
			break
		}

		if v.Uri == url {
			valid = true
			break
		}
	}

	return
}
