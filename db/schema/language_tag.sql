CREATE TABLE language_tags (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             iso_code_1 CHAR(2) NOT NULL,
                             iso_code_2 CHAR(3) NOT NULL
);
