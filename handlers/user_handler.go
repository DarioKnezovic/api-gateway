package handlers

import (
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"io/ioutil"
	"log"
	"net/http"
)

var routeMapping = map[string]string{
	"/register": "/api/user/register",
	"/login":    "/api/user/login",
}

// UserHandler handles requests related to user management
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP client
	client := http.Client{}
	cfg := config.LoadConfig()

	fmt.Println("We have following request")
	fmt.Println(r)
	// Create a new request to forward to the User Management service
	forwardRequest, err := http.NewRequest(r.Method, cfg.UserServiceURL+routeMapping[r.URL.Path], r.Body)
	if err != nil {
		log.Printf("Failed to create forward request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set headers from the original request
	forwardRequest.Header = r.Header

	// Send the request to the User Management service
	forwardResponse, err := client.Do(forwardRequest)
	if err != nil {
		log.Printf("Failed to forward request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer forwardResponse.Body.Close()

	// Read the response body from the User Management service
	responseBody, err := ioutil.ReadAll(forwardResponse.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers from the User Management service
	for key, values := range forwardResponse.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code and write the response body to the client
	w.WriteHeader(forwardResponse.StatusCode)
	w.Write(responseBody)
}
