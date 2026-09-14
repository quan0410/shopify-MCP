package tools

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerShopTools(add toolAdder) {
	add(
		"shopify_get_shop",
		"Retrieve configuration and metadata for the authenticated Shopify store (name, domain, currency, timezone, email).",
		baseProps(nil),
		nil,
		handleShopifyGetShop,
	)
}

func handleShopifyGetShop(args json.RawMessage) map[string]interface{} {
	client, _, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/shop.json", url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
