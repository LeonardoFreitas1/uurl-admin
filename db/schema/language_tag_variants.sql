CREATE TABLE variant (
                          id SERIAL PRIMARY KEY,
                          language_id INT REFERENCES language_tags(id) ON DELETE CASCADE,
                          country_id INT REFERENCES country(id) ON DELETE SET NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          variant_tag VARCHAR(255) NOT NULL,
                          description TEXT
);

CREATE INDEX idx_variants_language_id ON language_tag_variants(language_id);
