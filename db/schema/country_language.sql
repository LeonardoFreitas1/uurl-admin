CREATE TABLE country_language (
                                  country_id INT NOT NULL,
                                  language_id INT NOT NULL,
                                  PRIMARY KEY (country_id, language_id),
                                  FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE CASCADE,
                                  FOREIGN KEY (language_id) REFERENCES language(id) ON DELETE CASCADE
);

