// mcp-omniserp is an MCP server for web search with multi-engine support and optional secure credentials.
//
// This server supports multiple search engine backends (Serper, SerpAPI) and can optionally
// use VaultGuard for secure credential management:
//
// Standard Mode (environment variables):
//
//	export SERPER_API_KEY="your-key"    # or SERPAPI_API_KEY
//	export SEARCH_ENGINE="serper"       # optional, defaults to serper
//	./mcp-omniserp
//
// Secure Mode (OS keychain + policy):
//
//  1. Store your API key in the keychain:
//     security add-generic-password -s "omnivault" -a "SERPER_API_KEY" -w "your-key"
//
//  2. Create ~/.vaultguard/policy.json:
//     {
//     "version": 1,
//     "local": {
//     "require_encryption": true,
//     "min_security_score": 50
//     }
//     }
//
//  3. Run the server:
//     ./mcp-omniserp
//
// When a policy file exists, credentials are retrieved from the OS keychain with security
// posture validation. Without a policy file, standard environment variables are used.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	keyring "github.com/plexusone/omni-keyring/omnivault"
	"github.com/plexusone/vaultguard"

	"github.com/plexusone/omniserp"
	"github.com/plexusone/omniserp/client"
	"github.com/plexusone/omniserp/client/serpapi"
	"github.com/plexusone/omniserp/client/serper"
)

// ToolDefinition defines a search tool with its metadata
type ToolDefinition struct {
	Name        string
	Description string
	SearchFunc  func(context.Context, omniserp.SearchParams) (*omniserp.SearchResult, error)
}

func main() {
	ctx := context.Background()

	// Load policy from config files (or nil for permissive mode)
	policy, err := vaultguard.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// Initialize search client based on credential mode
	var searchClient *client.Client
	if policy == nil {
		log.Println("No policy configured - using environment variables")
		searchClient, err = initWithEnvCredentials()
	} else {
		log.Println("Policy loaded - using secure credential access")
		searchClient, err = initWithSecureCredentials(ctx, policy)
	}
	if err != nil {
		log.Fatalf("Failed to initialize search client: %v", err)
	}

	runServer(ctx, searchClient)
}

// initWithEnvCredentials initializes the client using environment variables.
func initWithEnvCredentials() (*client.Client, error) {
	return client.New()
}

// initWithSecureCredentials initializes the client using VaultGuard and OS keychain.
func initWithSecureCredentials(ctx context.Context, policy *vaultguard.Policy) (*client.Client, error) {
	// Create keyring provider for OS credential store
	keyringVault := keyring.New(keyring.Config{
		ServiceName: "omnivault",
	})

	// Create VaultGuard with the keyring and loaded policy
	sv, err := vaultguard.New(&vaultguard.Config{
		CustomVault: keyringVault,
		Policy:      policy,
	})
	if err != nil {
		return nil, fmt.Errorf("security check failed: %w", err)
	}
	defer sv.Close()

	// Log security status
	result := sv.SecurityResult()
	if result != nil {
		log.Printf("Security check passed: score=%d, level=%s", result.Score, result.Level)
		if result.Details.Local != nil {
			log.Printf("  Platform: %s, Encrypted: %v, Biometrics: %v",
				result.Details.Local.Platform,
				result.Details.Local.DiskEncrypted,
				result.Details.Local.BiometricsConfigured)
		}
	}

	// Determine which engine to use
	engineName := os.Getenv("SEARCH_ENGINE")
	if engineName == "" {
		engineName = "serper"
	}

	// Create registry and register engines with secure credentials
	registry := omniserp.NewRegistry()

	switch engineName {
	case "serper":
		apiKey, err := sv.GetValue(ctx, "SERPER_API_KEY")
		if err != nil {
			return nil, fmt.Errorf("failed to get SERPER_API_KEY from keychain: %w", err)
		}
		if apiKey == "" {
			return nil, fmt.Errorf("SERPER_API_KEY not found in keychain. Add it with:\n" +
				"  security add-generic-password -s \"omnivault\" -a \"SERPER_API_KEY\" -w \"your-key\"")
		}
		log.Println("SERPER_API_KEY retrieved from keychain successfully")

		engine, err := serper.NewWithAPIKey(apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create serper engine: %w", err)
		}
		registry.Register(engine)

	case "serpapi":
		apiKey, err := sv.GetValue(ctx, "SERPAPI_API_KEY")
		if err != nil {
			return nil, fmt.Errorf("failed to get SERPAPI_API_KEY from keychain: %w", err)
		}
		if apiKey == "" {
			return nil, fmt.Errorf("SERPAPI_API_KEY not found in keychain. Add it with:\n" +
				"  security add-generic-password -s \"omnivault\" -a \"SERPAPI_API_KEY\" -w \"your-key\"")
		}
		log.Println("SERPAPI_API_KEY retrieved from keychain successfully")

		engine, err := serpapi.NewWithAPIKey(apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create serpapi engine: %w", err)
		}
		registry.Register(engine)

	default:
		return nil, fmt.Errorf("unsupported engine: %s", engineName)
	}

	return client.NewWithRegistry(registry, engineName)
}

