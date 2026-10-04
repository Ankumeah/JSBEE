package apis

import (
	"strings"
	"testing"
)

func TestValidatePaperMetadata(t *testing.T) {
	cases := []struct {
		name     string
		location string
		category string
		wantErr  string
		wantLoc  string
		wantCat  string
	}{
		{name: "valid", location: "Mumbai, India", category: "Economics", wantLoc: "Mumbai, India", wantCat: "Economics"},
		{name: "trims whitespace", location: "  Berlin  ", category: "\tPhysics\n", wantLoc: "Berlin", wantCat: "Physics"},
		{name: "empty location", location: "", category: "Economics", wantErr: "Location is required"},
		{name: "blank location", location: "   ", category: "Economics", wantErr: "Location is required"},
		{name: "empty category", location: "Mumbai, India", category: "", wantErr: "Category is required"},
		{name: "blank category", location: "Mumbai, India", category: "  \t ", wantErr: "Category is required"},
		{name: "location too long", location: strings.Repeat("a", 121), category: "Economics", wantErr: "Location must be at most 120 characters"},
		{name: "category too long", location: "Mumbai, India", category: strings.Repeat("a", 121), wantErr: "Category must be at most 120 characters"},
		{name: "location at limit", location: strings.Repeat("a", 120), category: "Economics", wantLoc: strings.Repeat("a", 120), wantCat: "Economics"},
		// 120 emoji are 480 bytes but only 120 characters: length is counted in runes.
		{name: "unicode counted in runes", location: strings.Repeat("🌍", 120), category: "Economics", wantLoc: strings.Repeat("🌍", 120), wantCat: "Economics"},
		{name: "unicode over limit", location: strings.Repeat("🌍", 121), category: "Economics", wantErr: "Location must be at most 120 characters"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loc, cat, err := validatePaperMetadata(tc.location, tc.category)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("Expected error %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if loc != tc.wantLoc || cat != tc.wantCat {
				t.Fatalf("Got %q, %q; want %q, %q", loc, cat, tc.wantLoc, tc.wantCat)
			}
		})
	}
}
