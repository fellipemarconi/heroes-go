-- name: CreateFile :exec
INSERT INTO files (
    id,
    path,
    type,
    metadata,
    user_id
) VALUES (
     $1,
     $2,
     $3,
     $4,
     $5
);


-- name: GetFileByPath :one
SELECT id, path, type, metadata, created_at
FROM files
WHERE path = $1 LIMIT 1;

-- name: DeleteFileByPath :exec
DELETE FROM files
WHERE path = $1 AND user_id = $2;

-- name: ListFilePathsByUser :many
SELECT path
FROM files
WHERE user_id = $1;
