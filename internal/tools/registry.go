package tools

import (
	"github.com/datumbridge/shopify-mcp/internal/mcp"
	registrydocs "github.com/datumbridge/shopify-mcp/registry_docs"
)

// Register returns the full shopify-mcp tool catalogue: products, orders,
// customers, inventory, collections, and shop metadata. Full scope is the
// default; SHOPIFY_TOOLS narrows it to an explicit allowlist if set.
func Register() ([]mcp.ToolDesc, map[string]mcp.ToolHandler) {
	cfg := LoadServerConfig()
	handlers := map[string]mcp.ToolHandler{}
	var descs []mcp.ToolDesc

	add := func(name, desc string, props map[string]interface{}, required []string, h mcp.ToolHandler) {
		if !cfg.toolEnabled(name) {
			return
		}
		input := schema(props, required)
		if doc := registrydocs.Markdown(name); doc != "" {
			input["x-datumbridge-docs"] = doc
		}
		descs = append(descs, mcp.ToolDesc{
			Name:        name,
			Description: desc,
			InputSchema: input,
		})
		handlers[name] = h
	}

	registerProductTools(add)
	registerOrderTools(add)
	registerCustomerTools(add)
	registerInventoryTools(add)
	registerCollectionTools(add)
	registerShopTools(add)
	registerFulfillmentTools(add)
	registerDiscountTools(add)
	registerContentTools(add)
	registerGraphQLTools(add)

	return descs, handlers
}
