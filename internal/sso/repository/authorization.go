package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Authorization interface {
	CreateSession(ctx context.Context, id, clientId, codeChallenge, codeChallengeMethod, scopes string) error
}

type AuthorizationService struct {
	Authorization
}

func AuthorizationRepository(authorization Authorization) *AuthorizationService {
	return &AuthorizationService{
		Authorization: authorization,
	}
}

func (r *Repository) CreateSession(ctx context.Context, id, clientId, codeChallenge, codeChallengeMethod, scopes string) error {
	args := database.CreateSessionParams{
		ID:                  id,
		ClientID:            clientId,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Scopes: pgtype.Text{
			String: scopes,
			Valid:  true,
		},
	}

	err := r.write.CreateSession(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
