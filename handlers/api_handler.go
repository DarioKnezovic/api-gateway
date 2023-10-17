package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"github.com/DarioKnezovic/api-gateway/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// RouteInfo Define a struct that matches the structure of your JSON data
type RouteInfo struct {
	Method                 string `json:"method"`
	IncomingPath           string `json:"incoming_path"`
	OutgoingPath           string `json:"outgoing_path"`
	BackendService         string `json:"backend_service"`
	RequiresAuthentication bool   `json:"requires_authentication"`
	DynamicParameter       bool   `json:"dynamic_parameter"`
}

var routes map[string]RouteInfo

const RoutesJsonPath = "./gateway-routes.json"

func init() {
	log.Println("Fetching API Gateway routes from " + RoutesJsonPath)

	jsonFile, err := os.Open(RoutesJsonPath)
	if err != nil {
		log.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	// Read the JSON data from the file
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("Error reading JSON data:", err)
		return
	}

	// Unmarshal the JSON data into the 'routes' map
	err = json.Unmarshal(jsonData, &routes)
	if err != nil {
		log.Println("Error unmarshalling JSON data:", err)
		return
	}

	log.Println("Reading API Gateway routes was done successfully!")
}

func HandleServiceUrl(backendService string) string {
	cfg := config.LoadConfig()

	switch backendService {
	case "UserService":
		return cfg.UserServiceURL
	case "CampaignService":
		return cfg.CampaignServiceURL
	default:
		return ""
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

const (
	INTERNAL_SERVER_ERROR = "Internal server error"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received request for route %s %s", r.Method, r.URL.Path)

	id := mux.Vars(r)["id"]
	var forwardRequestUrl string

	routeKey := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	if id != "" {
		routeKey = ReplaceLastSegmentWithID(routeKey)
	}

	route := routes[routeKey]
	backendUrl := HandleServiceUrl(route.BackendService)
	if backendUrl == "" {
		log.Printf("Cannot find backend URL for this route: %s", routeKey)
	}

	if id == "" {
		forwardRequestUrl = fmt.Sprintf("%s%s", backendUrl, route.OutgoingPath)
	} else {
		forwardRequestUrl = fmt.Sprintf("%s%s", backendUrl, ReplaceIDInPath(route.OutgoingPath, id))
	}

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.Printf("Forwarding %s request to: %s", r.Method, forwardRequestUrl)
	forwardRequest, err := http.NewRequest(r.Method, forwardRequestUrl, r.Body)
	if err != nil {
		log.Printf("Failed to create forward request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}

	// Set headers from the original request
	forwardRequest.Header = r.Header

	// Send the request to the User Management service
	forwardResponse, err := client.Do(forwardRequest)
	if err != nil {
		log.Printf("Failed to forward request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}
	defer forwardResponse.Body.Close()

	// Set the response headers from the User Management service
	for key, values := range forwardResponse.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code and stream the response body to the client
	err = utils.WriteJSONResponse(w, forwardResponse)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}
}
