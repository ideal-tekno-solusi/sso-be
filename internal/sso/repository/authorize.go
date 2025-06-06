package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Authorize interface {
	GetSession(ctx context.Context, id string) (*database.GetSessionRow, error)
	CreateAuthToken(ctx context.Context, authToken, sessionId string) error
	GetRefreshToken(ctx context.Context, refreshToken string) (*database.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	CreateSession(ctx context.Context, id, userId, clientId, codeChallenge, codeChallengeMethod, scopes, redirectUrl string) error
}

type AuthorizeService struct {
	Authorize
}

func AuthorizeRepository(authorize Authorize) *AuthorizeService {
	return &AuthorizeService{
		Authorize: authorize,
	}
}

func (r *Repository) GetSession(ctx context.Context, id string) (*database.GetSessionRow, error) {
	data, err := r.read.GetSession(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) CreateAuthToken(ctx context.Context, authToken, sessionId string) error {
	args := database.CreateAuthTokenParams{
		ID: authToken,
		SessionID: pgtype.Text{
			String: sessionId,
			Valid:  true,
		},
	}

	err := r.write.CreateAuthToken(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, refreshToken string) (*database.RefreshToken, error) {
	data, err := r.read.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	err := r.write.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
