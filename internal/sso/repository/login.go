package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
)

type Login interface {
	GetUser(ctx context.Context, id string) (*database.GetUserRow, error)
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
