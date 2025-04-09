package repository

import (
	database "app/database/main"
	"context"
)

type Authorization interface {
	CreateSession(ctx context.Context, id, clientId, codeChallenge, codeChallengeMethod string) error
}

type AuthorizationService struct {
	Authorization
}

func AuthorizationRepository(authorization Authorization) *AuthorizationService {
	return &AuthorizationService{
		Authorization: authorization,
	}
}

func (r *Repository) CreateSession(ctx context.Context, id, clientId, codeChallenge, codeChallengeMethod string) error {
	args := database.CreateSessionParams{
		ID:                  id,
		ClientID:            clientId,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	err := r.write.CreateSession(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
