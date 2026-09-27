package tools

import (
	"encoding/json"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerGraphQLTools(add toolAdder) {
	add(
		"shopify_graphql",
		"Execute an arbitrary GraphQL query or mutation against Shopify Admin API (2024-04). Recommended for fetching orders or customers with specific non-PII fields to bypass Protected Customer Data restrictions.",
		baseProps(map[string]interface{}{
			"query":     map[string]interface{}{"type": "string", "description": "The GraphQL query or mutation string"},
			"variables": map[string]interface{}{"type": "object", "description": "Optional variables map for the GraphQL query"},
		}),
		[]string{"query"},
		handleShopifyGraphQL,
	)
}

func handleShopifyGraphQL(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	query := strArg(m, "query")
	if query == "" {
		return mcp.ToolResultError("query is required")
	}
	vars := mapArg(m, "variables")

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.GraphQL(c, query, vars)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
