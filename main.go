package main

import (
	"encoding/json"
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"github.com/DarioKnezovic/api-gateway/handlers"
	"github.com/DarioKnezovic/api-gateway/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Running %s on port %s...", cfg.ProjectName, cfg.APIPort)

	jsonData, err := ioutil.ReadFile("./gateway-routes.json")
	if err != nil {
		log.Println("Error reading JSON file: ", err)
		return
	}

	var routes map[string]handlers.RouteInfo

	// Unmarshal the JSON data into the 'routes' map
	err = json.Unmarshal(jsonData, &routes)
	if err != nil {
		log.Println("Error unmarshalling JSON data:", err)
		return
	}

	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware)

	// Routes
	for _, route := range routes {
		if route.RequiresAuthentication {
			router.HandleFunc(route.IncomingPath, middleware.AuthenticationMiddleware(handlers.ApiHandler)).Methods(route.Method)
		} else {
			router.HandleFunc(route.IncomingPath, handlers.ApiHandler).Methods(route.Method)
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(router)

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIPort), handler))
}
