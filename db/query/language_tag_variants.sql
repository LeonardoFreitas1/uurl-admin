-- name: GetVariantsByLanguageTagID :many
SELECT id, created_at, updated_at, variant_tag, description
FROM variant WHERE language_id = $1;

-- name: InsertVariant :exec
INSERT INTO variant (language_id, variant_tag, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateVariant :exec
UPDATE variant set language_id = $2, variant_tag = $3, description = $4, updated_at = $5 where id = $1;

-- name: GetVariantCount :one
SELECT count(id) FROM variant WHERE language_id = $1;

-- name: GetPaginatedVariantsWithFilter :many
SELECT * FROM variant
WHERE language_id = sqlc.arg(language_id)::integer
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetPaginatedVariantsWithoutFilter :many
SELECT * FROM variant
ORDER BY id
LIMIT $1 OFFSET $2;
