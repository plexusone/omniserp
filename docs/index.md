# OmniSerp

A modular, plugin-based search engine abstraction package for Go that provides a unified interface for multiple search engines.

## Overview

The `omniserp` package provides:

- **Unified Client SDK**: Single API that fronts multiple search engine backends
- **Normalized Responses**: Optional unified response structures across all engines (engine-agnostic)
- **Capability Checking**: Automatic validation of operation support across different backends
- **Unified Interface**: Common `Engine` interface for all search providers
- **Plugin Architecture**: Easy addition of new search engines
- **Multiple Providers**: Built-in support for Serper, SerpAPI, Brave Search, and Exa.ai
- **Type Safety**: Structured parameter and result types
- **Registry System**: Automatic discovery and management of engines
- **MCP Server**: Model Context Protocol server for AI integration with optional secure credentials
- **CLI Tool**: Command-line interface for quick searches

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
│   ├── serpapi/            # SerpAPI implementation
│   ├── brave/              # Brave Search API
│   └── exa/                # Exa.ai neural search
├── cmd/                    # Executable applications
│   ├── mcp-omniserp/       # MCP server for AI integration
│   └── omniserp/           # CLI tool
├── examples/               # Example programs
│   ├── capability_check/   # Capability checking demo
│   └── normalized_search/  # Normalized responses demo
├── types.go                # Core types and Engine interface
├── normalized.go           # Normalized response types
├── normalizer.go           # Response normalizer
└── omniserp.go             # Utility functions
```

## Installation

```bash
go get github.com/plexusone/omniserp@latest
```

## Environment Configuration

```bash
# Choose which engine to use (optional, defaults to "serper")
export SEARCH_ENGINE="serper"  # or "serpapi", "brave", "exa"

# API keys for respective engines
export SERPER_API_KEY="your_serper_key"
export SERPAPI_API_KEY="your_serpapi_key"
export BRAVE_API_KEY="your_brave_key"
export EXA_API_KEY="your_exa_key"
```

## Supported Engines

| Engine | Best For | Key Features |
|--------|----------|--------------|
| **Serper** | Full SERP data | All 12 search types, Lens support |
| **SerpAPI** | Google scraping | Reliable, well-documented |
| **Brave** | Privacy, speed | Free tier, summarizer, goggles |
| **Exa** | LLM/AI apps | Neural search, content extraction |

See [Engines Overview](engines/overview.md) for detailed feature comparison.
