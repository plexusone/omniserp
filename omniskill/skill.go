// Package omniskill provides an omniskill wrapper for omniserp search capabilities.
//
// This package adapts the omniserp search client to the omniskill interface,
// allowing web search capabilities to be used by any omniskill-compatible agent.
package omniskill

import (
	"context"
	"fmt"
	"strings"

	"github.com/plexusone/omniserp"
	"github.com/plexusone/omniserp/client"
	"github.com/plexusone/omniskill/skill"
)

// Skill implements omniskill.Skill for web search.
type Skill struct {
	client *client.Client
	config Config
}

// Config configures the search skill.
type Config struct {
	// Provider specifies the search provider (e.g., "serper", "brave", "serpapi").
	// If empty, uses the default from environment variables.
	Provider string

	// MaxResults is the maximum number of results to return (default: 5).
	MaxResults int
}

// New creates a new search skill.
func New(cfg Config) (*Skill, error) {
	c, err := client.NewWithOptions(&client.Options{
		Silent: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create search client: %w", err)
	}

	if cfg.MaxResults == 0 {
		cfg.MaxResults = 5
	}

	return &Skill{
		client: c,
		config: cfg,
	}, nil
}

// Name returns the skill name.
func (s *Skill) Name() string {
	return "search"
}

// Description returns the skill description.
func (s *Skill) Description() string {
	return "Web search tools for finding current information, news, and facts"
}

// Tools returns the tools provided by this skill.
func (s *Skill) Tools() []skill.Tool {
	return []skill.Tool{
		skill.NewTool(
			"web_search",
			"Search the web for current information. Use this when you need up-to-date information, news, or facts.",
			map[string]skill.Parameter{
				"query": {
					Type:        "string",
					Description: "The search query",
					Required:    true,
				},
			},
			s.webSearch,
		),
		skill.NewTool(
			"news_search",
			"Search for recent news articles on a topic.",
			map[string]skill.Parameter{
				"query": {
					Type:        "string",
					Description: "The news search query",
					Required:    true,
				},
			},
			s.newsSearch,
		),
		skill.NewTool(
			"image_search",
			"Search for images related to a query.",
			map[string]skill.Parameter{
				"query": {
					Type:        "string",
					Description: "The image search query",
					Required:    true,
				},
			},
			s.imageSearch,
		),
	}
}

// Init initializes the skill.
func (s *Skill) Init(ctx context.Context) error {
	return nil
}

// Close releases resources.
func (s *Skill) Close() error {
	return nil
}

// webSearch performs a general web search.
func (s *Skill) webSearch(ctx context.Context, params map[string]any) (any, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query is required")
	}

	result, err := s.client.SearchNormalized(ctx, omniserp.SearchParams{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return formatSearchResults(result, s.config.MaxResults), nil
}

// newsSearch performs a news search.
func (s *Skill) newsSearch(ctx context.Context, params map[string]any) (any, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query is required")
	}

	result, err := s.client.SearchNewsNormalized(ctx, omniserp.SearchParams{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("news search failed: %w", err)
	}

	return formatNewsResults(result, s.config.MaxResults), nil
}

// imageSearch performs an image search.
func (s *Skill) imageSearch(ctx context.Context, params map[string]any) (any, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query is required")
	}

	result, err := s.client.SearchImagesNormalized(ctx, omniserp.SearchParams{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("image search failed: %w", err)
	}

	return formatImageResults(result, s.config.MaxResults), nil
}

// formatSearchResults formats web search results as a readable string.
func formatSearchResults(result *omniserp.NormalizedSearchResult, maxResults int) string {
	var sb strings.Builder

	// Answer box if available
	if result.AnswerBox != nil && result.AnswerBox.Answer != "" {
		sb.WriteString("Direct Answer:\n")
		sb.WriteString(fmt.Sprintf("  %s\n\n", result.AnswerBox.Answer))
	}

	// Knowledge graph if available
	if result.KnowledgeGraph != nil && result.KnowledgeGraph.Title != "" {
		sb.WriteString("Knowledge Panel:\n")
		sb.WriteString(fmt.Sprintf("  %s\n", result.KnowledgeGraph.Title))
		if result.KnowledgeGraph.Description != "" {
			sb.WriteString(fmt.Sprintf("  %s\n", result.KnowledgeGraph.Description))
		}
		sb.WriteString("\n")
	}

	// Organic results
	for i, item := range result.OrganicResults {
		if i >= maxResults {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("   URL: %s\n", item.Link))
		if item.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", item.Snippet))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatNewsResults formats news search results.
func formatNewsResults(result *omniserp.NormalizedSearchResult, maxResults int) string {
	var sb strings.Builder

	for i, item := range result.NewsResults {
		if i >= maxResults {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("   Source: %s | %s\n", item.Source, item.Date))
		sb.WriteString(fmt.Sprintf("   URL: %s\n", item.Link))
		if item.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", item.Snippet))
		}
		sb.WriteString("\n")
	}

	if len(result.NewsResults) == 0 {
		sb.WriteString("No news results found.\n")
	}

	return sb.String()
}

// formatImageResults formats image search results.
func formatImageResults(result *omniserp.NormalizedSearchResult, maxResults int) string {
	var sb strings.Builder

	for i, item := range result.ImageResults {
		if i >= maxResults {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("   Image: %s\n", item.ImageURL))
		sb.WriteString(fmt.Sprintf("   Source: %s\n", item.Source))
		sb.WriteString("\n")
	}

	if len(result.ImageResults) == 0 {
		sb.WriteString("No image results found.\n")
	}

	return sb.String()
}
