package main

import (
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"log"
	"net/http"

	"github.com/DarioKnezovic/api-gateway/handlers"
	"github.com/DarioKnezovic/api-gateway/middleware"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Running %s on port %s...", cfg.ProjectName, cfg.APIPort)
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware)

	// Routes
	router.HandleFunc("/api/login", handlers.ApiHandler).Methods("POST")
	router.HandleFunc("/api/register", handlers.ApiHandler).Methods("POST")
	router.HandleFunc("/api/logout", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("POST")

	router.HandleFunc("/api/campaigns", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("GET")
	router.HandleFunc("/api/campaign", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("POST")
	router.HandleFunc("/api/campaign/{id}", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("GET")
	router.HandleFunc("/api/campaign/{id}", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("PUT")
	router.HandleFunc("/api/campaign/{id}", middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods("DELETE")
	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), router))
}
