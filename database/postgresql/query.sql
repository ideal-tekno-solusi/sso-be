-- name: CreateSession :exec
insert into sessions (
    id,
    client_id,
    code_challenge,
    code_challenge_method,
    scopes,
    insert_date
)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
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
    auth.session_id = $1
order by
    auth.insert_date desc;

-- name: DeleteSession :exec
delete from sessions
where id = $1;

-- name: DeleteAuthToken :exec
delete from authorization_tokens
where session_id = $1;

-- name: UpdateUserIdSession :exec
update sessions
set
    user_id = $1
where
    id = $2;