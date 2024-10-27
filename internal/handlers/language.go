package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
	"github.com/LeonardoFreitas1/uurl-admin/pkg/config"
)

var database = config.GetDB()
var queries = sqlc.New(database)

type CustomString struct {
	Value string `json:"value"`
}

type VariantWithCustomString struct {
	VariantTag              string       `json:"variantTag"`
	Description             CustomString `json:"description"`
	IsIANALanguageSubTag    bool         `json:"isIANALanguageSubTag"`
	InstancesOnDomainsCount int          `json:"instancesOnDomainsCount"`
}

type LanguageTagWithVariants struct {
	sqlc.LanguageTag
	Variants []VariantWithCustomString `json:"variants"`
}

// LanguageTagHandler godoc
// @Summary Manage Language Tags
// @Description Endpoint to handle operations on language tags by method
// @Tags LanguageTags
// @Accept json
// @Produce json
// @Param id path int false "Language Tag ID"
// @Success 200 {object} LanguageTagWithVariants "Language Tag with variants"
// @Failure 400 {string} string "Invalid item ID"
// @Failure 405 {string} string "Method not allowed"
// @Router /language/{id} [get]
// @Router /language [post]
func LanguageTagHandler(w http.ResponseWriter, r *http.Request) {
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

// getAllLanguageTags godoc
// @Summary Get all language tags
// @Description Retrieve all language tags with their associated variants
// @Tags LanguageTags
// @Produce json
// @Success 200 {array} LanguageTagWithVariants "List of Language Tags with variants"
// @Failure 500 {string} string "Failed to get language tags"
// @Router /language [get]
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

		var variantList []VariantWithCustomString
		for _, variant := range variants {
			variantList = append(variantList, VariantWithCustomString{
				VariantTag:              variant.VariantTag,
				Description:             CustomString{Value: variant.Description.String},
				IsIANALanguageSubTag:    variant.IsIanaLanguageSubTag,
				InstancesOnDomainsCount: int(variant.InstancesOnDomainsCount),
			})
		}

		result = append(result, LanguageTagWithVariants{
			LanguageTag: tag,
			Variants:    variantList,
		})
	}

	json.NewEncoder(w).Encode(result)
}

// getLanguageTagByID godoc
// @Summary Get language tag by ID
// @Description Retrieve a specific language tag and its variants by ID
// @Tags LanguageTags
// @Produce json
// @Param id path int true "Language Tag ID"
// @Success 200 {object} LanguageTagWithVariants "Language Tag with variants"
// @Failure 404 {string} string "Language tag not found"
// @Failure 500 {string} string "Failed to get variants"
// @Router /language/{id} [get]
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

	var variantList []VariantWithCustomString
	for _, variant := range variants {
		variantList = append(variantList, VariantWithCustomString{
			VariantTag:              variant.VariantTag,
			Description:             CustomString{Value: variant.Description.String},
			IsIANALanguageSubTag:    variant.IsIanaLanguageSubTag,
			InstancesOnDomainsCount: int(variant.InstancesOnDomainsCount),
		})
	}

	result := LanguageTagWithVariants{
		LanguageTag: tag,
		Variants:    variantList,
	}

	json.NewEncoder(w).Encode(result)
}

// postLanguageTag godoc
// @Summary Create a new language tag
// @Description Insert a new language tag and its associated variants
// @Tags LanguageTags
// @Accept json
// @Produce json
// @Param languageTag body LanguageTagWithVariants true "Language Tag with Variants"
// @Success 201 {object} LanguageTagWithVariants "Created Language Tag with variants"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to insert language tag or variants"
// @Router /language [post]
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

	var variantList []VariantWithCustomString
	for _, variant := range variants {
		variantList = append(variantList, VariantWithCustomString{
			VariantTag:              variant.VariantTag,
			Description:             CustomString{Value: variant.Description.String},
			IsIANALanguageSubTag:    variant.IsIanaLanguageSubTag,
			InstancesOnDomainsCount: int(variant.InstancesOnDomainsCount),
		})
	}

	result := LanguageTagWithVariants{
		LanguageTag: tag,
		Variants:    variantList,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
