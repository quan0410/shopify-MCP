package tools

import (
	"github.com/datumbridge/shopify-mcp/internal/mcp"
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
		descs = append(descs, mcp.ToolDesc{
			Name:        name,
			Description: desc,
			InputSchema: schema(props, required),
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

	return descs, handlers
}
