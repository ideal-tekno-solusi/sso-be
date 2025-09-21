package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Token interface {
	GetAuth(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error)
	UpdateAuth(ctx context.Context, code pgtype.Text) error
}

type TokenService struct {
	Token
}

func TokenRepository(token Token) *TokenService {
	return &TokenService{
		Token: token,
	}
}

func (r *Repository) GetAuth(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error) {
	data, err := r.read.GetAuth(ctx, code)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) UpdateAuth(ctx context.Context, code pgtype.Text) error {
	err := r.write.UpdateAuth(ctx, code)
	if err != nil {
		return err
	}
