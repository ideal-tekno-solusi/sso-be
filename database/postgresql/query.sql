-- name: CreateSession :exec
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
);

-- name: GetAuthorization :one
select
    id,
    user_id
from
    authorization_tokens
where
    id = $1;

-- name: GetUser :one
select
    id,
    name,
    dot,
    password
from
    users
where
    id = $1;