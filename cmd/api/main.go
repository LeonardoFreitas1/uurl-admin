package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/LeonardoFreitas1/uurl-admin/cmd/api/docs"
	"github.com/LeonardoFreitas1/uurl-admin/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title UURL Admin API
// @version 1.0
// @description API documentation for UURL Admin service.
// @host localhost:8080
// @BasePath /
func main() {
	http.HandleFunc("/language/", handlers.LanguageTagHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
