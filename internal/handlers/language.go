package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
	"github.com/LeonardoFreitas1/uurl-admin/pkg/config"
)

var database = config.GetDB()
var queries = sqlc.New(database)

type LanguageTagWithVariants struct {
	sqlc.LanguageTag
	Variants []sqlc.GetVariantsByLanguageTagIDRow `json:"variants"`
}

func LanguageTagHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("here:")
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Path[len("/language/"):]
		if idStr == "" || idStr == "/language" {
			getAllLanguageTags(w, r)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		getLanguageTagByID(w, r, int32(id))
	case http.MethodPost:
		postLanguageTag(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllLanguageTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	languageTags, err := queries.GetAllLanguageTags(ctx)
	if err != nil {
		http.Error(w, "Failed to get language tags", http.StatusInternalServerError)
		return
	}

	var result []LanguageTagWithVariants
	for _, tag := range languageTags {
		tagIDNull := sql.NullInt32{
			Int32: int32(tag.ID),
			Valid: tag.ID != 0,
		}

		variants, err := queries.GetVariantsByLanguageTagID(ctx, tagIDNull)
		if err != nil {
			http.Error(w, "Failed to get variants for language tag", http.StatusInternalServerError)
			return
		}
		result = append(result, LanguageTagWithVariants{
			LanguageTag: tag,
			Variants:    variants,
		})
	}

	json.NewEncoder(w).Encode(result)
}

func getLanguageTagByID(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := r.Context()

	tag, err := queries.GetLanguageTagByID(ctx, id)
	if err != nil {
		http.Error(w, "Language tag not found", http.StatusNotFound)
		return
	}

	tagIDNull := sql.NullInt32{
		Int32: int32(tag.ID),
		Valid: tag.ID != 0,
	}

	variants, err := queries.GetVariantsByLanguageTagID(ctx, tagIDNull)
	if err != nil {
		http.Error(w, "Failed to get variants", http.StatusInternalServerError)
		return
	}

	result := LanguageTagWithVariants{
		LanguageTag: tag,
		Variants:    variants,
	}

	json.NewEncoder(w).Encode(result)
}

func postLanguageTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Name     string `json:"name"`
		IsoCode1 string `json:"isoCode1"`
		IsoCode2 string `json:"isoCode2"`
		Variants []struct {
			VariantTag              string `json:"variantTag"`
			Description             string `json:"description"`
			IsIANALanguageSubTag    bool   `json:"isIANALanguageSubTag"`
			InstancesOnDomainsCount int    `json:"instancesOnDomainsCount"`
		} `json:"variants"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tagParams := sqlc.InsertLanguageTagParams{
		Name:     input.Name,
		IsoCode1: input.IsoCode1,
		IsoCode2: input.IsoCode2,
	}

	tagID, err := queries.InsertLanguageTag(ctx, tagParams)
	if err != nil {
		http.Error(w, "Failed to insert language tag", http.StatusInternalServerError)
		return
	}

	println(tagID)

	tagIDNull := sql.NullInt32{
		Int32: int32(tagID),
		Valid: tagID != 0,
	}

	for _, variant := range input.Variants {
		descriptionNull := sql.NullString{
			String: variant.Description,
			Valid:  variant.Description != "",
		}

		variantParams := sqlc.InsertVariantParams{
			LanguageTagID:           tagIDNull,
			VariantTag:              variant.VariantTag,
			Description:             descriptionNull,
			IsIanaLanguageSubTag:    variant.IsIANALanguageSubTag,
			InstancesOnDomainsCount: int32(variant.InstancesOnDomainsCount),
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}

		if err := queries.InsertVariant(ctx, variantParams); err != nil {
			http.Error(w, "Failed to insert variant", http.StatusInternalServerError)
			return
		}
	}

	tag, err := queries.GetLanguageTagByID(ctx, tagID)
	if err != nil {
		http.Error(w, "Failed to retrieve inserted language tag", http.StatusInternalServerError)
		return
	}

	variants, err := queries.GetVariantsByLanguageTagID(ctx, tagIDNull)
	if err != nil {
		http.Error(w, "Failed to retrieve variants", http.StatusInternalServerError)
		return
	}

	result := LanguageTagWithVariants{
		LanguageTag: tag,
		Variants:    variants,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
