package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]interface{}{
		"error": message,
	}

	WriteJSONResponse(w, code, response)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	WriteJSONResponse(w, code, payload)
}

func WriteJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
