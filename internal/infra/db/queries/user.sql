-- name: GetUserByEmail :one
SELECT id, email, password_hash FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, email, name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetUserByID :one
SELECT id, email, name, password_hash, image,created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, password_hash = $3
WHERE id = $1;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: UpdateUserImage :exec
UPDATE users
SET image = $2, updated_at = now()
WHERE id = $1;