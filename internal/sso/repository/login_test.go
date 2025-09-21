package repository

import (
	database "app/database/main"
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
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
	assert.NotNil(t, user)
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
