-- name: GetAllCountries :many
SELECT id, name, official_state_name, tld, iso3166_2_A1, iso3166_2_A3 FROM country;

-- name: InsertCountry :one
INSERT INTO country(
    name,
    official_state_name,
    tld,
    iso3166_2_a1,
    iso3166_2_a3,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: GetCountryById :one
SELECT id, name, official_state_name, tld, iso3166_2_A1, iso3166_2_A3 FROM country where id = $1;

-- name: GetFilteredCountry :many
SELECT DISTINCT ON (ctr.id) id, name, official_state_name, tld, iso3166_2_A1, iso3166_2_A3
FROM country ctr
         LEFT JOIN country_language cl ON ctr.id = cl.country_id
WHERE ($1::int[] IS NULL OR language_id = ANY($1::int[]));
