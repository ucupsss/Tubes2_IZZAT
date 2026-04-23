package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NormalizeURL(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}

	if !strings.Contains(trimmed, "://") {
		return "http://" + trimmed
	}

	return trimmed
}

// FetchHTML performs an HTTP GET request to the target URL and returns the response body as a string.
func FetchHTML(rawURL string) (string, error) {
	normalizedURL := NormalizeURL(rawURL)
	parsedURL, err := url.ParseRequestURI(normalizedURL)
	if err != nil {
		return "", fmt.Errorf("format URL tidak valid")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("URL harus menggunakan http:// atau https://")
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL tidak memiliki host yang valid")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request ke URL")
	}
	req.Header.Set("User-Agent", "Tubes2_IZZAT/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengakses URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("website mengembalikan status error: %d", resp.StatusCode)
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	if contentType != "" &&
		!strings.Contains(contentType, "text/html") &&
		!strings.Contains(contentType, "application/xhtml+xml") {
		return "", fmt.Errorf("URL tidak mengembalikan dokumen HTML")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca isi HTML: %v", err)
	}

	body := string(bodyBytes)
	if strings.TrimSpace(body) == "" {
		return "", fmt.Errorf("halaman yang diambil kosong")
	}

	return body, nil
}
