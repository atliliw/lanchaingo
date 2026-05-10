package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DuckDuckGoSearchTool performs web searches.
type DuckDuckGoSearchTool struct {
	client *http.Client
}

func NewDuckDuckGoSearchTool() *DuckDuckGoSearchTool {
	return &DuckDuckGoSearchTool{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *DuckDuckGoSearchTool) Name() string { return "search" }

func (s *DuckDuckGoSearchTool) Description() string {
	return "search the internet and return results"
}

func (s *DuckDuckGoSearchTool) Run(ctx context.Context, input string) (string, error) {
	apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1", url.QueryEscape(input))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "langchaingo/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	content := string(body)
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}
	return content, nil
}
