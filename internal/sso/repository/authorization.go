package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
)

type Authorization interface {
	CreateSession(ctx context.Context, id, clientId, codeChallenge, codeChallengeMethod string) error
	GetAuthorization(ctx context.Context, id string) (*database.GetAuthorizationRow, error)
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

func (r *Repository) GetAuthorization(ctx context.Context, id string) (*database.GetAuthorizationRow, error) {
	data, err := r.read.GetAuthorization(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}
