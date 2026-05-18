package utils

import (
	"context"
	"io"
	"net/http"
	"time"
)

// FetchText fetches the content of a URL and returns it as a string.
// Useful for searching keywords in the detail page.
func FetchText(ctx context.Context, url string) string {
	if url == "" {
		return ""
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (canarias.run bot)")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 50000)) // Limit to 50KB to be fast
	if err != nil {
		return ""
	}
	return string(body)
}
