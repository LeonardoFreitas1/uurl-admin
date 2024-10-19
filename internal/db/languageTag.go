package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/LeonardoFreitas1/uurl-admin/internal/models"
)

func GetAllLanguageTags(db *sql.DB) ([]models.LanguageTag, error) {
	rows, err := db.Query("SELECT id, name, iso_code_1, iso_code_2 FROM language_tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languageTags []models.LanguageTag
	for rows.Next() {
		var tag models.LanguageTag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.IsoCode1, &tag.IsoCode2); err != nil {
			log.Println(err)
			continue
		}
		languageTags = append(languageTags, tag)
	}
	return languageTags, nil
}

func GetLanguageTagByID(db *sql.DB, id int) (models.LanguageTag, error) {
	var tag models.LanguageTag
	query := "SELECT id, name, iso_code_1, iso_code_2 FROM language_tags WHERE id = $1"

	row := db.QueryRow(query, id)
	err := row.Scan(&tag.ID, &tag.Name, &tag.IsoCode1, &tag.IsoCode2)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func GetVariantsByLanguageTagID(db *sql.DB, languageTagID int) ([]models.Variants, error) {
	query := `
		SELECT id, created_at, updated_at, variant_tag, description, is_iana_language_sub_tag, instances_on_domains_count
		FROM language_tag_variants
		WHERE language_tag_id = $1
	`
	rows, err := db.Query(query, languageTagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []models.Variants
	for rows.Next() {
		var variant models.Variants
		err := rows.Scan(
			&variant.ID,
			&variant.CreatedAt,
			&variant.UpdatedAt,
			&variant.VariantTag,
			&variant.Description,
			&variant.IsIANALanguageSubTag,
			&variant.InstancesOnDomainsCount,
		)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	return variants, nil
}

func InsertLanguageTag(db *sql.DB, tag models.LanguageTag) (int, error) {
	var newID int
	query := `
		INSERT INTO language_tags (name, iso_code_1, iso_code_2)
		VALUES ($1, $2, $3) RETURNING id
	`
	err := db.QueryRow(query, tag.Name, tag.IsoCode1, tag.IsoCode2).Scan(&newID)
	return newID, err
}

func InsertVariant(db *sql.DB, variant models.Variants) error {
	query := `
        INSERT INTO language_tag_variants (language_tag_id, variant_tag, description, is_iana_language_sub_tag, instances_on_domains_count, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := db.Exec(query,
		variant.LanguageTagID,
		variant.VariantTag,
		variant.Description,
		variant.IsIANALanguageSubTag,
		variant.InstancesOnDomainsCount,
		time.Now(),
		time.Now(),
	)
	return err
}
