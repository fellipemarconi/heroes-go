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
SELECT *
FROM files
WHERE path = $1 LIMIT 1;

-- name: DeleteFileByPath :exec
DELETE FROM files
WHERE path = $1;
