package repository

import (
	database "app/database/main"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
	mock := &mockRepository{
		getUserFunc: func(ctx context.Context, id string) (database.GetUserRow, error) {
			return database.GetUserRow{ID: "123", Name: "test user"}, nil
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}
	user, err := repo.GetUser(context.Background(), "123")

	assert.NoError(t, err)
	assert.NotZero(t, user)
	assert.Equal(t, "123", user.ID)
}

func TestGetUser_NotFound(t *testing.T) {
	mock := &mockRepository{
		getUserFunc: func(ctx context.Context, id string) (database.GetUserRow, error) {
			return database.GetUserRow{}, pgx.ErrNoRows
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}
	user, err := repo.GetUser(context.Background(), "123")

	assert.NoError(t, err)
	assert.Zero(t, user)
}

func TestGetUser_Error(t *testing.T) {
	mock := &mockRepository{
		getUserFunc: func(ctx context.Context, id string) (database.GetUserRow, error) {
			return database.GetUserRow{}, errors.New("mock error")
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}
	user, err := repo.GetUser(context.Background(), "123")

	assert.Zero(t, user)
	assert.Error(t, err)
}

func TestCreateSession_Success(t *testing.T) {
	mock := &mockRepository{
		createSessionFunc: func(ctx context.Context, arg database.CreateSessionParams) error {
			return nil
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	args := database.CreateSessionParams{
		ID: pgtype.Text{
			String: "mock session",
			Valid:  true,
		},
		UserID: pgtype.Text{
			String: "mock user",
			Valid:  true,
		},
	}

	err := repo.CreateSession(context.Background(), args)

	assert.NoError(t, err)
}

func TestCreateSession_Error(t *testing.T) {
	mock := &mockRepository{
		createSessionFunc: func(ctx context.Context, arg database.CreateSessionParams) error {
			return errors.New("mock error")
		},
	}

	repo := &Repository{
		read:  mock,
		write: mock,
	}

	args := database.CreateSessionParams{
		ID: pgtype.Text{
			String: "mock session",
			Valid:  true,
		},
		UserID: pgtype.Text{
			String: "mock user",
			Valid:  true,
		},
	}

	err := repo.CreateSession(context.Background(), args)

	assert.Error(t, err)
}
