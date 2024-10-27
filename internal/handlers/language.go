package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
	"github.com/LeonardoFreitas1/uurl-admin/pkg/config"
)

var database = config.GetDB()
var queries = sqlc.New(database)

type LanguageTagGetAllResponse struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	IsoCode1      string `json:"iso_code_1"`
	IsoCode2      string `json:"iso_code_2"`
	VariantsCount int32  `json:"variants_count"`
}

type LanguageTagResponse struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	IsoCode1 string `json:"iso_code_1"`
	IsoCode2 string `json:"iso_code_2"`
}

type LanguageTagBody struct {
	Name     string `json:"name"`
	IsoCode1 string `json:"iso_code_1"`
	IsoCode2 string `json:"iso_code_2"`
}

// LanguageTagHandler godoc
// @Summary Manage Language Tags
// @Description Endpoint to handle operations on language tags by method
// @Tags LanguageTags
// @Accept json
// @Produce json
// @Param id path int false "Language Tag ID"
// @Success 200 {object} LanguageTagResponse "Language Tag with variants"
// @Failure 400 {string} string "Invalid item ID"
// @Failure 405 {string} string "Method not allowed"
// @Router /language/{id} [get]
// @Router /language [post]
func LanguageTagHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := r.URL.Path

		if path == "/language" || path == "/language/" {
			getAllLanguageTags(w, r)
			return
		}

		idStr := strings.TrimPrefix(path, "/language/")

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
// @Success 200 {array} LanguageTagGetAllResponse "List of Language Tags with variants"
// @Failure 500 {string} string "Failed to get language tags"
// @Router /language [get]
func getAllLanguageTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	languageTags, err := queries.GetAllLanguageTags(ctx)
	if err != nil {
		http.Error(w, "Failed to get language tags", http.StatusInternalServerError)
		return
	}

	var result []LanguageTagGetAllResponse
	for _, tag := range languageTags {
		tagIDNull := sql.NullInt32{
			Int32: int32(tag.ID),
			Valid: tag.ID != 0,
		}

		variantCount, err := queries.GetVariantCount(ctx, tagIDNull)
		if err != nil {
			http.Error(w, "Failed to get variants for language tag", http.StatusInternalServerError)
			return
		}

		result = append(result, LanguageTagGetAllResponse{
			ID:            tag.ID,
			IsoCode1:      tag.IsoCode1,
			Name:          tag.Name,
			IsoCode2:      tag.IsoCode2,
			VariantsCount: int32(variantCount),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// getLanguageTagByID godoc
// @Summary Get language tag by ID
// @Description Retrieve a specific language tag and its variants by ID
// @Tags LanguageTags
// @Produce json
// @Param id path int true "Language Tag ID"
// @Success 200 {object} LanguageTagResponse "Language Tag with variants"
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

	result := LanguageTagResponse{
		ID:       tag.ID,
		Name:     tag.Name,
		IsoCode1: tag.IsoCode1,
		IsoCode2: tag.IsoCode2,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// postLanguageTag godoc
// @Summary Create a new language tag
// @Description Insert a new language tag and its associated variants
// @Tags LanguageTags
// @Accept json
// @Produce json
// @Param languageTag body LanguageTagBody true "Language Tag with Variants"
// @Success 201 {object} LanguageTagResponse "Created Language Tag with variants"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to insert language tag or variants"
// @Router /language [post]
func postLanguageTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input LanguageTagBody
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

	tag, err := queries.GetLanguageTagByID(ctx, tagID)
	if err != nil {
		http.Error(w, "Failed to retrieve inserted language tag", http.StatusInternalServerError)
		return
	}

	result := LanguageTagResponse{
		ID:       tag.ID,
		Name:     tag.Name,
		IsoCode1: tag.IsoCode1,
		IsoCode2: tag.IsoCode2,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
