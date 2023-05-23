package main

import (
	"github.com/DarioKnezovic/api-gateway/config"
	"log"
	"net/http"

	"github.com/DarioKnezovic/api-gateway/handlers"
	"github.com/DarioKnezovic/api-gateway/middleware"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Running %s ...", cfg.ProjectName)
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware)

	// Routes
	router.HandleFunc("/login", handlers.AuthHandler).Methods("POST")
	router.HandleFunc("/register", handlers.AuthHandler).Methods("POST")

	// Start the server
	log.Fatal(http.ListenAndServe(cfg.APIPort, router))
}
