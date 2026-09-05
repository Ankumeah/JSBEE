package components

import "encoding/json"

// jsObjectFromJSON parses a JSON object string so that templ can render it
func jsObjectFromJSON(s string) map[string]any {
	var obj map[string]any
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return map[string]any{}
	}
	return obj
}
