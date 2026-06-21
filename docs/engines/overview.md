# Supported Engines

OmniSerp supports multiple search engine backends through a unified interface.

## Available Engines

### Serper

- **Package**: `github.com/plexusone/omniserp/client/serper`
- **Environment Variable**: `SERPER_API_KEY`
- **Website**: [serper.dev](https://serper.dev)
- **Supported Operations**: All 12 search types including Lens

### SerpAPI

- **Package**: `github.com/plexusone/omniserp/client/serpapi`
- **Environment Variable**: `SERPAPI_API_KEY`
- **Website**: [serpapi.com](https://serpapi.com)
- **Supported Operations**: 11 search types (no Lens support)

!!! note
    `SearchLens()` is not supported by SerpAPI and will return `ErrOperationNotSupported`

### Brave Search

- **Package**: `github.com/plexusone/omniserp/client/brave`
- **Environment Variable**: `BRAVE_API_KEY`
- **Website**: [brave.com/search/api](https://brave.com/search/api)
- **Supported Operations**: Web, News, Images, Videos, Autocomplete

Brave Search provides privacy-focused search results with excellent coverage and fast response times. It includes:

- **Goggles** - Customizable result filtering
- **Summarizer** - AI-powered result summaries
- **Discussions** - Forum and community results

```go
import "github.com/plexusone/omniserp/client/brave"

engine := brave.New(brave.Config{
    APIKey: os.Getenv("BRAVE_API_KEY"),
})

// Web search
result, err := engine.Search(ctx, client.SearchParams{
    Query:      "golang concurrency patterns",
    NumResults: 10,
})

// News search
news, err := engine.SearchNews(ctx, client.SearchParams{
    Query: "tech news",
})
```

### Exa.ai

- **Package**: `github.com/plexusone/omniserp/client/exa`
- **Environment Variable**: `EXA_API_KEY`
- **Website**: [exa.ai](https://exa.ai)
- **Supported Operations**: Web Search (multiple modes)

Exa.ai provides neural search optimized for LLM applications with multiple search modes:

| Mode | Description | Use Case |
|------|-------------|----------|
| `auto` | Automatic mode selection | General purpose |
| `instant` | Fastest results | Real-time applications |
| `fast` | Balanced speed/quality | Most queries |
| `deep` | High-quality results | Research tasks |
| `deep-lite` | Deep with lower latency | Time-sensitive research |
| `deep-reasoning` | Best quality | Complex queries |

```go
import "github.com/plexusone/omniserp/client/exa"

engine := exa.New(exa.Config{
    APIKey: os.Getenv("EXA_API_KEY"),
})

// Neural search with deep mode
result, err := engine.Search(ctx, client.SearchParams{
    Query:      "latest AI research papers 2026",
    NumResults: 10,
    Options: map[string]any{
        "mode": "deep",
    },
})
```

#### Exa Features

- **Contents Extraction** - Get full page content, not just snippets
- **Highlights** - Key passages relevant to query
- **Autoprompt** - Automatic query optimization
- **Date Filtering** - Filter by publication date

## Feature Comparison

| Operation | Serper | SerpAPI | Brave | Exa |
|-----------|:------:|:-------:|:-----:|:---:|
| Web Search | ✓ | ✓ | ✓ | ✓ |
| News Search | ✓ | ✓ | ✓ | ✗ |
| Image Search | ✓ | ✓ | ✓ | ✗ |
| Video Search | ✓ | ✓ | ✓ | ✗ |
| Places Search | ✓ | ✓ | ✗ | ✗ |
| Maps Search | ✓ | ✓ | ✗ | ✗ |
| Reviews Search | ✓ | ✓ | ✗ | ✗ |
| Shopping Search | ✓ | ✓ | ✗ | ✗ |
| Scholar Search | ✓ | ✓ | ✗ | ✗ |
| Lens Search | ✓ | ✗ | ✗ | ✗ |
| Autocomplete | ✓ | ✓ | ✓ | ✗ |
| Webpage Scrape | ✓ | ✓ | ✗ | ✗ |
| **Neural Search** | ✗ | ✗ | ✗ | **✓** |
| **Content Extract** | ✗ | ✗ | ✗ | **✓** |
| **Summarizer** | ✗ | ✗ | **✓** | ✗ |

## Engine Interface

All engines implement the `Engine` interface:

```go
type Engine interface {
    // Metadata
    GetName() string
    GetVersion() string
    GetSupportedTools() []string

    // Search methods
    Search(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchNews(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchImages(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchVideos(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchPlaces(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchMaps(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchReviews(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchShopping(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchScholar(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchLens(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchAutocomplete(ctx context.Context, params SearchParams) (*SearchResult, error)

    // Utility
    ScrapeWebpage(ctx context.Context, params ScrapeParams) (*SearchResult, error)
}
```

## Selecting an Engine

### Via Environment Variable

```bash
export SEARCH_ENGINE="serper"  # or "serpapi", "brave", "exa"
```

### Programmatically

```go
// Use default (from SEARCH_ENGINE env var)
c, err := client.New()

// Specify explicitly
c, err := client.NewWithEngine("brave")

// Switch at runtime
c.SetEngine("exa")
```

### Engine Selection Guide

| Engine | Best For | Pricing |
|--------|----------|---------|
| **Serper** | Full-featured SERP data | Pay per search |
| **SerpAPI** | Google SERP scraping | Pay per search |
| **Brave** | Privacy-focused, fast | Free tier + paid |
| **Exa** | LLM applications, research | Pay per search |

## Checking Engine Capabilities

```go
c, _ := client.New()

// Check current engine
fmt.Printf("Engine: %s v%s\n", c.GetName(), c.GetVersion())

// List supported tools
tools := c.GetSupportedTools()
fmt.Printf("Supported: %v\n", tools)

// Check specific operation
if c.SupportsOperation(client.OpSearchLens) {
    // Lens is supported
}
```
