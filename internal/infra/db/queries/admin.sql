-- name: GetStats :one
SELECT
  (SELECT COUNT(*)::int FROM users)  AS users,
  (SELECT COUNT(*)::int FROM heroes) AS heroes,
  (SELECT COUNT(*)::int FROM files)  AS images,
  COALESCE(
    (SELECT SUM((metadata->>'size')::bigint) FROM files),
    0::bigint
)::bigint AS storage_bytes;
