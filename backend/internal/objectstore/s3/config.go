package s3

import "strings"

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
	AccessKey      string
	SecretKey      string
	Region         string
	MaxDBSnapshots uint
	Url            string
	Secure         bool
}

func NewConfig(
	accessKey string,
	secretKey string,
	region string,
	maxDBSnapshots uint,
	url string,
	secure bool,
) s3StaticConfig {
	return s3StaticConfig{
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		Region:         region,
		MaxDBSnapshots: maxDBSnapshots,
		Url:            url,
		Secure:         secure,
	}
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
