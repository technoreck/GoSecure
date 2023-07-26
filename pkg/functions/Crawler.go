package functions

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

func AddSchemaIfMissing(siteURL string) string {
	if !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		return "https://" + siteURL
	}
	return siteURL
}

func FetchRobotsTXT(siteURL string) ([]byte, error) {
	parsedURL, err := url.Parse(siteURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url query parameter")
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Hostname())

	statusCode, body, err := fasthttp.Get(nil, robotsURL)
	if err != nil {
		return nil, err
	}

	if statusCode == fasthttp.StatusOK {
		return body, nil
	}

	return nil, fmt.Errorf("failed to fetch robots.txt, statusCode: %d", statusCode)
}
