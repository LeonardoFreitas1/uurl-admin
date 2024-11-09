package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/LeonardoFreitas1/uurl-admin/db/sqlc"
)

type GetAllCountriesResponse struct {
	ID                int32  `json:"id"`
	Name              string `json:"name"`
	OfficialStateName string `json:"official_state_name"`
	Tld               string `json:"tld"`
	Iso31662A1        string `json:"iso3166_2_a1"`
	Iso31662A3        string `json:"iso3166_2_a3"`
}

type InsertCountryRequest struct {
	Name              string `json:"name"`
	OfficialStateName string `json:"official_state_name"`
	Tld               string `json:"tld"`
	Iso31662A1        string `json:"iso3166_2_a1"`
	Iso31662A3        string `json:"iso3166_2_a3"`
}

// CountryHandler handles requests for country-related operations
// @Summary Country-related operations
// @Description Handles get, post, and retrieve by ID requests for countries
// @tags Country
// @Accept  json
// @Produce  json
// @Param   id   path   int    false  "Country ID"
// @Success 200  {object}  GetAllCountriesResponse
// @Failure 400  {string}  string  "Invalid item ID"
// @Failure 405  {string}  string  "Method not allowed"
// @Router /country [get]
// @Router /country [post]
func CountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := r.URL.Path

		if path == "/country" || path == "/country/" {
			getAllCountries(w, r)
			return
		}

		idStr := strings.TrimPrefix(path, "/country/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}
		getCountryByID(w, r, int32(id))
	case http.MethodPost:
		createCountry(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllCountries retrieves all countries
// @Summary Get all countries
// @Description Retrieves a list of all countries
// @tags Country
// @Accept  json
// @Produce  json
// @Success 200  {array}   GetAllCountriesResponse
// @Router /country [get]
func getAllCountries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	countries, err := queries.GetAllCountries(ctx)
	if err != nil {
		http.Error(w, "Failed to get countrys", http.StatusInternalServerError)
		return
	}

	var result []GetAllCountriesResponse
	for _, country := range countries {

		if err != nil {
			http.Error(w, "Failed to get variants for country", http.StatusInternalServerError)
			return
		}

		result = append(result, GetAllCountriesResponse{
			ID:                country.ID,
			Name:              country.Name,
			Tld:               country.Tld,
			Iso31662A1:        country.Iso31662A1,
			Iso31662A3:        country.Iso31662A3,
			OfficialStateName: country.OfficialStateName.String,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// getCountryByID retrieves a country by its ID
// @Summary Get country by ID
// @Description Retrieves a country by the provided ID
// @tags Country
// @Accept  json
// @Produce  json
// @Param   id   path   int   true  "Country ID"
// @Success 200  {object}  GetAllCountriesResponse
// @Failure 400  {string}  string  "Invalid item ID"
// @Router /country/{id} [get]
func getCountryByID(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := r.Context()

	country, err := queries.GetCountryById(ctx, id)
	if err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	result := GetAllCountriesResponse{
		ID:                country.ID,
		Name:              country.Name,
		OfficialStateName: country.OfficialStateName.String,
		Iso31662A1:        country.Iso31662A1,
		Iso31662A3:        country.Iso31662A3,
		Tld:               country.Tld,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// createCountry creates a new country
// @Summary Create a new country
// @Description Creates a new country with the provided information
// @tags Country
// @Accept  json
// @Produce  json
// @Param   country  body  InsertCountryRequest  true  "Country Data"
// @Success 201  {object}  GetAllCountriesResponse
// @Failure 400  {string}  string  "Invalid input"
// @Router /country [post]
func createCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input InsertCountryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	countryParams := sqlc.InsertCountryParams{
		Name:              input.Name,
		OfficialStateName: sql.NullString{String: input.OfficialStateName},
		Tld:               input.Tld,
		Iso31662A1:        input.Iso31662A1,
		Iso31662A3:        input.Iso31662A3,
	}

	countryID, err := queries.InsertCountry(ctx, countryParams)
	if err != nil {
		http.Error(w, "Failed to insert country", http.StatusInternalServerError)
		return
	}

	country, err := queries.GetCountryById(ctx, countryID)
	if err != nil {
		http.Error(w, "Failed to retrieve inserted country", http.StatusInternalServerError)
		return
	}

	result := GetAllCountriesResponse{
		ID:                country.ID,
		Name:              country.Name,
		Tld:               country.Tld,
		OfficialStateName: country.OfficialStateName.String,
		Iso31662A1:        country.Iso31662A1,
		Iso31662A3:        country.Iso31662A3,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
