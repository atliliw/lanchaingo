package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// URLFetchTool fetches content from URLs.
type URLFetchTool struct {
	client *http.Client
}

func NewURLFetchTool() *URLFetchTool {
	return &URLFetchTool{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (u *URLFetchTool) Name() string { return "url_fetch" }

func (u *URLFetchTool) Description() string {
	return "fetch content from a URL and return the text"
}

func (u *URLFetchTool) Run(ctx context.Context, input string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Limit output to 2000 chars
	content := string(body)
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}

	return content, nil
}
