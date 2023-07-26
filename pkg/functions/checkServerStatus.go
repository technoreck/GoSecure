package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func CheckServerStatus(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %s", err.Error())
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	startTime := time.Now()
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return "", fmt.Errorf("error during operation: %s", err.Error())
	}
	defer resp.Body.Close()

	responseTime := time.Since(startTime).Seconds()

	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 400 {
		return fmt.Sprintf("Received non-success response code: %d", statusCode), nil
	}

	return fmt.Sprintf("Server is up!\nResponse Code: %d\nResponse Time: %.2f seconds", statusCode, responseTime), nil
}
