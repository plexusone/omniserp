// Package brave implements the omniserp.Engine interface for Brave Search API.
package brave

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/plexusone/omniserp"
)

const (
	baseURL       = "https://api.search.brave.com/res/v1"
	engineName    = "brave"
	engineVersion = "1.0.0"
)

// Engine implements the omniserp.Engine interface for Brave Search API.
type Engine struct {
	apiKey string
	client *http.Client
}

// New creates a new Brave Search engine instance using BRAVE_API_KEY env var.
func New() (*Engine, error) {
	apiKey := os.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BRAVE_API_KEY environment variable is required")
	}
	return NewWithAPIKey(apiKey)
}

// NewWithAPIKey creates a new Brave Search engine instance with the provided API key.
func NewWithAPIKey(apiKey string) (*Engine, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	return &Engine{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

// GetName returns the engine name.
func (e *Engine) GetName() string {
	return engineName
}

// GetVersion returns the engine version.
func (e *Engine) GetVersion() string {
	return engineVersion
}

// GetSupportedTools returns the list of supported tools.
func (e *Engine) GetSupportedTools() []string {
	return []string{
		"brave_search",
		"brave_search_news",
		"brave_search_images",
		"brave_search_videos",
	}
}

// makeRequest performs HTTP GET request to Brave Search API.
func (e *Engine) makeRequest(endpoint string, params url.Values) (*omniserp.SearchResult, error) {
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("X-Subscription-Token", e.apiKey)

	// #nosec G704 -- request to hardcoded Brave API endpoint
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &omniserp.SearchResult{
		Data: result,
		Raw:  string(body),
	}, nil
}

// buildParams converts SearchParams to URL query parameters.
func (e *Engine) buildParams(params omniserp.SearchParams) url.Values {
	q := url.Values{}
	q.Set("q", params.Query)

	if params.Country != "" {
		q.Set("country", params.Country)
	}
	if params.Language != "" {
		q.Set("search_lang", params.Language)
	}
	if params.NumResults > 0 && params.NumResults <= 20 {
		q.Set("count", strconv.Itoa(params.NumResults))
	}

	// Enable extra snippets for better LLM context
	q.Set("extra_snippets", "true")

	return q
}

// Search performs a general web search.
func (e *Engine) Search(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return e.makeRequest("/web/search", e.buildParams(params))
}

// SearchNews performs a news search.
func (e *Engine) SearchNews(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return e.makeRequest("/news/search", e.buildParams(params))
}

// SearchImages performs an image search.
func (e *Engine) SearchImages(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return e.makeRequest("/images/search", e.buildParams(params))
}

// SearchVideos performs a video search.
func (e *Engine) SearchVideos(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return e.makeRequest("/videos/search", e.buildParams(params))
}

// SearchPlaces performs a places search (not supported by Brave, returns error).
func (e *Engine) SearchPlaces(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("places search not supported by Brave Search API")
}

// SearchMaps performs a maps search (not supported by Brave, returns error).
func (e *Engine) SearchMaps(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("maps search not supported by Brave Search API")
}

// SearchReviews performs a reviews search (not supported by Brave, returns error).
func (e *Engine) SearchReviews(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("reviews search not supported by Brave Search API")
}

// SearchShopping performs a shopping search (not supported by Brave, returns error).
func (e *Engine) SearchShopping(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("shopping search not supported by Brave Search API")
}

// SearchScholar performs a scholar search (not supported by Brave, returns error).
func (e *Engine) SearchScholar(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("scholar search not supported by Brave Search API")
}

// SearchLens performs a visual search (not supported by Brave, returns error).
func (e *Engine) SearchLens(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("lens search not supported by Brave Search API")
}

// SearchAutocomplete gets search suggestions.
func (e *Engine) SearchAutocomplete(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	if params.Country != "" {
		q.Set("country", params.Country)
	}

	return e.makeRequest("/suggest/search", q)
}

// ScrapeWebpage scrapes content from a webpage (not supported by Brave, returns error).
func (e *Engine) ScrapeWebpage(ctx context.Context, params omniserp.ScrapeParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("webpage scraping not supported by Brave Search API")
}

// Ensure Engine implements omniserp.Engine at compile time.
var _ omniserp.Engine = (*Engine)(nil)
