package repository

import (
	database "app/database/main"
	"context"
)

type Token interface {
	GetClient(ctx context.Context, id string) (*database.GetClientRow, error)
}

type TokenService struct {
	Token
}

func TokenRepository(token Token) *TokenService {
	return &TokenService{
		Token: token,
	}
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