// runServer starts the MCP server with the configured search client.
func runServer(ctx context.Context, searchClient *client.Client) {
	log.Printf("Using engine: %s v%s", searchClient.GetName(), searchClient.GetVersion())
	log.Printf("Available engines: %v", searchClient.ListEngines())

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-omniserp",
		Version: "2.0.0",
	}, nil)

	// Define all possible search tools with their operation names
	allTools := []ToolDefinition{
		{client.OpSearch, "Perform a Google web search", searchClient.Search},
		{client.OpSearchNews, "Search for news articles using Google News", searchClient.SearchNews},
		{client.OpSearchImages, "Search for images using Google Images", searchClient.SearchImages},
		{client.OpSearchVideos, "Search for videos using Google Videos", searchClient.SearchVideos},
		{client.OpSearchPlaces, "Search for places using Google Places", searchClient.SearchPlaces},
		{client.OpSearchMaps, "Search for locations using Google Maps", searchClient.SearchMaps},
		{client.OpSearchReviews, "Search for reviews", searchClient.SearchReviews},
		{client.OpSearchShopping, "Search for products using Google Shopping", searchClient.SearchShopping},
		{client.OpSearchScholar, "Search for academic papers using Google Scholar", searchClient.SearchScholar},
		{client.OpSearchLens, "Perform visual search using Google Lens", searchClient.SearchLens},
		{client.OpSearchAutocomplete, "Get search suggestions using Google Autocomplete", searchClient.SearchAutocomplete},
	}

	// Register tools only if supported by the current engine
	registeredTools := []string{}
	skippedTools := []string{}

	for _, tool := range allTools {
		if searchClient.SupportsOperation(tool.Name) {
			// Register this tool
			toolName := tool.Name
			toolDesc := tool.Description
			searchFunc := tool.SearchFunc

			mcp.AddTool(server, &mcp.Tool{
				Name:        toolName,
				Description: toolDesc,
			}, func(ctx context.Context, req *mcp.CallToolRequest, args omniserp.SearchParams) (*mcp.CallToolResult, any, error) {
				result, err := searchFunc(ctx, args)
				if err != nil {
					return nil, nil, fmt.Errorf("%s failed: %w", toolName, err)
				}

				resultJSON, _ := json.MarshalIndent(result.Data, "", "  ")
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: string(resultJSON)},
					},
				}, nil, nil
			})

			registeredTools = append(registeredTools, tool.Name)
		} else {
			skippedTools = append(skippedTools, tool.Name)
		}
	}

	// Register web scraping tool if supported
	if searchClient.SupportsOperation(client.OpScrapeWebpage) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        client.OpScrapeWebpage,
			Description: "Scrape content from a webpage",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args omniserp.ScrapeParams) (*mcp.CallToolResult, any, error) {
			result, err := searchClient.ScrapeWebpage(ctx, args)
			if err != nil {
				return nil, nil, fmt.Errorf("scraping failed: %w", err)
			}

			resultJSON, _ := json.MarshalIndent(result.Data, "", "  ")
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(resultJSON)},
				},
			}, nil, nil
		})
		registeredTools = append(registeredTools, client.OpScrapeWebpage)
	} else {
		skippedTools = append(skippedTools, client.OpScrapeWebpage)
	}

	// Log tool registration summary
	log.Printf("Registered %d tools: %v", len(registeredTools), registeredTools)
	if len(skippedTools) > 0 {
		log.Printf("Skipped %d unsupported tools: %v", len(skippedTools), skippedTools)
	}

	log.Printf("Starting OmniSerp MCP Server with %s engine...", searchClient.GetName())
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
