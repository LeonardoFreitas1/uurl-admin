CREATE TABLE country (
    id SERIAL PRIMARY KEY,
    name varchar(100) not null,
    official_state_name varchar(100),
    tld varchar(3) not null,
    iso3166_2_A1 varchar(2) not null,
    iso3166_2_A3 varchar(3) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
