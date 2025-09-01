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

-- name: GetClient :one
select
    c.id,
    c.name,
    t.name as type,
    c.secret,
    c.token_livetime
from
    clients as c
join
    client_types as t
on
    c.type = t.id
where
    c.id = $1;

-- name: FetchClientRedirects :many
select
    c.id,
    r.uri
from
    client_redirects as r
join
    clients as c
on
    r.client_id = c.id
where
    c.id = $1;

-- name: CreateSession :exec
insert into sessions
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

-- name: UpdateSession :exec
update sessions
set
    user_id = $1
where
    id = $2;

-- name: CreateAuth :exec
insert into auths
(
    code,
    scope,
    type,
    user_id,
    insert_date
)
values
(
    $1,
    $2,
    $3,
    $4,
    now()
);

-- name: GetSession :one
select
    id,
    user_id,
    insert_date
from
    sessions
where
    id = $1;

-- name: GetAuth :one
select
    a.code,
    a.scope,
    t.name as type,
    a.user_id
from
    auths as a
join
    auth_types as t
on
    a.type = t.id
where
    code = $1
and
    use_date is null;

-- name: UpdateAuth :exec
update auths
set
    use_date = now()
where
    code = $1;