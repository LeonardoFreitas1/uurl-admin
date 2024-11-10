-- name: GetAllLanguageTags :many
SELECT id, name, iso_639_1, iso_639_2 FROM language;

-- name: GetLanguageTagByID :one
SELECT id, name, iso_639_1, iso_639_2 FROM language WHERE id = $1;

-- name: InsertLanguageTag :one
INSERT INTO language (name, iso_639_1, iso_639_2) VALUES ($1, $2, $3) RETURNING id;
