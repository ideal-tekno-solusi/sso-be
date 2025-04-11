// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :exec
insert into sessions (
    id,
    client_id,
    code_challenge,
    code_challenge_method,
    insert_date
)
values (
    $1,
    $2,
    $3,
    $4,
    now()
)
`

type CreateSessionParams struct {
	ID                  string
	ClientID            string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.Exec(ctx, createSession,
		arg.ID,
		arg.ClientID,
		arg.CodeChallenge,
		arg.CodeChallengeMethod,
	)
	return err
}

const getAuthorization = `-- name: GetAuthorization :one
select
    id,
    user_id
from
    authorization_tokens
where
    id = $1
`

type GetAuthorizationRow struct {
	ID     string
	UserID pgtype.Text
}

func (q *Queries) GetAuthorization(ctx context.Context, id string) (GetAuthorizationRow, error) {
	row := q.db.QueryRow(ctx, getAuthorization, id)
	var i GetAuthorizationRow
	err := row.Scan(&i.ID, &i.UserID)
	return i, err
}

const getUser = `-- name: GetUser :one
select
    id,
    name,
    dot,
    password
from
    users
where
    id = $1
`

type GetUserRow struct {
	ID       string
	Name     string
	Dot      pgtype.Timestamp
	Password string
}

func (q *Queries) GetUser(ctx context.Context, id string) (GetUserRow, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Dot,
		&i.Password,
	)
	return i, err
}
