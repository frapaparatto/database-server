package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	writeJSONResponse(w, code, errorResponse{
		Error: msg,
	})
}

func writeJSONResponse(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
		return
	}
	w.WriteHeader(code)
	_, _ = w.Write(data)
}
