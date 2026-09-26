package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Type string         `json:"type"`
	Body map[string]any `json:"body"`
}

func event(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request Request
	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Body == nil {
		request.Body = map[string]any{}
	}

	if !validateRequest(request) {
		return
	}

	if err := sendEmail(request); err != nil {
		log.Printf("Error: sendEmail: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
