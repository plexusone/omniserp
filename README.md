# OmniSerp Multi-Search Client and MCP Server

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omniserp/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omniserp/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omniserp/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omniserp/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omniserp/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omniserp/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omniserp
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omniserp
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omniserp
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omniserp
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomniserp
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omniserp
 [repo-url]: https://github.com/plexusone/omniserp
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omniserp/blob/master/LICENSE

A modular, plugin-based search engine abstraction package for Go that provides a unified interface for multiple search engines.

## Overview

The `omniserp` package provides:

- 📦 **Unified Client SDK**: Single API that fronts multiple search engine backends (`client/client.go`)
- 📐 **Normalized Responses**: Optional unified response structures across all engines (engine-agnostic)
- ✅ **Capability Checking**: Automatic validation of operation support across different backends
- 🔌 **Unified Interface**: Common `Engine` interface for all search providers
- 🧩 **Plugin Architecture**: Easy addition of new search engines
- 🤝 **Multiple Providers**: Built-in support for Serper and SerpAPI
- 🔒 **Type Safety**: Structured parameter and result types
- 📋 **Registry System**: Automatic discovery and management of engines
- 🤖 **MCP Server**: Model Context Protocol server for AI integration with optional secure credentials (`cmd/mcp-omniserp`)
- ⌨️ **CLI Tool**: Command-line interface for quick searches (`cmd/omniserp`)

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client"
)

