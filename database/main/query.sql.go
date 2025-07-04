// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthToken = `-- name: CreateAuthToken :exec
insert into authorization_tokens
(
    id,
    session_id,
    insert_date
)
values
(
    $1,
    $2,
    now()
)
`

type CreateAuthTokenParams struct {
	ID        string
	SessionID pgtype.Text
}

func (q *Queries) CreateAuthToken(ctx context.Context, arg CreateAuthTokenParams) error {
	_, err := q.db.Exec(ctx, createAuthToken, arg.ID, arg.SessionID)
	return err
}

const createRefreshToken = `-- name: CreateRefreshToken :exec
insert into refresh_tokens
(
    id,
    user_id,
    insert_date
)
values
(
    $1,
    $2,
    now()
)
`

type CreateRefreshTokenParams struct {
	ID     string
	UserID pgtype.Text
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, createRefreshToken, arg.ID, arg.UserID)
	return err
}

const createSession = `-- name: CreateSession :exec
insert into sessions (
    id,
    user_id,
    client_id,
    code_challenge,
    code_challenge_method,
    scopes,
    redirect_url,
    insert_date
)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    now()
)
`

type CreateSessionParams struct {
	ID                  string
	UserID              pgtype.Text
	ClientID            string
	CodeChallenge       string
	CodeChallengeMethod string
	Scopes              pgtype.Text
	RedirectUrl         pgtype.Text
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.Exec(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.ClientID,
		arg.CodeChallenge,
		arg.CodeChallengeMethod,
		arg.Scopes,
		arg.RedirectUrl,
	)
	return err
}

const deleteAuthToken = `-- name: DeleteAuthToken :exec
delete from authorization_tokens
where session_id = $1
`

func (q *Queries) DeleteAuthToken(ctx context.Context, sessionID pgtype.Text) error {
	_, err := q.db.Exec(ctx, deleteAuthToken, sessionID)
	return err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
delete from refresh_tokens
where id = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteRefreshToken, id)
	return err
}

const deleteSession = `-- name: DeleteSession :exec
delete from sessions
where id = $1
`

func (q *Queries) DeleteSession(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteSession, id)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
select
    id,
    user_id,
    insert_date
from
    refresh_tokens
where
    id = $1
`

func (q *Queries) GetRefreshToken(ctx context.Context, id string) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshToken, id)
	var i RefreshToken
	err := row.Scan(&i.ID, &i.UserID, &i.InsertDate)
	return i, err
}

const getSession = `-- name: GetSession :one
select
    client_id,
    code_challenge,
    code_challenge_method
from
    sessions
where
    id = $1
`

type GetSessionRow struct {
	ClientID            string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (q *Queries) GetSession(ctx context.Context, id string) (GetSessionRow, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i GetSessionRow
	err := row.Scan(&i.ClientID, &i.CodeChallenge, &i.CodeChallengeMethod)
	return i, err
}

const getToken = `-- name: GetToken :one
select
    auth.id,
    auth.session_id,
    sess.code_challenge,
    sess.scopes,
    sess.redirect_url,
    us.id as username,
    us.name
from
    authorization_tokens auth
join
    sessions sess
on
    auth.session_id = sess.id
join
    users us
on
    sess.user_id = us.id
where
    sess.code_challenge = $1
order by
    auth.insert_date desc
`

type GetTokenRow struct {
	ID            string
	SessionID     pgtype.Text
	CodeChallenge string
	Scopes        pgtype.Text
	RedirectUrl   pgtype.Text
	Username      string
	Name          string
}

func (q *Queries) GetToken(ctx context.Context, codeChallenge string) (GetTokenRow, error) {
	row := q.db.QueryRow(ctx, getToken, codeChallenge)
	var i GetTokenRow
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.CodeChallenge,
		&i.Scopes,
		&i.RedirectUrl,
		&i.Username,
		&i.Name,
	)
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
