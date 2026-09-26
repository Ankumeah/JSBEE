//go:build netlify

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", event)

	return mux
}

func handler(
	ctx context.Context,
	input events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	body := input.Body

	if input.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Invalid base64 body",
			}, nil
		}

		body = string(decoded)
	}

	query := url.Values{}

	for key, values := range input.MultiValueQueryStringParameters {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	// Fallback for requests where only single-value query parameters exist.
	if len(query) == 0 {
		for key, value := range input.QueryStringParameters {
			query.Set(key, value)
		}
	}

	requestURL := "https://netlify.local"

	if input.Path != "" {
		requestURL += input.Path
	} else {
		requestURL += "/event"
	}

	if encodedQuery := query.Encode(); encodedQuery != "" {
		requestURL += "?" + encodedQuery
	}

	req, err := http.NewRequestWithContext(
		ctx,
		input.HTTPMethod,
		requestURL,
		bytes.NewBufferString(body),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "could not create request",
		}, nil
	}

	for key, value := range input.Headers {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()

	getMux().ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()

	headers := map[string]string{}

	for key, values := range response.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      response.StatusCode,
		Headers:         headers,
		Body:            recorder.Body.String(),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
