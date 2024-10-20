-- name: GetAllLanguageTags :many
SELECT id, name, iso_code_1, iso_code_2 FROM language_tags;

-- name: GetLanguageTagByID :one
SELECT id, name, iso_code_1, iso_code_2 FROM language_tags WHERE id = $1;

-- name: InsertLanguageTag :one
INSERT INTO language_tags (name, iso_code_1, iso_code_2) VALUES ($1, $2, $3) RETURNING id;
