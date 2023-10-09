package handlers

import (
	"github.com/DarioKnezovic/api-gateway/config"
	"github.com/DarioKnezovic/api-gateway/utils"
	_ "github.com/markdingo/cslb"
	"log"
	"net/http"
)

var campaignRouteMapping = map[string]string{
	"/api/campaigns": "/api/campaign/all",
}

func CampaignHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP client
	client := http.Client{}
	cfg := config.LoadConfig()

	log.Printf("Forwarding %s request to: %s", r.Method, cfg.CampaignServiceURL+campaignRouteMapping[r.URL.Path])
	forwardRequest, err := http.NewRequest(r.Method, cfg.CampaignServiceURL+campaignRouteMapping[r.URL.Path], r.Body)
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
