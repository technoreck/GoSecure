package handlers

import (
	"GoSecure/pkg/functions"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/miekg/dns"
)

func Handler(c *fiber.Ctx) error {
	// Get the hostname from the form data submitted in the request body
	form := new(struct {
		Hostname string `form:"hostname"`
	})

	if err := c.BodyParser(form); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
	}

	hostname := form.Hostname

	// Handle URLs by extracting hostname
	if strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://") {
		u, err := url.Parse(hostname)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid URL")
		}
		hostname = u.Hostname()
	}

	// Perform DNS lookup using the performDNSLookup function
	resp, err := functions.PerformDNSLookup(hostname)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("DNS lookup error")
	}

	// Set the response Content-Type to application/json
	c.Set("Content-Type", "application/json")

	// Encode the response object as JSON and send it back to the client
	return c.JSON(resp)
}

func Headerhandler(c *fiber.Ctx) error {
	urlParam := c.Query("url")

	if urlParam == "" {
		response := map[string]string{"error": "url query string parameter is required"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// If the URL does not start with a protocol scheme, add "http://" as the default protocol
	if !functions.HasProtocolScheme(urlParam) {
		urlParam = "http://" + urlParam
	}

	headers, err := functions.GetRemoteData(urlParam)
	if err != nil {
		fmt.Println("Error:", err.Error())
		response := map[string]string{"error": err.Error()}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(headers)
}

func ScanHandler(c *fiber.Ctx) error {
	hostname := c.FormValue("hostname")

	results, err := functions.ScanPorts(hostname)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.SendString(fmt.Sprintf("Open ports with service names: %v", results))
}

func HstsHandler(c *fiber.Ctx) error {
	urlString := c.FormValue("url")

	if urlString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": "URL parameter is missing!",
		})
	}

	// Check and add the URL scheme if missing
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": "Invalid URL format!",
		})
	}

	result, err := functions.CheckHSTS(u.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error checking HSTS policy: %s", err.Error()),
		})
	}

	return c.JSON(result)
}

func Servstatushandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}

	result, err := functions.CheckServerStatus(url)
	if err != nil {
		return c.SendString(fmt.Sprintf("Error: %s", err.Error()))
	}

	return c.SendString(result)
}

func Dnssechandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}
	rrsigRecords, dnskeyRecords, err := functions.GetRRSIGWithKey(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error while querying DNS records: %v", err.Error()),
		})
	}
	responseData := struct {
		RRIGRecords   []*dns.RRSIG
		DNSKEYRecords []*dns.DNSKEY
	}{
		RRIGRecords:   rrsigRecords,
		DNSKEYRecords: dnskeyRecords,
	}

	return c.JSON(responseData)
}

func Screenshothandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}
	screenshotData, err := functions.CaptureScreenshot(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error while capturing screenshot: %v", err.Error()),
		})
	}

	responseData := struct {
		ScreenshotBase64 string
	}{
		ScreenshotBase64: base64.StdEncoding.EncodeToString(screenshotData),
	}

	return c.JSON(responseData)
}
func Dnsserverhandler(c *fiber.Ctx) error {
	domain := c.FormValue("url")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	addresses, err := net.LookupHost(domain)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("An error occurred while resolving DNS. %s", err.Error()),
		})
	}

	results := make([]functions.DNSResult, 0)
	for _, address := range addresses {
		hostname, err := net.LookupAddr(address)
		if err != nil {
			hostname = nil
		}

		dohDirectSupports := functions.CheckDOHSupport(fmt.Sprintf("https://%s/dns-query", address))

		results = append(results, functions.DNSResult{
			Address:          address,
			Hostname:         functions.GetHostname(hostname),
			DOHDirectSupport: dohDirectSupports,
		})
	}

	response := functions.Response{
		Domain: domain,
		DNS:    results,
	}

	return c.JSON(response)
}

func SSLhandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}

	sslInfo, err := functions.GetSSLInfo(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error while fetching SSL Details: %v", err.Error()),
		})
	}
	return c.JSON(sslInfo)
}

func CookieHandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}

	cookieInfo, err := functions.FindCookies(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error while fetching cookie: %v", err.Error()),
		})
	}
	return c.JSON(cookieInfo)
}

func WhoisHandler(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.SendString("You must provide a URL!")
	}

	whoisInfo, err := functions.GetWHOISInfo(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": fmt.Sprintf("Error fetching WHOIS information: %v", err.Error()),
		})
	}

	return c.JSON(whoisInfo)
}
