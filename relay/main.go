//go:build bare

package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/event", event)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error: http.ListenAndServe: %v\n", err)
	}
}
