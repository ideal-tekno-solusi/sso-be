package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Login interface {
	GetUser(ctx context.Context, id string) (*database.GetUserRow, error)
	CreateSession(ctx context.Context, id, userId, clientId, codeChallenge, codeChallengeMethod, scopes string) error
}

type LoginService struct {
	Login
}

func LoginRepository(login Login) *LoginService {
	return &LoginService{
		Login: login,
	}
}

func (r *Repository) GetUser(ctx context.Context, id string) (*database.GetUserRow, error) {
	data, err := r.read.GetUser(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) CreateSession(ctx context.Context, id, userId, clientId, codeChallenge, codeChallengeMethod, scopes string) error {
	args := database.CreateSessionParams{
		ID: id,
		UserID: pgtype.Text{
			String: userId,
			Valid:  true,
		},
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
