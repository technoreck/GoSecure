package functions

import (
	"fmt"
	"net/http"
)

func CheckHSTS(url string) (map[string]interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err.Error())
	}
	defer response.Body.Close()

	headers := response.Header
	hstsHeader := headers.Get("Strict-Transport-Security")

	if hstsHeader == "" {
		return map[string]interface{}{"message": "Site does not serve any HSTS headers.", "compatible": false}, nil
	}

	if len(hstsHeader) < 10886400 {
		return map[string]interface{}{"message": "HSTS max-age is less than 10886400.", "compatible": false}, nil
	}

	if !ContainsSubDomains(hstsHeader) {
		return map[string]interface{}{"message": "HSTS header does not include all subdomains.", "compatible": false}, nil
	}

	if !ContainsPreload(hstsHeader) {
		return map[string]interface{}{"message": "HSTS header does not contain the preload directive.", "compatible": false}, nil
	}

	return map[string]interface{}{
		"message":    "Site is compatible with the HSTS preload list!",
		"compatible": true,
		"hstsHeader": hstsHeader,
	}, nil
}

func ContainsSubDomains(hstsHeader string) bool {
	return hstsHeader == "includeSubDomains" || hstsHeader == "includeSubDomains; preload"
}

func ContainsPreload(hstsHeader string) bool {
	return hstsHeader == "preload" || hstsHeader == "includeSubDomains; preload"
}
