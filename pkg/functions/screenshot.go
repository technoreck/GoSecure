package functions

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func CaptureScreenshot(domain string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	if err := chromedp.Run(ctx, chromedp.Navigate(domain)); err != nil {
		return nil, err
	}

	chromedp.Sleep(2 * time.Second)

	var buf []byte
	if err := chromedp.Run(ctx, ScreenshotTask(&buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

func ScreenshotTask(res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.CaptureScreenshot(res),
	}
}

func GetFileNameFromURL(domain string) string {
	u, err := url.Parse(domain)
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}
	return u.Hostname() + ".png"
}

func SaveScreenshot(fileName string, data []byte) error {
	return os.WriteFile(fileName, data, 0644)
}
