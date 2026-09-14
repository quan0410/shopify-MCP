package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerCollectionTools(add toolAdder) {
	add(
		"shopify_list_custom_collections",
		"List custom (manually-curated) product collections from Shopify.",
		baseProps(map[string]interface{}{
			"limit":    map[string]interface{}{"type": "integer", "description": "Number of collections to return (default 50)"},
			"since_id": map[string]interface{}{"type": "string", "description": "Restrict results to after specified collection ID"},
			"title":    map[string]interface{}{"type": "string", "description": "Filter collections by title"},
		}),
		nil,
		handleShopifyListCustomCollections,
	)

	add(
		"shopify_get_custom_collection",
		"Retrieve a specific custom collection by its collection ID.",
		baseProps(map[string]interface{}{
			"collection_id": map[string]interface{}{"type": "string", "description": "Collection ID"},
		}),
		[]string{"collection_id"},
		handleShopifyGetCustomCollection,
	)

	add(
		"shopify_list_smart_collections",
		"List automated (smart) product collections with rule conditions.",
		baseProps(map[string]interface{}{
			"limit": map[string]interface{}{"type": "integer", "description": "Number of collections to return (default 50)"},
			"title": map[string]interface{}{"type": "string", "description": "Filter smart collections by title"},
		}),
		nil,
		handleShopifyListSmartCollections,
	)

	add(
		"shopify_create_custom_collection",
		"Create a new custom collection with title, description, and optional image.",
		baseProps(map[string]interface{}{
			"title":     map[string]interface{}{"type": "string", "description": "Collection title"},
			"body_html": map[string]interface{}{"type": "string", "description": "Collection description"},
			"published": map[string]interface{}{"type": "boolean", "description": "Whether collection is published (default true)"},
		}),
		[]string{"title"},
		handleShopifyCreateCustomCollection,
	)
}

func handleShopifyListCustomCollections(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	q := url.Values{}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}
	if s := strArg(m, "since_id"); s != "" {
		q.Set("since_id", s)
	}
	if t := strArg(m, "title"); t != "" {
		q.Set("title", t)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/custom_collections.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetCustomCollection(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	colID := strArg(m, "collection_id")
	if colID == "" {
		return mcp.ToolResultError("collection_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/custom_collections/%s.json", colID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListSmartCollections(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	q := url.Values{}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}
	if t := strArg(m, "title"); t != "" {
		q.Set("title", t)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/smart_collections.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateCustomCollection(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	title := strArg(m, "title")
	if title == "" {
		return mcp.ToolResultError("title is required")
	}

	col := map[string]interface{}{
		"title": title,
	}
	if desc := strArg(m, "body_html"); desc != "" {
		col["body_html"] = desc
	}
	if pub := boolArg(m, "published", true); !pub {
		col["published"] = false
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/custom_collections.json", url.Values{}, map[string]interface{}{"custom_collection": col})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
