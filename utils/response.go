package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// RespondWithError sends an error message as a JSON response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]string{
		"error": message,
	}

	// Convert the response map to a JSON byte slice
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new http.Response object
	resp := &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
		Header:     make(http.Header),
	}

	// Use the existing WriteJSONResponse function
	if err := WriteJSONResponse(w, resp); err != nil {
		log.Printf("Error while sending error response: %v", err)
	}
}

func WriteJSONResponse(w http.ResponseWriter, resp *http.Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		return err
	}
	return nil
}
