package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Token interface {
	GetToken(ctx context.Context, sessionId string) (*database.GetTokenRow, error)
	DeleteAuthToken(ctx context.Context, sessionId string) error
	DeleteSession(ctx context.Context, sessionId string) error
}

type TokenService struct {
	Token
}

func TokenRepository(token Token) *TokenService {
	return &TokenService{
		Token: token,
	}
}

func (r *Repository) GetToken(ctx context.Context, sessionId string) (*database.GetTokenRow, error) {
	args := pgtype.Text{
		String: sessionId,
		Valid:  true,
	}

	data, err := r.read.GetToken(ctx, args)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) DeleteAuthToken(ctx context.Context, sessionId string) error {
	args := pgtype.Text{
		String: sessionId,
		Valid:  true,
	}

	err := r.write.DeleteAuthToken(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionId string) error {
	err := r.write.DeleteSession(ctx, sessionId)
	if err != nil {
		return err
	}

	return nil
}
