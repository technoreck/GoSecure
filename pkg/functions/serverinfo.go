package functions

import (
	"net/http"
)

type DNSResult struct {
	Address          string `json:"address"`
	Hostname         string `json:"hostname"`
	DOHDirectSupport bool   `json:"dohDirectSupports"`
}

type Response struct {
	Domain            string      `json:"domain"`
	DNS               []DNSResult `json:"dns"`
	DOHMozillaSupport bool        `json:"dohMozillaSupport"`
}

func CheckDOHSupport(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func GetHostname(hostnames []string) string {
	if len(hostnames) > 0 {
		return hostnames[0]
	}
	return ""
}
