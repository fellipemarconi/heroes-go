-- name: GetUserByEmail :one
SELECT id, email FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, email, name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id;

