package main

const maxRetries = 5
const backoffMultiplier = 2

var envVars = map[string]string{
	"PORT":        "",
	"API_VERSION": "",
	"FROM_NAME":   "",
	"FROM_EMAIL":  "",
	"DB_URL":      "",
}
