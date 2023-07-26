package functions

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ViewData struct {
	URL     string
	Sitemap string
}

func GetSitemapURL(urlStr string) (string, error) {
	urlHTTP := "http://" + urlStr
	urlHTTPS := "https://" + urlStr

	resp, err := http.Get(urlHTTP + "/robots.txt")
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return GetValidSitemapURL(urlHTTP)
	}

	resp, err = http.Get(urlHTTPS + "/robots.txt")
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return GetValidSitemapURL(urlHTTPS)
	}

	return "", fmt.Errorf("sitemap not found in robots.txt for both http and https versions")
}

func GetValidSitemapURL(urlStr string) (string, error) {
	resp, err := http.Get(urlStr + "/robots.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	robotsTxt := strings.Split(string(body), "\n")
	var sitemapURL string
	for _, line := range robotsTxt {
		if strings.HasPrefix(line, "Sitemap:") {
			sitemapURL = strings.TrimSpace(strings.TrimPrefix(line, "Sitemap:"))
			break
		}
	}

	if sitemapURL == "" {
		return "", fmt.Errorf("sitemap not found in robots.txt")
	}

	parsedSitemapURL, err := url.Parse(sitemapURL)
	if err != nil {
		return "", err
	}
	if !parsedSitemapURL.IsAbs() {
		baseURL, err := url.Parse(urlStr)
		if err != nil {
			return "", err
		}
		sitemapURL = baseURL.ResolveReference(parsedSitemapURL).String()
	}

	return sitemapURL, nil
}

func FetchSitemap(sitemapURL string) ([]byte, error) {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
