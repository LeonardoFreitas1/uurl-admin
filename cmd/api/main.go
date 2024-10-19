package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/LeonardoFreitas1/uurl-admin/internal/db"
	"github.com/LeonardoFreitas1/uurl-admin/internal/models"
	"github.com/LeonardoFreitas1/uurl-admin/pkg/config"
	_ "github.com/lib/pq"
)

var database = config.GetDB()

func main() {
	http.HandleFunc("/language/", languageTagHandler)

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func languageTagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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

		getLanguageTagByID(w, r, id)
	} else if r.Method == http.MethodPost {
		postLanguageTag(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllLanguageTags(w http.ResponseWriter, r *http.Request) {
	languageTags, err := db.GetAllLanguageTags(database)
	if err != nil {
		http.Error(w, "Failed to get language tags", http.StatusInternalServerError)
		return
	}

	for i, tag := range languageTags {
		variants, err := db.GetVariantsByLanguageTagID(database, tag.ID)
		if err != nil {
			http.Error(w, "Failed to get variants for language tag", http.StatusInternalServerError)
			return
		}
		languageTags[i].Variants = variants
	}

	json.NewEncoder(w).Encode(languageTags)
}

func getLanguageTagByID(w http.ResponseWriter, r *http.Request, id int) {
	tag, err := db.GetLanguageTagByID(database, id)
	if err != nil {
		println(err.Error())
		http.Error(w, "Language tag not found", http.StatusNotFound)
		return
	}

	variants, err := db.GetVariantsByLanguageTagID(database, id)
	if err != nil {
		http.Error(w, "Failed to get variants", http.StatusInternalServerError)
		return
	}

	tag.Variants = variants

	json.NewEncoder(w).Encode(tag)
}

func postLanguageTag(w http.ResponseWriter, r *http.Request) {
	var tag models.LanguageTag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newID, err := db.InsertLanguageTag(database, tag)
	if err != nil {
		http.Error(w, "Failed to insert language tag", http.StatusInternalServerError)
		return
	}

	for _, variant := range tag.Variants {
		variant.LanguageTagID = newID
		err := db.InsertVariant(database, variant)
		if err != nil {
			println(err.Error())
			http.Error(w, "Failed to insert variant", http.StatusInternalServerError)
			return
		}
	}

	tag.ID = newID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}
