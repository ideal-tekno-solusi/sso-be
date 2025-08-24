package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Token interface {
	GetClient(ctx context.Context, id string) (*database.GetClientRow, error)
	GetAuth(ctx context.Context, code string) (*database.GetAuthRow, error)
	GetSession(ctx context.Context, id string) (*database.Session, error)
	UpdateAuth(ctx context.Context, code string) error
	CreateAuth(ctx context.Context, authorizeCode, scope, userId string, authType int) error
}

type TokenService struct {
	Token
}

func TokenRepository(token Token) *TokenService {
	return &TokenService{
		Token: token,
	}
}

func (r *Repository) GetAuth(ctx context.Context, code string) (*database.GetAuthRow, error) {
	data, err := r.read.GetAuth(ctx, pgtype.Text{String: code, Valid: true})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) UpdateAuth(ctx context.Context, code string) error {
	err := r.write.UpdateAuth(ctx, pgtype.Text{String: code, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

// func (r *Repository) GetToken(ctx context.Context, codeChallenge string) (*database.GetTokenRow, error) {
// 	data, err := r.read.GetToken(ctx, codeChallenge)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &data, nil
// }

// func (r *Repository) DeleteAuthToken(ctx context.Context, sessionId string) error {
// 	args := pgtype.Text{
// 		String: sessionId,
// 		Valid:  true,
// 	}

// 	err := r.write.DeleteAuthToken(ctx, args)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *Repository) DeleteSession(ctx context.Context, sessionId string) error {
// 	err := r.write.DeleteSession(ctx, sessionId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *Repository) CreateRefreshToken(ctx context.Context, refreshToken, userId string) error {
// 	args := database.CreateRefreshTokenParams{
// 		ID: refreshToken,
// 		UserID: pgtype.Text{
// 			String: userId,
// 			Valid:  true,
// 		},
// 	}

// 	err := r.write.CreateRefreshToken(ctx, args)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
