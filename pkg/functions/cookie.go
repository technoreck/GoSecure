package functions

import (
	"net/http"
)

type CookieInfo struct {
	Status  string `json:"status"`
	Cookies []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"cookies"`
}

// FindCookies retrieves cookie information for the given domain.
func FindCookies(domain string) (*CookieInfo, error) {
	client := &http.Client{}
	url := "https://" + domain
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cookieInfo := &CookieInfo{
		Status: resp.Status,
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		cookieInfo.Cookies = append(cookieInfo.Cookies, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	return cookieInfo, nil
}
