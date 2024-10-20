package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LeonardoFreitas1/uurl-admin/internal/handlers"
)

func main() {
	http.HandleFunc("/language/", handlers.LanguageTagHandler)

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
