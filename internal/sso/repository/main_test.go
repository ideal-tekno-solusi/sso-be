package repository

import (
	database "app/database/main"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type mockRepository struct {
	getUserFunc              func(ctx context.Context, id string) (database.GetUserRow, error)
	createSessionFunc        func(ctx context.Context, arg database.CreateSessionParams) error
	getAuthFunc              func(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error)
	updateAuthFunc           func(ctx context.Context, code pgtype.Text) error
	getClientFunc            func(ctx context.Context, id string) (database.GetClientRow, error)
	fetchClientRedirectsFunc func(ctx context.Context, id string) ([]database.FetchClientRedirectsRow, error)
	getSessionFunc           func(ctx context.Context, id pgtype.Text) (database.Session, error)
	createAuthFunc           func(ctx context.Context, arg database.CreateAuthParams) error
}

func (m *mockRepository) GetUser(ctx context.Context, id string) (database.GetUserRow, error) {
	return m.getUserFunc(ctx, id)
}

func (m *mockRepository) CreateSession(ctx context.Context, arg database.CreateSessionParams) error {
	return m.createSessionFunc(ctx, arg)
}

func (m *mockRepository) GetAuth(ctx context.Context, code pgtype.Text) (database.GetAuthRow, error) {
	return m.getAuthFunc(ctx, code)
}

func (m *mockRepository) UpdateAuth(ctx context.Context, code pgtype.Text) error {
	return m.updateAuthFunc(ctx, code)
}

func (m *mockRepository) GetClient(ctx context.Context, id string) (database.GetClientRow, error) {
	return m.getClientFunc(ctx, id)
}

func (m *mockRepository) FetchClientRedirects(ctx context.Context, id string) ([]database.FetchClientRedirectsRow, error) {
	return m.fetchClientRedirectsFunc(ctx, id)
}

func (m *mockRepository) GetSession(ctx context.Context, id pgtype.Text) (database.Session, error) {
	return m.getSessionFunc(ctx, id)
}

func (m *mockRepository) CreateAuth(ctx context.Context, arg database.CreateAuthParams) error {
	return m.createAuthFunc(ctx, arg)
}