func main() {
    // Set API key
    // export SERPER_API_KEY="your_key"

    // Create client (auto-selects engine)
    c, err := client.New()
    if err != nil {
        log.Fatal(err)
    }

    // Perform a search
    result, err := c.Search(context.Background(), omniserp.SearchParams{
        Query: "golang programming",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Results: %+v\n", result.Data)
}
```

## Project Structure

```
omniserp/
├── client/                 # Search engine client implementations
│   ├── client.go           # Unified Client SDK with capability checking
│   ├── serper/             # Serper.dev implementation
│   └── serpapi/            # SerpAPI implementation
├── cmd/                    # Executable applications
│   ├── mcp-omniserp/       # MCP server for AI integration (with optional secure credentials)
│   └── omniserp/           # CLI tool
├── examples/               # Example programs
│   ├── capability_check/   # Capability checking demo
│   └── normalized_search/  # Normalized responses demo
├── types.go                # Core types and Engine interface
├── normalized.go           # Normalized response types
├── normalizer.go           # Response normalizer
├── omniserp.go             # Utility functions
└── README.md
```

## Applications

### CLI Tool

#### Installation
```bash
go build ./cmd/omniserp
```

#### Basic Usage
```bash
# Set API key
export SERPER_API_KEY="your_api_key"

# Basic search (specify engine and query)
./omniserp -e serper -q "golang programming"

# Or use long flags
./omniserp --engine serpapi --query "golang programming"

# With SerpAPI
export SERPAPI_API_KEY="your_api_key"
./omniserp -e serpapi -q "golang programming"
```

### MCP Server

The Model Context Protocol (MCP) server enables AI assistants to perform web searches through this package.

#### Installation
```bash
go install github.com/plexusone/omniserp/cmd/mcp-omniserp@latest
```

Or build from source:
```bash
go build ./cmd/mcp-omniserp
```

#### Configuration

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "omniserp": {
      "command": "mcp-omniserp",
      "env": {
        "SERPER_API_KEY": "your_serper_api_key",
        "SEARCH_ENGINE": "serper"
      }
    }
  }
}
```

#### Features

The MCP server **dynamically registers only the tools supported by the current search engine backend**. This means:

- When using **Serper**, all 12 tools are available including Lens search
- When using **SerpAPI**, 11 tools are available (Lens is excluded)

Available tool categories:
- **Web Search**: General web searches with customizable parameters
- **News Search**: Search news articles
- **Image Search**: Search for images
- **Video Search**: Search for videos
- **Places Search**: Search for locations and businesses
- **Maps Search**: Search maps data
- **Reviews Search**: Search reviews
- **Shopping Search**: Search shopping/product listings
- **Scholar Search**: Search academic papers
- **Lens Search**: Visual search capabilities (Serper only)
- **Autocomplete**: Get search suggestions
- **Webpage Scrape**: Extract content from webpages

All searches support parameters like location, language, country, and number of results.

**Server Logs**: The MCP server logs which tools were registered and which were skipped:
```
2025/12/13 19:00:00 Using engine: serpapi v1.0.0
2025/12/13 19:00:00 Registered 11 tools: [google_search, google_search_news, ...]
2025/12/13 19:00:00 Skipped 1 unsupported tools: [google_search_lens]
```

#### Secure Mode (Optional)

The MCP server supports optional secure credential management using VaultGuard. When a policy file exists, API keys are retrieved from the OS keychain instead of environment variables.

**Setup for Secure Mode:**

1. Store your API key in the keychain:
   ```bash
   security add-generic-password -s "omnivault" -a "SERPER_API_KEY" -w "your-key"
   ```

2. Create a security policy (`~/.vaultguard/policy.json`):
   ```json
   {
     "version": 1,
     "local": {
       "require_encryption": true,
       "min_security_score": 50
     }
   }
   ```

3. Update your Claude Desktop config (no `env` section needed):
   ```json
   {
     "mcpServers": {
       "omniserp": {
         "command": "mcp-omniserp"
       }
     }
   }
   ```

**Note**: Without a policy file, the server works exactly as before using environment variables.

## Client SDK

The `client` package provides a high-level SDK that simplifies working with multiple search engines:

### Key Features

- **Auto-registration**: Automatically discovers and registers all available engines
- **Smart selection**: Uses `SEARCH_ENGINE` environment variable or defaults to Serper
- **Runtime switching**: Switch between engines without recreating the client
- **Capability checking**: Validates operations before calling backends
- **Error handling**: Returns `ErrOperationNotSupported` for unsupported operations
- **Clean API**: Implements the same `Engine` interface, proxying to the selected backend

### Quick Start

```go
import "github.com/plexusone/omniserp/client"

// Create client - auto-selects engine based on SEARCH_ENGINE env var
c, err := client.New()

// Or specify engine explicitly
c, err := client.NewWithEngine("serper")

// Check support before calling
if c.SupportsOperation(client.OpSearchLens) {
    result, _ := c.SearchLens(ctx, params)
}

// Switch engines at runtime
c.SetEngine("serpapi")
```

### Operation Constants

The SDK provides constants for all operations:
- `client.OpSearch` - Web search
- `client.OpSearchNews` - News search
- `client.OpSearchImages` - Image search
- `client.OpSearchVideos` - Video search
- `client.OpSearchPlaces` - Places search
- `client.OpSearchMaps` - Maps search
- `client.OpSearchReviews` - Reviews search
- `client.OpSearchShopping` - Shopping search
- `client.OpSearchScholar` - Scholar search
- `client.OpSearchLens` - Lens search (Serper only)
- `client.OpSearchAutocomplete` - Autocomplete
- `client.OpScrapeWebpage` - Webpage scraping

### Normalized Responses

The client SDK provides **optional normalized response methods** that return unified structures across all search engines:

```go
// Use *Normalized() methods for engine-agnostic response structures
normalized, err := c.SearchNormalized(ctx, params)

// Access results in a consistent format regardless of engine
for _, result := range normalized.OrganicResults {
    fmt.Printf("%s: %s\n", result.Title, result.Link)
}

// Switch engines without changing your code!
c.SetEngine("serpapi")
normalized, err = c.SearchNormalized(ctx, params) // Same structure!
```

**Available Normalized Methods:**
- `SearchNormalized()` - Web search with normalized results
- `SearchNewsNormalized()` - News search with normalized results
- `SearchImagesNormalized()` - Image search with normalized results

**Benefits:**
- **Engine-Agnostic**: Same code works with any backend
- **Type-Safe**: Strongly-typed result structures
- **Optional**: Raw responses still available via standard methods
- **Complete**: Preserves original response in `Raw` field

**Example Normalized Structure:**
```go
type NormalizedSearchResult struct {
    OrganicResults  []OrganicResult    // Standard search results
    AnswerBox       *AnswerBox         // Featured answer
    KnowledgeGraph  *KnowledgeGraph    // Knowledge panel
    RelatedSearches []RelatedSearch    // Related queries
    PeopleAlsoAsk   []PeopleAlsoAsk   // PAA questions
    NewsResults     []NewsResult       // News articles
    ImageResults    []ImageResult      // Images
    SearchMetadata  SearchMetadata     // Search info
    Raw             *SearchResult      // Original response
}
```

**Comparison: Raw vs Normalized**

| Aspect | Raw Response | Normalized Response |
|--------|--------------|---------------------|
| Field names | Engine-specific | Unified |
| Structure | Varies by engine | Consistent |
| Engine switching | Requires code changes | No changes needed |
| Type safety | `interface{}` | Strongly typed |
| Use case | Engine-specific features | Engine-agnostic apps |

## Library Usage

### Basic Usage with Client SDK

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client"
)

func main() {
    // Create client (auto-registers all engines and selects based on SEARCH_ENGINE env var)
    c, err := client.New()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    log.Printf("Using engine: %s v%s", c.GetName(), c.GetVersion())

    // Perform a search
    result, err := c.Search(context.Background(), omniserp.SearchParams{
        Query:      "golang programming",
        NumResults: 10,
        Language:   "en",
        Country:    "us",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Use the result
    log.Printf("Search completed: %+v", result.Data)
}
```

### Selecting a Specific Engine

```go
// Create client with a specific engine
c, err := client.NewWithEngine("serpapi")
if err != nil {
    log.Fatal(err)
}

// Or switch engines at runtime
c.SetEngine("serper")
```

### Capability Checking

The client SDK automatically checks if operations are supported by the current backend:

```go
c, _ := client.New()

// Check if an operation is supported
if c.SupportsOperation(client.OpSearchLens) {
    result, err := c.SearchLens(ctx, params)
    // ...
} else {
    log.Println("Current engine doesn't support Lens search")
}

// Or let the client return an error
result, err := c.SearchLens(ctx, params)
if errors.Is(err, client.ErrOperationNotSupported) {
    log.Printf("Operation not supported: %v", err)
}
```

### Advanced Usage with Registry

For direct registry access:

```go
import (
    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client/serper"
    "github.com/plexusone/omniserp/client/serpapi"
)

func main() {
    // Create registry and manually register engines
    registry := omniserp.NewRegistry()

    // Register engines (handle errors as needed)
    if serperEngine, err := serper.New(); err == nil {
        registry.Register(serperEngine)
    }
    if serpApiEngine, err := serpapi.New(); err == nil {
        registry.Register(serpApiEngine)
    }

    // Get default engine (based on SEARCH_ENGINE env var)
    engine, err := omniserp.GetDefaultEngine(registry)
    if err != nil {
        log.Printf("Warning: %v", err)
    }

    // Perform a search
    result, err := engine.Search(context.Background(), omniserp.SearchParams{
        Query: "golang programming",
    })
    // ...
}
```

## Supported Engines

### Serper
- **Package**: `github.com/plexusone/omniserp/client/serper`
- **Environment Variable**: `SERPER_API_KEY`
- **Website**: [serper.dev](https://serper.dev)
- **Supported Operations**: All search types including Lens

### SerpAPI
- **Package**: `github.com/plexusone/omniserp/client/serpapi`
- **Environment Variable**: `SERPAPI_API_KEY`
- **Website**: [serpapi.com](https://serpapi.com)
- **Supported Operations**: All search types except Lens
- **Note**: `SearchLens()` is not supported and will return `ErrOperationNotSupported`

| Operation | Serper | SerpAPI |
|-----------|--------|---------|
| Web Search | ✓ | ✓ |
| News Search | ✓ | ✓ |
| Image Search | ✓ | ✓ |
| Video Search | ✓ | ✓ |
| Places Search | ✓ | ✓ |
| Maps Search | ✓ | ✓ |
| Reviews Search | ✓ | ✓ |
| Shopping Search | ✓ | ✓ |
| Scholar Search | ✓ | ✓ |
| **Lens Search** | **✓** | **✗** |
| Autocomplete | ✓ | ✓ |
| Webpage Scrape | ✓ | ✓ |

## Available Search Methods

All engines implement these methods:

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

## Types

### SearchParams
```go
type SearchParams struct {
    Query      string `json:"query"`                    // Required: search query
    Location   string `json:"location,omitempty"`       // Optional: search location
    Language   string `json:"language,omitempty"`       // Optional: language code (e.g., "en")
    Country    string `json:"country,omitempty"`        // Optional: country code (e.g., "us")
    NumResults int    `json:"num_results,omitempty"`    // Optional: number of results (1-100)
}
```

### ScrapeParams
```go
type ScrapeParams struct {
    URL string `json:"url"` // Required: URL to scrape
}
```

### SearchResult
```go
type SearchResult struct {
    Data interface{} `json:"data"`          // Parsed response data
    Raw  string      `json:"raw,omitempty"` // Raw response (optional)
}
```

## Registry Usage

### Basic Registry Operations
```go
import (
    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client/serper"
)

// Create new registry and register engines
registry := omniserp.NewRegistry()

// Register engines manually
if serperEngine, err := serper.New(); err == nil {
    registry.Register(serperEngine)
}

// List available engines
engines := registry.List()
log.Printf("Available engines: %v", engines)

// Get specific engine
if engine, exists := registry.Get("serper"); exists {
    log.Printf("Using engine: %s v%s", engine.GetName(), engine.GetVersion())
}

// Get all engines
allEngines := registry.GetAll()
```

### Engine Information
```go
// Get info about specific engine
engine, _ := registry.Get("serper")
info := omniserp.GetEngineInfo(engine)
log.Printf("Engine: %s v%s, Tools: %v", info.Name, info.Version, info.SupportedTools)

// Get info about all engines
allInfo := omniserp.GetAllEngineInfo(registry)
```

## Environment Configuration

The package uses environment variables for configuration:

```bash
# Choose which engine to use (optional, defaults to "serper")
export SEARCH_ENGINE="serper"  # or "serpapi"

# API keys for respective engines
export SERPER_API_KEY="your_serper_key"
export SERPAPI_API_KEY="your_serpapi_key"
```

## Adding New Engines

To add a new search engine:

1. **Create engine package under `client/`**:
```go
// client/newengine/newengine.go
package newengine

import (
    "context"
    "fmt"
    "os"
    "github.com/plexusone/omniserp"
)

type Engine struct {
    apiKey string
    // other fields
}

func New() (*Engine, error) {
    apiKey := os.Getenv("NEWENGINE_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("NEWENGINE_API_KEY required")
    }
    return &Engine{apiKey: apiKey}, nil
}

func (e *Engine) GetName() string { return "newengine" }
func (e *Engine) GetVersion() string { return "1.0.0" }
func (e *Engine) GetSupportedTools() []string { /* return supported tools */ }

// Implement all other omniserp.Engine methods...
func (e *Engine) Search(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
    // Implementation
}
// ... implement all other interface methods
```

2. **Register in your application**:
```go
// In your application code (e.g., cmd/yourapp/main.go)
import (
    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client/newengine"
    "github.com/plexusone/omniserp/client/serper"
)

func createRegistry() *omniserp.Registry {
    registry := omniserp.NewRegistry()

    // Register existing engines
    if serperEngine, err := serper.New(); err == nil {
        registry.Register(serperEngine)
    }

    // Register new engine
    if newEng, err := newengine.New(); err == nil {
        registry.Register(newEng)
    }

    return registry
}
```

3. **Update CLI (optional)**:
   Add the new engine import and registration to `cmd/omniserp/main.go`

## Error Handling

The package provides consistent error handling:

```go
engine, err := omniserp.GetDefaultEngine(registry)
if err != nil {
    // Handle engine selection error
    log.Printf("Engine selection warning: %v", err)
}

result, err := engine.Search(ctx, params)
if err != nil {
    // Handle search error
    log.Printf("Search failed: %v", err)
}
```

## Examples

See the `examples/` directory for working examples:

- **`capability_check/`**: Demonstrates capability checking, engine switching, and operation support matrix
- **`normalized_search/`**: Shows normalized responses and engine-agnostic code

To run an example:
```bash
export SERPER_API_KEY="your_key"
export SERPAPI_API_KEY="your_key"  # optional

# Check capabilities
go run examples/capability_check/main.go

# Demonstrate normalized responses
go run examples/normalized_search/main.go "golang programming"
```

## Testing

Run tests without API keys (tests will skip gracefully):
```bash
go test ./...
```

Run tests with API calls (requires API keys):
```bash
export SERPER_API_KEY="your_key"
export SERPAPI_API_KEY="your_key"
go test -v ./client
```

## Thread Safety

The registry is safe for concurrent read operations. Engine implementations should be thread-safe for concurrent use.
