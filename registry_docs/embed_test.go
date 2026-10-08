package registrydocs

import (
	"strings"
	"testing"
)

var allShopifyTools = []string{
	// Products & Variants (8)
	"shopify_list_products", "shopify_get_product", "shopify_create_product",
	"shopify_update_product", "shopify_delete_product", "shopify_list_variants",
	"shopify_get_variant", "shopify_update_variant",
	// Orders & Draft Orders (7)
	"shopify_list_orders", "shopify_get_order", "shopify_create_order",
	"shopify_cancel_order", "shopify_close_order", "shopify_list_draft_orders",
	"shopify_create_draft_order",
	// Customers (5)
	"shopify_list_customers", "shopify_get_customer", "shopify_search_customers",
	"shopify_create_customer", "shopify_update_customer",
	// Inventory & Locations (4)
	"shopify_list_locations", "shopify_list_inventory_levels",
	"shopify_adjust_inventory_level", "shopify_set_inventory_level",
	// Collections (4)
	"shopify_list_custom_collections", "shopify_get_custom_collection",
	"shopify_list_smart_collections", "shopify_create_custom_collection",
	// Shop (1)
	"shopify_get_shop",
	// Fulfillments (4)
	"shopify_list_fulfillments", "shopify_get_fulfillment",
	"shopify_create_fulfillment", "shopify_cancel_fulfillment",
	// Discounts & Price Rules (7)
	"shopify_list_price_rules", "shopify_get_price_rule", "shopify_create_price_rule",
	"shopify_delete_price_rule", "shopify_list_discount_codes", "shopify_create_discount_code",
	"shopify_lookup_discount_code",
	// Content: Pages & Blogs/Articles (10)
	"shopify_list_pages", "shopify_get_page", "shopify_create_page",
	"shopify_update_page", "shopify_delete_page", "shopify_list_blogs",
	"shopify_list_articles", "shopify_get_article", "shopify_create_article",
	"shopify_update_article",
	// GraphQL (1)
	"shopify_graphql",
}

func TestMarkdown_AllToolsExist(t *testing.T) {
	for _, name := range allShopifyTools {
		doc := Markdown(name)
		if doc == "" {
			t.Errorf("missing documentation file for tool: %s", name)
		}
		if !strings.HasPrefix(doc, "# "+name) {
			t.Errorf("doc for %s should start with header '# %s'", name, name)
		}
		if !strings.Contains(doc, "## Parameters") {
			t.Errorf("doc for %s missing '## Parameters' section", name)
		}
		if !strings.Contains(doc, "## Cases") {
			t.Errorf("doc for %s missing '## Cases' section", name)
		}
	}
}

func TestMarkdown_SecurityAndErrors(t *testing.T) {
	if Markdown("no_such_tool") != "" {
		t.Fatal("missing tool should return empty string")
	}
	if Markdown("../embed") != "" {
		t.Fatal("path traversal should return empty string")
	}
	if Markdown("sub/tool") != "" {
		t.Fatal("subpath should return empty string")
	}
	if Markdown("") != "" {
		t.Fatal("empty tool name should return empty string")
	}
}
