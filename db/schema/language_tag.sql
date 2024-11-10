CREATE TABLE language (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             iso_639_1 CHAR(2) NOT NULL,
                             iso_639_2 CHAR(3) NOT NULL
);
