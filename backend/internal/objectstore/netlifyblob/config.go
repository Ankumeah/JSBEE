package netlifyblob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	SiteID string
	Token string
	APIURL string
	Store string
	Timeout time.Duration
	MaxDBSnapshots uint
	UploadSecret string
}

type rawConfig struct {
	SiteID         string `json:"site_id"`
	Token          string `json:"token"`
	APIURL         string `json:"api_url"`
	Store          string `json:"store"`
	TimeoutMS      int64  `json:"timeout_ms"`
	MaxDBSnapshots uint   `json:"max_db_snapshots"`
	UploadSecret   string `json:"upload_secret"`
}

// Parses a netlifyblob config from JSON
func ParseConfig(raw []byte) (Config, error) {
	var r rawConfig
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		return Config{}, fmt.Errorf("netlifyblob: bad config: %w", err)
	}

	return Config{
		SiteID:         strings.TrimSpace(r.SiteID),
		Token:          strings.TrimSpace(r.Token),
		APIURL:         strings.TrimSpace(r.APIURL),
		Store:          strings.TrimSpace(r.Store),
		MaxDBSnapshots: r.MaxDBSnapshots,
		UploadSecret:   r.UploadSecret,
	}, nil
}
