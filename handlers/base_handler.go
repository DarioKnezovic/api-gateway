package handlers

import (
	"encoding/json"
	"net/http"
)

type BaseHandler struct{}

func (bh *BaseHandler) RespondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]interface{}{
		"error": message,
	}

	WriteJSONResponse(w, code, response)
}

func (bh *BaseHandler) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	WriteJSONResponse(w, code, payload)
}

func WriteJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}