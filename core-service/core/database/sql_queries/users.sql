-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
Insert into users (email, password, enabled) values ($1, $2, $3) RETURNING *;

-- name: EmailTaken :one
SELECT EXISTS( SELECT 1 FROM users WHERE email = $1 );
