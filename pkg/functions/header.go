package functions

import (
	"net/http"
	"net/url"
)

func HasProtocolScheme(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != ""
}

func GetRemoteData(urlParam string) (map[string][]string, error) {
	resp, err := http.Get(urlParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	headers := resp.Header
	return headers, nil
}
