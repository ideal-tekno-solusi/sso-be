package repository

import (
	database "app/database/main"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestGetAuth_Success(t *testing.T) {
	mock := &mockRepository{
		getAuthFunc: func(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error) {
			return database.GetAuthRow{
				Code: pgtype.Text{
					String: "testcodehere",
					Valid:  true,
				},
				Scope: pgtype.Text{
					String: "auth",
					Valid:  true,
				},
				Type: "code",
				UserID: pgtype.Text{
					String: "testuser",
					Valid:  true,
				},
			}, nil
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	auth, err := repo.GetAuth(context.Background(), pgtype.Text{String: "testcodehere", Valid: true})

	assert.NoError(t, err)
	assert.NotZero(t, auth)
	assert.Equal(t, pgtype.Text{String: "testcodehere", Valid: true}, auth.Code)
}

func TestGetAuth_NotFound(t *testing.T) {
	mock := &mockRepository{
		getAuthFunc: func(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error) {
			return database.GetAuthRow{}, nil
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	auth, err := repo.GetAuth(context.Background(), pgtype.Text{String: "testcodehere", Valid: true})

	assert.NoError(t, err)
	assert.Zero(t, auth)
}

func TestGetAuth_Error(t *testing.T) {
	mock := &mockRepository{
		getAuthFunc: func(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error) {
			return database.GetAuthRow{}, errors.New("mock error")
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	auth, err := repo.GetAuth(context.Background(), pgtype.Text{String: "testcodehere", Valid: true})

	assert.Zero(t, auth)
	assert.Error(t, err)
}

func TestUpdateAuth_Success(t *testing.T) {
	mock := &mockRepository{
		updateAuthFunc: func(ctx context.Context, code pgtype.Text) error {
			return nil
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	err := repo.UpdateAuth(context.Background(), pgtype.Text{String: "testcodehere", Valid: true})

	assert.NoError(t, err)
}

func TestUpdateAuth_Error(t *testing.T) {
	mock := &mockRepository{
		updateAuthFunc: func(ctx context.Context, code pgtype.Text) error {
			return errors.New("mock error")
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	err := repo.UpdateAuth(context.Background(), pgtype.Text{String: "testcodehere", Valid: true})

	assert.Error(t, err)
}
