package main

import (
	"log"
	"net/http"

	"github.com/DarioKnezovic/api-gateway/handlers"
	"github.com/DarioKnezovic/api-gateway/middleware"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Running server...")
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware)

	// Routes
	router.HandleFunc("/login", handlers.AuthHandler).Methods("POST")
	router.HandleFunc("/register", handlers.AuthHandler).Methods("POST")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
