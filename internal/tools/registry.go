package tools

import (
	"strings"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
	registrydocs "github.com/datumbridge/shopify-mcp/registry_docs"
)

var toolOutputs = map[string]string{
	// Products & Variants (8)
	"shopify_list_products":  "products",
	"shopify_get_product":    "product, id, title, variants",
	"shopify_create_product": "product, id, title, handle, variants",
	"shopify_update_product": "product, id, title, handle, variants",
	"shopify_delete_product": "",
	"shopify_list_variants":  "variants",
	"shopify_get_variant":    "variant, id, product_id, title, price, sku",
	"shopify_update_variant": "variant, id, product_id, title, price, sku",

	// Orders & Draft Orders (7)
	"shopify_list_orders":        "orders",
	"shopify_get_order":          "order, id, name, financial_status, total_price, line_items",
	"shopify_create_order":       "order, id, name, total_price, line_items",
	"shopify_cancel_order":       "order, id, cancelled_at, cancel_reason",
	"shopify_close_order":        "order, id, closed_at",
	"shopify_list_draft_orders":  "draft_orders",
	"shopify_create_draft_order": "draft_order, id, name, status, total_price",

	// Customers (5)
	"shopify_list_customers":   "customers",
	"shopify_get_customer":     "customer, id, email, first_name, last_name, phone",
	"shopify_search_customers": "customers",
	"shopify_create_customer":  "customer, id, email, first_name, last_name",
	"shopify_update_customer":  "customer, id, email, first_name, last_name",

	// Inventory & Locations (4)
	"shopify_list_locations":         "locations",
	"shopify_list_inventory_levels":  "inventory_levels",
	"shopify_adjust_inventory_level": "inventory_level, inventory_item_id, location_id, available",
	"shopify_set_inventory_level":    "inventory_level, inventory_item_id, location_id, available",

	// Collections (4)
	"shopify_list_custom_collections":  "custom_collections",
	"shopify_get_custom_collection":    "custom_collection, id, title, handle",
	"shopify_list_smart_collections":   "smart_collections",
	"shopify_create_custom_collection": "custom_collection, id, title, handle",

	// Shop (1)
	"shopify_get_shop": "shop, id, name, email, domain, currency",

	// Fulfillments (4)
	"shopify_list_fulfillments":  "fulfillments",
	"shopify_get_fulfillment":    "fulfillment, id, order_id, status, tracking_company, tracking_number",
	"shopify_create_fulfillment": "fulfillment, id, order_id, status",
	"shopify_cancel_fulfillment": "fulfillment, id, order_id, status",

	// Discounts & Price Rules (7)
	"shopify_list_price_rules":     "price_rules",
	"shopify_get_price_rule":       "price_rule, id, title, value_type, value",
	"shopify_create_price_rule":    "price_rule, id, title, value_type, value",
	"shopify_delete_price_rule":    "",
	"shopify_list_discount_codes":  "discount_codes",
	"shopify_create_discount_code": "discount_code, id, price_rule_id, code",
	"shopify_lookup_discount_code": "discount_code, id, price_rule_id, code",

	// Content: Pages & Blogs/Articles (10)
	"shopify_list_pages":     "pages",
	"shopify_get_page":       "page, id, title, handle, body_html",
	"shopify_create_page":    "page, id, title, handle",
	"shopify_update_page":    "page, id, title, handle",
	"shopify_delete_page":    "",
	"shopify_list_blogs":     "blogs",
	"shopify_list_articles":  "articles",
	"shopify_get_article":    "article, id, blog_id, title, body_html",
	"shopify_create_article": "article, id, blog_id, title",
	"shopify_update_article": "article, id, blog_id, title",

	// GraphQL (1)
	"shopify_graphql": "data",
}

type relatedMeta struct {
	tool  string
	field string
}

var paramRelated = map[string]relatedMeta{
	"product_id":         {tool: "shopify_list_products", field: "products[0].id"},
	"order_id":           {tool: "shopify_list_orders", field: "orders[0].id"},
	"customer_id":        {tool: "shopify_list_customers", field: "customers[0].id"},
	"variant_id":         {tool: "shopify_list_variants", field: "variants[0].id"},
	"location_id":        {tool: "shopify_list_locations", field: "locations[0].id"},
	"location_ids":       {tool: "shopify_list_locations", field: "locations[0].id"},
	"inventory_item_id":  {tool: "shopify_get_variant", field: "variant.inventory_item_id"},
	"inventory_item_ids": {tool: "shopify_get_variant", field: "variant.inventory_item_id"},
	"collection_id":      {tool: "shopify_list_custom_collections", field: "custom_collections[0].id"},
	"price_rule_id":      {tool: "shopify_list_price_rules", field: "price_rules[0].id"},
	"page_id":            {tool: "shopify_list_pages", field: "pages[0].id"},
	"blog_id":            {tool: "shopify_list_blogs", field: "blogs[0].id"},
	"article_id":         {tool: "shopify_list_articles", field: "articles[0].id"},
	"fulfillment_id":     {tool: "shopify_list_fulfillments", field: "fulfillments[0].id"},
}

func annotateRelated(toolName string, props map[string]interface{}) {
	for k, v := range props {
		propMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if rel, found := paramRelated[k]; found {
			propMap["x-datumbridge-related-tool"] = rel.tool
			propMap["x-datumbridge-output-field"] = rel.field
		}
	}
}

func formatToolDescription(name, desc string) string {
	res := strings.TrimSpace(desc)
	if !strings.Contains(strings.ToLower(res), "capabilities:") {
		res += "\n\nCapabilities: shopify." + name
	}
	if !strings.Contains(strings.ToLower(res), "outputs:") {
		outs := toolOutputs[name]
		if outs != "" {
			res += "\nOutputs: " + outs + ", success"
		} else {
			res += "\nOutputs: success"
		}
	}
	return res
}

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
		annotateRelated(name, props)
		input := schema(props, required)
		input["x-datumbridge-capabilities"] = []string{"shopify." + name}
		if doc := registrydocs.Markdown(name); doc != "" {
			input["x-datumbridge-docs"] = doc
		}
		descs = append(descs, mcp.ToolDesc{
			Name:        name,
			Description: formatToolDescription(name, desc),
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
