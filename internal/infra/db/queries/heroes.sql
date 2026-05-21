-- name: CreateHero :exec
INSERT INTO heroes (
    id,
    user_id,
    name,
    slug,
    universe,
    powers,
    alignment,
    description
) VALUES (
     $1,
     $2,
     $3,
     $4,
     $5,
     $6,
     $7,
     $8
    );

-- name: GetHeroBySlug :one
SELECT id, user_id, name, slug, universe, powers, alignment, description, image, is_active, created_at
FROM heroes
WHERE slug = $1 LIMIT 1;

-- name: ListHeroes :many
SELECT id, user_id, name, slug, universe, powers, alignment, description, image, is_active, created_at
FROM heroes
WHERE
    (sqlc.narg(universe)::text IS NULL OR universe = sqlc.narg(universe))
  AND
    (sqlc.narg(alignment)::text IS NULL OR alignment = sqlc.narg(alignment))
  AND
    (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.narg(is_active))
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);

-- name: UpdateHero :exec
UPDATE heroes
SET
    name = $3,
    universe = $4,
    powers = $5,
    alignment = $6,
    description = $7,
    updated_at = now()
WHERE id = $1 AND user_id = $2;

-- name: SetHeroActiveStatus :exec
UPDATE heroes
SET
    is_active = $3,
    updated_at = now()
WHERE id = $1 AND user_id = $2;

-- name: UpdateHeroImage :exec
UPDATE heroes
SET image = $3, updated_at = now()
WHERE id = $1 AND user_id = $2;