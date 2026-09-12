package components

import (
	"encoding/json"
	"log"
)

// jsObjectFromJSON parses a JSON object string so that templ can render it
func jsObjectFromJSON(s string) map[string]any {
	var obj map[string]any
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		log.Printf("Invalid FIREBASE_CLIENT_CONFIG: %v", err)
		return map[string]any{}
	}
	return obj
}
