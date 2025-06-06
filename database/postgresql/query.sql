-- name: CreateSession :exec
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
);

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

-- name: GetSession :one
select
    client_id,
    code_challenge,
    code_challenge_method
from
    sessions
where
    id = $1;

-- name: CreateAuthToken :exec
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
);

-- name: GetToken :one
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
    auth.insert_date desc;

-- name: DeleteSession :exec
delete from sessions
where id = $1;

-- name: DeleteAuthToken :exec
delete from authorization_tokens
where session_id = $1;

-- name: CreateRefreshToken :exec
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
);

-- name: GetRefreshToken :one
select
    id,
    user_id,
    insert_date
from
    refresh_tokens
where
    id = $1;

-- name: DeleteRefreshToken :exec
delete from refresh_tokens
where id = $1;