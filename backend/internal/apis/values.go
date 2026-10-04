package apis

import (
	"errors"
	"strings"
	"time"
)

const uploadURLTTL = 15 * time.Minute

// Longest accepted location/category for a paper, in characters.
// Mirrors the CHECK constraints added to the papers table in migration v7.
const maxPaperFieldLength = 120

// validatePaperMetadata trims and checks the location and category
// submitted with a paper. Both are required; overlong values are
// rejected so inserts can never violate the v7 CHECK constraints.
func validatePaperMetadata(location, category string) (string, string, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", "", errors.New("Location is required")
	}
	if len([]rune(location)) > maxPaperFieldLength {
		return "", "", errors.New("Location must be at most 120 characters")
	}

	category = strings.TrimSpace(category)
	if category == "" {
		return "", "", errors.New("Category is required")
	}
	if len([]rune(category)) > maxPaperFieldLength {
		return "", "", errors.New("Category must be at most 120 characters")
	}

	return location, category, nil
}
