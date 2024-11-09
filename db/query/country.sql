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
