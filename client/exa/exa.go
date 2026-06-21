// Package exa implements the omniserp.Engine interface for Exa.ai Search API.
package exa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/plexusone/omniserp"
)

const (
	baseURL       = "https://api.exa.ai"
	engineName    = "exa"
	engineVersion = "1.0.0"
)

// SearchType represents the type of search to perform.
type SearchType string

const (
	// SearchTypeAuto intelligently selects the best search mode.
	SearchTypeAuto SearchType = "auto"
	// SearchTypeInstant provides sub-200ms latency for real-time applications.
	SearchTypeInstant SearchType = "instant"
	// SearchTypeFast uses lower-latency search models (~450ms).
	SearchTypeFast SearchType = "fast"
	// SearchTypeDeep performs in-depth research with synthesis.
	SearchTypeDeep SearchType = "deep"
	// SearchTypeDeepLite is lightweight synthesis.
	SearchTypeDeepLite SearchType = "deep-lite"
	// SearchTypeDeepReasoning adds more reasoning to deep search.
	SearchTypeDeepReasoning SearchType = "deep-reasoning"
)

// ContentOptions specifies what content to retrieve from search results.
type ContentOptions struct {
	Text       bool `json:"text,omitempty"`
	Highlights bool `json:"highlights,omitempty"`
	Summary    bool `json:"summary,omitempty"`
}

// SearchRequest represents a request to the Exa search API.
type SearchRequest struct {
	Query          string          `json:"query"`
	NumResults     int             `json:"numResults,omitempty"`
	Type           SearchType      `json:"type,omitempty"`
	IncludeDomains []string        `json:"includeDomains,omitempty"`
	ExcludeDomains []string        `json:"excludeDomains,omitempty"`
	Contents       *ContentOptions `json:"contents,omitempty"`
	SafeSearch     bool            `json:"safeSearch,omitempty"`
}

// Engine implements the omniserp.Engine interface for Exa.ai Search API.
type Engine struct {
	apiKey     string
	client     *http.Client
	searchType SearchType
}

// Option is a functional option for configuring the Engine.
type Option func(*Engine)

// WithSearchType sets the default search type for the engine.
func WithSearchType(st SearchType) Option {
	return func(e *Engine) {
		e.searchType = st
	}
}

// New creates a new Exa Search engine instance using EXA_API_KEY env var.
func New(opts ...Option) (*Engine, error) {
	apiKey := os.Getenv("EXA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("EXA_API_KEY environment variable is required")
	}
	return NewWithAPIKey(apiKey, opts...)
}

// NewWithAPIKey creates a new Exa Search engine instance with the provided API key.
func NewWithAPIKey(apiKey string, opts ...Option) (*Engine, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	e := &Engine{
		apiKey:     apiKey,
		client:     &http.Client{},
		searchType: SearchTypeAuto,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
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
		"exa_search",
		"exa_search_news",
	}
}

// makeRequest performs HTTP POST request to Exa Search API.
func (e *Engine) makeRequest(endpoint string, reqBody any) (*omniserp.SearchResult, error) {
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+endpoint, strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", e.apiKey)

	// #nosec G704 -- request to hardcoded Exa API endpoint
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

// buildRequest converts SearchParams to Exa SearchRequest.
func (e *Engine) buildRequest(params omniserp.SearchParams) SearchRequest {
	req := SearchRequest{
		Query:      params.Query,
		Type:       e.searchType,
		SafeSearch: true,
		Contents: &ContentOptions{
			Text:       true,
			Highlights: true,
			Summary:    true,
		},
	}

	if params.NumResults > 0 && params.NumResults <= 100 {
		req.NumResults = params.NumResults
	} else {
		req.NumResults = 10
	}

	return req
}

// Search performs a general web search.
func (e *Engine) Search(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return e.makeRequest("/search", e.buildRequest(params))
}

// SearchNews performs a news search.
func (e *Engine) SearchNews(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	req := e.buildRequest(params)
	// Filter to news domains
	req.IncludeDomains = []string{
		"reuters.com",
		"apnews.com",
		"bbc.com",
		"cnn.com",
		"theguardian.com",
		"nytimes.com",
		"washingtonpost.com",
		"wsj.com",
		"bloomberg.com",
	}
	return e.makeRequest("/search", req)
}

// SearchImages performs an image search (not supported by Exa, returns error).
func (e *Engine) SearchImages(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("image search not supported by Exa API")
}

// SearchVideos performs a video search (not supported by Exa, returns error).
func (e *Engine) SearchVideos(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("video search not supported by Exa API")
}

// SearchPlaces performs a places search (not supported by Exa, returns error).
func (e *Engine) SearchPlaces(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("places search not supported by Exa API")
}

// SearchMaps performs a maps search (not supported by Exa, returns error).
func (e *Engine) SearchMaps(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("maps search not supported by Exa API")
}

// SearchReviews performs a reviews search (not supported by Exa, returns error).
func (e *Engine) SearchReviews(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("reviews search not supported by Exa API")
}

// SearchShopping performs a shopping search (not supported by Exa, returns error).
func (e *Engine) SearchShopping(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("shopping search not supported by Exa API")
}

// SearchScholar performs a scholar search.
func (e *Engine) SearchScholar(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	req := e.buildRequest(params)
	// Filter to academic domains
	req.IncludeDomains = []string{
		"scholar.google.com",
		"arxiv.org",
		"pubmed.ncbi.nlm.nih.gov",
		"sciencedirect.com",
		"nature.com",
		"science.org",
		"ieee.org",
		"acm.org",
	}
	return e.makeRequest("/search", req)
}

// SearchLens performs a visual search (not supported by Exa, returns error).
func (e *Engine) SearchLens(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("lens search not supported by Exa API")
}

// SearchAutocomplete gets search suggestions (not supported by Exa, returns error).
func (e *Engine) SearchAutocomplete(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("autocomplete not supported by Exa API")
}

// ScrapeWebpage scrapes content from a webpage (not supported by Exa, returns error).
func (e *Engine) ScrapeWebpage(ctx context.Context, params omniserp.ScrapeParams) (*omniserp.SearchResult, error) {
	return nil, fmt.Errorf("webpage scraping not supported by Exa API, use GetContents endpoint instead")
}

// Ensure Engine implements omniserp.Engine at compile time.
var _ omniserp.Engine = (*Engine)(nil)
