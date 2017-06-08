package api

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, map[string]string{"error": message}, code)
}

func respondWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
