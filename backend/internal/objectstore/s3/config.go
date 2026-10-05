package s3

import (
	"encoding/json"
	"strings"
)

const backupBucket = "backup"
const publicBucket = "public"
const privateBucket = "private"

const backupBaseName = "jsbee.sql.bak."

var buckets []string = []string{
	backupBucket,
	publicBucket,
	privateBucket,
}

type s3StaticConfig struct {
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Region         string `json:"region"`
	MaxDBSnapshots uint   `json:"max_db_snapshots"`
	Url            string `json:"url"`
	Secure         bool   `json:"secure"`
}

func NewS3StaticConfigFromJSON(
	configJSON []byte,
) (s3StaticConfig, error) {
	var config s3StaticConfig
	err := json.Unmarshal(configJSON, &config)

	return config, err
}

func publicBaseURLFromConfig(config s3StaticConfig) string {
	scheme := "https"
	if !config.Secure {
		scheme = "http"
	}

	return scheme + "://" + config.Url + "/" + publicBucket
}

// Content type for stored objects based on file extension
func contentTypeFor(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown"):
		return "text/markdown; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}
