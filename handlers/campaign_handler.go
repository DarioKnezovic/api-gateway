package handlers

import (
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"github.com/DarioKnezovic/api-gateway/utils"
	"github.com/gorilla/mux"
	_ "github.com/markdingo/cslb"
	"log"
	"net/http"
	"regexp"
)

var campaignRouteMapping = map[string]string{
	"/api/campaigns":     "/api/campaign/all",
	"/api/campaign":      "/api/campaign/create",
	"/api/campaign/{id}": "/api/campaign/single/{id}",
}

func CampaignHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP client
	var forwardRequestKey string
	var forwardRequestUrl string
	client := http.Client{}
	cfg := config.LoadConfig()
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	if id != "" {
		forwardRequestKey = ReplaceLastSegmentWithID(r.URL.Path)
		forwardRequestUrl = ReplaceIDInPath(campaignRouteMapping[forwardRequestKey], id)
	} else {
		forwardRequestKey = r.URL.Path
		forwardRequestUrl = campaignRouteMapping[forwardRequestKey]
	}

	fmt.Println("forwardRequestUrl: " + forwardRequestUrl)

	log.Printf("Forwarding %s request to: %s", r.Method, cfg.CampaignServiceURL+forwardRequestUrl)
	forwardRequest, err := http.NewRequest(r.Method, cfg.CampaignServiceURL+forwardRequestUrl, r.Body)
	if err != nil {
		log.Printf("Failed to create forward request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}

	// Set headers from the original request
	forwardRequest.Header = r.Header

	forwardResponse, err := client.Do(forwardRequest)
	if err != nil {
		log.Printf("Failed to forward request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}
	defer forwardResponse.Body.Close()

	for key, values := range forwardResponse.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	log.Printf("%s request %s returned status code %d",
		r.Method, cfg.CampaignServiceURL+campaignRouteMapping[r.URL.Path], forwardResponse.StatusCode)

	err = utils.WriteJSONResponse(w, forwardResponse)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}
}

func ReplaceLastSegmentWithID(path string) string {
	// Define a regular expression pattern to match the last segment
	// of the path that looks like an integer (e.g., "/4").
	pattern := `/(\d+)$`
	r := regexp.MustCompile(pattern)

	// Find the last segment in the path that matches the pattern.
	// Replace it with "{id}" if found.
	newPath := r.ReplaceAllString(path, `/{id}`)

	return newPath
}

func ReplaceIDInPath(path string, id string) string {
	// Define a regular expression pattern to match "{id}".
	pattern := `/{id}`
	r := regexp.MustCompile(pattern)

	// Replace "{id}" with the provided ID.
	newPath := r.ReplaceAllString(path, "/"+id)

	return newPath
}
