package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
)

type LanguageTagVariantsRequest struct {
	LanguageTagID int32  `json:"language_id"`
	CountryID     int32  `json:"country_id"`
	VariantTag    string `json:"variant_tag"`
	Description   string `json:"description"`
}

type LanguageTagVariantsResponse struct {
	ID            int32  `json:"id"`
	LanguageTagID int32  `json:"language_tag_id"`
	VariantTag    string `json:"variant_tag"`
	Description   string `json:"description"`
}

type PaginatedVariantsResponse struct {
	Variants      []LanguageTagVariantsResponse `json:"variants"`
	NextPageToken string                        `json:"next_page_token,omitempty"`
}

// LanguageTagVariantHandler handles requests related to language tag variants
//
//	@Summary		Handles language tag variants
//	@Description	Get, create, or update language tag variants
//	@tags			Language variants
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	false	"Variant ID"	for	PUT	request
//	@Success		200	{object}	LanguageTagVariantsResponse
//	@Failure		400	{string}	string	"Invalid request"
//	@Failure		405	{string}	string	"Method not allowed"
func LanguageTagVariantHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPaginatedVariants(w, r)
	case http.MethodPost:
		postLanguageTagVariant(w, r)
	case http.MethodPut:
		path := r.URL.Path
		idStr := strings.TrimPrefix(path, "/language-variant/")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		updateLanguageTagVariant(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getPaginatedVariants returns paginated language tag variants
//
//	@Summary		Get paginated language tag variants
//	@Description	Get a list of language tag variants with pagination
//	@tags			Language variants
//	@Accept			json
//	@Produce		json
//	@Param			languageTagId	query		int	false	"Language Tag ID"
//	@Param			page_size		query		int	false	"Limit of items per page"	default(10)
//	@Param			page_token		query		int	false	"Offset for pagination"		default(0)
//	@Success		200				{object}	PaginatedVariantsResponse
//	@Failure		400				{string}	string	"Invalid languageTagId or page_token"
//	@Failure		500				{string}	string	"Database query error"
//	@Router			/language-variant [get]
func getPaginatedVariants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	languageTagIdStr := query.Get("languageTagId")
	pageSizeStr := query.Get("page_size")
	pageTokenStr := query.Get("page_token")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	var languageTagId *int32
	if languageTagIdStr != "" {
		id, err := strconv.Atoi(languageTagIdStr)
		if err != nil {
			http.Error(w, "Invalid languageTagId", http.StatusBadRequest)
			return
		}
		id32 := int32(id)
		languageTagId = &id32
	} else {
		languageTagId = nil
	}

	offset := 0
	if pageTokenStr != "" {
		offset, err = strconv.Atoi(pageTokenStr)
		if err != nil || offset < 0 {
			http.Error(w, "Invalid page_token", http.StatusBadRequest)
			return
		}
	}

	var variants []sqlc.Variant
	if languageTagId == nil {
		variants, err = queries.GetPaginatedVariantsWithoutFilter(ctx, sqlc.GetPaginatedVariantsWithoutFilterParams{
			Limit:  int32(pageSize),
			Offset: int32(offset),
		})
	} else {
		variants, err = queries.GetPaginatedVariantsWithFilter(ctx, sqlc.GetPaginatedVariantsWithFilterParams{
			LanguageID: *languageTagId,
			Limit:      int32(pageSize),
			Offset:     int32(offset),
		})
	}

	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	var response PaginatedVariantsResponse
	for _, v := range variants {
		response.Variants = append(response.Variants, LanguageTagVariantsResponse{
			ID:            v.ID,
			LanguageTagID: v.LanguageID.Int32,
			VariantTag:    v.VariantTag,
			Description:   v.Description.String,
		})
	}

	if len(variants) == pageSize {
		response.NextPageToken = strconv.Itoa(offset + pageSize)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// postLanguageTagVariant handles creating a new language tag variant
//
//	@Summary		Create a new language tag variant
//	@Description	Create a new language tag variant
//	@tags			Language variants
//	@Accept			json
//	@Produce		json
//	@Param			variant	body		[]LanguageTagVariantsRequest	true	"Language Tag Variant"
//	@Success		201		{object}	[]LanguageTagVariantsResponse
//	@Failure		400		{string}	string	"Invalid request payload"
//	@Failure		500		{string}	string	"Database query error"
//	@Router			/language-variant [post]
func postLanguageTagVariant(w http.ResponseWriter, r *http.Request) {
	var req []LanguageTagVariantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var response []LanguageTagVariantsResponse
	for _, v := range req {
		arg := sqlc.InsertVariantParams{
			LanguageID:  sql.NullInt32{Int32: v.LanguageTagID, Valid: false},
			CountryID:   sql.NullInt32{Int32: v.CountryID, Valid: false},
			VariantTag:  v.VariantTag,
			Description: sql.NullString{String: v.Description, Valid: true},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := queries.InsertVariant(r.Context(), arg)
		if err != nil {
			println(err.Error())
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}

		response = append(response, LanguageTagVariantsResponse{
			LanguageTagID: v.LanguageTagID,
			VariantTag:    v.VariantTag,
			Description:   v.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// updateLanguageTagVariant handles updating an existing language tag variant
//
//	@Summary		Update an existing language tag variant
//	@Description	Update an existing language tag variant
//	@tags			Language variants
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Variant ID"
//	@Param			variant	body		LanguageTagVariantsRequest	true	"Language Tag Variant"
//	@Success		200		{object}	LanguageTagVariantsResponse
//	@Failure		400		{string}	string	"Invalid request payload"
//	@Failure		404		{string}	string	"Variant not found"
//	@Failure		500		{string}	string	"Database query error"
//	@Router			/language-variant/{id} [put]
func updateLanguageTagVariant(w http.ResponseWriter, r *http.Request, LanguageTagVariantId int) {
	var req LanguageTagVariantsResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	arg := sqlc.UpdateVariantParams{
		ID:          int32(LanguageTagVariantId),
		LanguageID:  sql.NullInt32{Int32: req.LanguageTagID, Valid: true},
		VariantTag:  req.VariantTag,
		Description: sql.NullString{String: req.Description, Valid: true},
		UpdatedAt:   time.Now(),
	}

	err := queries.UpdateVariant(r.Context(), arg)
	if err != nil {
		println(err.Error())
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	response := LanguageTagVariantsResponse{
		ID:            int32(LanguageTagVariantId),
		LanguageTagID: req.LanguageTagID,
		VariantTag:    req.VariantTag,
		Description:   req.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
