package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchHTML performs an HTTP GET request to the target URL and returns the response body as a string.
func FetchHTML(url string) (string, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("gagal mengakses URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("website mengembalikan status error: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca isi HTML: %v", err)
	}

	return string(bodyBytes), nil
}