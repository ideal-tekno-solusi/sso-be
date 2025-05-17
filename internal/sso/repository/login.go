package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Login interface {
	GetUser(ctx context.Context, id string) (*database.GetUserRow, error)
	GetSession(ctx context.Context, id string) (*database.GetSessionRow, error)
	CreateAuthToken(ctx context.Context, authToken, sessionId string) error
	UpdateUserIdSession(ctx context.Context, userId, id string) error
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

func (r *Repository) UpdateUserIdSession(ctx context.Context, userId, id string) error {
	args := database.UpdateUserIdSessionParams{
		UserID: pgtype.Text{
			String: userId,
			Valid:  true,
		},
		ID: id,
	}

	err := r.write.UpdateUserIdSession(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
