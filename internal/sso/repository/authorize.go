package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Authorize interface {
	GetClient(ctx context.Context, id string) (*database.GetClientRow, error)
	FetchClientRedirects(ctx context.Context, id string) (*[]database.FetchClientRedirectsRow, error)
	GetSession(ctx context.Context, id string) (*database.Session, error)
	CreateAuth(ctx context.Context, authorizeCode, userId string) error
}

type AuthorizeService struct {
	Authorize
}

func AuthorizeRepository(authorize Authorize) *AuthorizeService {
	return &AuthorizeService{
		Authorize: authorize,
	}
}

func (r *Repository) GetClient(ctx context.Context, id string) (*database.GetClientRow, error) {
	data, err := r.read.GetClient(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) FetchClientRedirects(ctx context.Context, id string) (*[]database.FetchClientRedirectsRow, error) {
	data, err := r.read.FetchClientRedirects(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) GetSession(ctx context.Context, id string) (*database.Session, error) {
	args := pgtype.Text{
		String: id,
		Valid:  true,
	}

	data, err := r.read.GetSession(ctx, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}

func (r *Repository) CreateAuth(ctx context.Context, authorizeCode, userId string) error {
	args := database.CreateAuthParams{
		Code: pgtype.Text{
			String: authorizeCode,
			Valid:  true,
		},
		UserID: pgtype.Text{
			String: userId,
			Valid:  true,
		},
	}

	err := r.write.CreateAuth(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
