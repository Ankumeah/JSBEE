package components

import (
	"maps"
	"slices"

	"github.com/Ankumeah/JSBEE/backend/internal/provider"
)

type pageInfo struct {
	title       string
	description string
}

type socialLink struct {
	name    string
	text    string
	url     string
	isEmail bool
}

// contactEmail returns the raw email address from provider.Soicals
// ("" when absent, so callers can hide email lines as before).
func contactEmail() string {
	if entry, ok := provider.Soicals["Email"]; ok {
		return entry[1]
	}
	return ""
}

// orderedSocials flattens provider.Soicals into a stably-ordered slice
// (Go map iteration is random, so sort by name for a deterministic page).
func orderedSocials() []socialLink {
	links := make([]socialLink, 0, len(provider.Soicals))
	for _, name := range slices.Sorted(maps.Keys(provider.Soicals)) {
		entry := provider.Soicals[name]
		url := entry[1]
		isEmail := name == "Email"
		if isEmail {
			url = "mailto:" + url
		}
		links = append(links, socialLink{name: name, text: entry[0], url: url, isEmail: isEmail})
	}
	return links
}
