package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WikipediaTool searches Wikipedia.
type WikipediaTool struct {
	client *http.Client
}

func NewWikipediaTool() *WikipediaTool {
	return &WikipediaTool{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (w *WikipediaTool) Name() string { return "wikipedia" }

func (w *WikipediaTool) Description() string {
	return "search Wikipedia and return summary information"
}

func (w *WikipediaTool) Run(ctx context.Context, input string) (string, error) {
	apiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(input))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "langchaingo/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Wikipedia: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if extract, ok := result["extract"]; ok {
		return fmt.Sprintf("%v", extract), nil
	}
	if title, ok := result["title"]; ok {
		return fmt.Sprintf("Title: %v", title), nil
	}

	return fmt.Sprintf("No results found for: %s", input), nil
}
