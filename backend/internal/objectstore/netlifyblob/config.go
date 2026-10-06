package netlifyblob

import "time"

type Config struct {
	SiteID         string
	Token          string
	APIURL         string
	Store          string
	Timeout        time.Duration
	MaxDBSnapshots uint
	UploadSecret   string
}

func NewConfig(
	siteID string,
	token string,
	apiURL string,
	store string,
	timeout time.Duration,
	maxDBSnapshots uint,
	uploadSecret string,
) Config {
	return Config{
		SiteID:         siteID,
		Token:          token,
		APIURL:         apiURL,
		Store:          store,
		Timeout:        timeout,
		MaxDBSnapshots: maxDBSnapshots,
		UploadSecret:   uploadSecret,
	}
}
