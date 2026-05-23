-- name: SearchHeroes :many
SELECT 
  h.id,
  h.user_id,
  h.name,
  h.slug,
  h.universe,
  h.alignment,
  h.powers,
  h.description,
  h.image,
  h.is_active,
  h.created_at,
  h.updated_at,
  ts_rank(h.search_vector, plainto_tsquery('english', $1)) as rank
FROM heroes h
WHERE h.search_vector @@ plainto_tsquery('english', $1)
  AND h.is_active = true
  AND (NULLIF($2, '') IS NULL OR h.universe = $2)
  AND (NULLIF($3, '') IS NULL OR h.alignment = $3)
ORDER BY rank DESC, h.created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountSearchHeroes :one
SELECT COUNT(*) as total
FROM heroes h
WHERE h.search_vector @@ plainto_tsquery('english', $1)
  AND h.is_active = true
  AND (NULLIF($2, '') IS NULL OR h.universe = $2)
  AND (NULLIF($3, '') IS NULL OR h.alignment = $3);
