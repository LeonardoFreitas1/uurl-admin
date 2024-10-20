-- name: GetVariantsByLanguageTagID :many
SELECT id, created_at, updated_at, variant_tag, description, is_iana_language_sub_tag, instances_on_domains_count
FROM language_tag_variants WHERE language_tag_id = $1;

-- name: InsertVariant :exec
INSERT INTO language_tag_variants (language_tag_id, variant_tag, description, is_iana_language_sub_tag, instances_on_domains_count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);
