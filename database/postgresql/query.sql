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