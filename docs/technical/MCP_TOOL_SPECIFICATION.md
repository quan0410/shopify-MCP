# MCP Tool Specification — shopify-mcp

| Field | Value |
|---|---|
| Folder | `shopify-mcp` |
| mcpServer id | `shopify-mcp` |
| Transport | Streamable HTTP `POST /mcp` |
| Protocol | `2024-11-05` |
| Health | `GET /health` |
| Default port | `8013` |
| Upstream | Shopify Admin REST API (`2024-04`) |

## Tool visibility

Full scope by default (49 tools). `SHOPIFY_TOOLS=tool1,tool2` narrows to an explicit allowlist.

## Tools (49)

| Name | Group | Method | Endpoint |
|---|---|---|---|
| `shopify_list_products` | products | GET | `/products.json` |
| `shopify_get_product` | products | GET | `/products/{id}.json` |
| `shopify_create_product` | products | POST | `/products.json` |
| `shopify_update_product` | products | PUT | `/products/{id}.json` |
| `shopify_delete_product` | products | DELETE | `/products/{id}.json` |
| `shopify_list_variants` | products | GET | `/products/{id}/variants.json` |
| `shopify_get_variant` | products | GET | `/variants/{id}.json` |
| `shopify_update_variant` | products | PUT | `/variants/{id}.json` |
| `shopify_list_orders` | orders | GET | `/orders.json` |
| `shopify_get_order` | orders | GET | `/orders/{id}.json` |
| `shopify_create_order` | orders | POST | `/orders.json` |
| `shopify_cancel_order` | orders | POST | `/orders/{id}/cancel.json` |
| `shopify_close_order` | orders | POST | `/orders/{id}/close.json` |
| `shopify_list_draft_orders` | orders | GET | `/draft_orders.json` |
| `shopify_create_draft_order` | orders | POST | `/draft_orders.json` |
| `shopify_list_customers` | customers | GET | `/customers.json` |
| `shopify_get_customer` | customers | GET | `/customers/{id}.json` |
| `shopify_search_customers` | customers | GET | `/customers/search.json` |
| `shopify_create_customer` | customers | POST | `/customers.json` |
| `shopify_update_customer` | customers | PUT | `/customers/{id}.json` |
| `shopify_list_locations` | inventory | GET | `/locations.json` |
| `shopify_list_inventory_levels` | inventory | GET | `/inventory_levels.json` |
| `shopify_adjust_inventory_level` | inventory | POST | `/inventory_levels/adjust.json` |
| `shopify_set_inventory_level` | inventory | POST | `/inventory_levels/set.json` |
| `shopify_list_custom_collections` | collections | GET | `/custom_collections.json` |
| `shopify_get_custom_collection` | collections | GET | `/custom_collections/{id}.json` |
| `shopify_list_smart_collections` | collections | GET | `/smart_collections.json` |
| `shopify_create_custom_collection` | collections | POST | `/custom_collections.json` |
| `shopify_get_shop` | shop | GET | `/shop.json` |
| `shopify_list_fulfillments` | fulfillments | GET | `/orders/{order_id}/fulfillments.json` |
| `shopify_get_fulfillment` | fulfillments | GET | `/orders/{order_id}/fulfillments/{id}.json` |
| `shopify_create_fulfillment` | fulfillments | POST | `/fulfillments.json` |
| `shopify_cancel_fulfillment` | fulfillments | POST | `/fulfillments/{id}/cancel.json` |
| `shopify_list_price_rules` | discounts | GET | `/price_rules.json` |
| `shopify_get_price_rule` | discounts | GET | `/price_rules/{id}.json` |
| `shopify_create_price_rule` | discounts | POST | `/price_rules.json` |
| `shopify_delete_price_rule` | discounts | DELETE | `/price_rules/{id}.json` |
| `shopify_list_discount_codes` | discounts | GET | `/price_rules/{id}/discount_codes.json` |
| `shopify_create_discount_code` | discounts | POST | `/price_rules/{id}/discount_codes.json` |
| `shopify_lookup_discount_code` | discounts | GET | `/discount_codes/lookup.json` |
| `shopify_list_pages` | content | GET | `/pages.json` |
| `shopify_get_page` | content | GET | `/pages/{id}.json` |
| `shopify_create_page` | content | POST | `/pages.json` |
| `shopify_update_page` | content | PUT | `/pages/{id}.json` |
| `shopify_delete_page` | content | DELETE | `/pages/{id}.json` |
| `shopify_list_blogs` | content | GET | `/blogs.json` |
| `shopify_list_articles` | content | GET | `/blogs/{id}/articles.json` |
| `shopify_get_article` | content | GET | `/blogs/{id}/articles/{id}.json` |
| `shopify_create_article` | content | POST | `/blogs/{id}/articles.json` |
