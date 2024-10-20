CREATE TABLE language_tag_variants (
                          id SERIAL PRIMARY KEY,
                          language_tag_id INT REFERENCES language_tags(id) ON DELETE CASCADE,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          variant_tag VARCHAR(255) NOT NULL,
                          description TEXT,
                          is_iana_language_sub_tag BOOLEAN NOT NULL,
                          instances_on_domains_count INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_variants_language_tag_id ON language_tag_variants(language_tag_id);
