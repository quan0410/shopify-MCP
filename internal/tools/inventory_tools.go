package tools

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerInventoryTools(add toolAdder) {
	add(
		"shopify_list_locations",
		"List all inventory locations configured on the Shopify store.",
		baseProps(nil),
		nil,
		handleShopifyListLocations,
	)

	add(
		"shopify_list_inventory_levels",
		"Retrieve inventory levels for specific inventory_item_ids or location_ids. Shopify requires at least one of location_ids or inventory_item_ids; if neither is provided, active store locations are automatically queried.",
		baseProps(map[string]interface{}{
			"location_ids":       map[string]interface{}{"type": "string", "description": "Comma-separated list or array of location IDs. Required by Shopify unless inventory_item_ids is provided (auto-detected from active store locations if omitted)."},
			"inventory_item_ids": map[string]interface{}{"type": "string", "description": "Comma-separated list or array of inventory item IDs. Required by Shopify unless location_ids is provided."},
			"limit":              map[string]interface{}{"type": "integer", "description": "Number of inventory levels to return (default 50, max 250)"},
		}),
		nil,
		handleShopifyListInventoryLevels,
	)

	add(
		"shopify_adjust_inventory_level",
		"Adjust available inventory quantity for an inventory item at a location by a delta amount (e.g. +5 or -2).",
		baseProps(map[string]interface{}{
			"inventory_item_id":    map[string]interface{}{"type": "integer", "description": "Numeric inventory item ID"},
			"location_id":          map[string]interface{}{"type": "integer", "description": "Numeric location ID"},
			"available_adjustment": map[string]interface{}{"type": "integer", "description": "Quantity adjustment delta (positive or negative integer)"},
		}),
		[]string{"inventory_item_id", "location_id", "available_adjustment"},
		handleShopifyAdjustInventoryLevel,
	)

	add(
		"shopify_set_inventory_level",
		"Set exact available inventory quantity for an inventory item at a location.",
		baseProps(map[string]interface{}{
			"inventory_item_id": map[string]interface{}{"type": "integer", "description": "Numeric inventory item ID"},
			"location_id":       map[string]interface{}{"type": "integer", "description": "Numeric location ID"},
			"available":         map[string]interface{}{"type": "integer", "description": "Exact available quantity to set"},
		}),
		[]string{"inventory_item_id", "location_id", "available"},
		handleShopifySetInventoryLevel,
	)
}

func handleShopifyListLocations(args json.RawMessage) map[string]interface{} {
	client, _, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/locations.json", url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListInventoryLevels(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	itemIDs := strOrSliceArg(m, "inventory_item_ids")
	if itemIDs == "" {
		itemIDs = strOrSliceArg(m, "inventory_item_id")
	}
	locIDs := strOrSliceArg(m, "location_ids")
	if locIDs == "" {
		locIDs = strOrSliceArg(m, "location_id")
	}

	c, cancel := ctx()
	defer cancel()

	// Shopify requires at least one of inventory_item_ids or location_ids.
	// If neither is provided, automatically query active store locations to prevent HTTP 422.
	if itemIDs == "" && locIDs == "" {
		locData, _, locErr := client.Request(c, http.MethodGet, "/locations.json", url.Values{}, nil)
		if locErr == nil {
			var locResp struct {
				Locations []struct {
					ID     int64 `json:"id"`
					Active bool  `json:"active"`
				} `json:"locations"`
			}
			if err := json.Unmarshal(locData, &locResp); err == nil && len(locResp.Locations) > 0 {
				var ids []string
				for _, loc := range locResp.Locations {
					if loc.Active {
						ids = append(ids, strconv.FormatInt(loc.ID, 10))
					}
				}
				if len(ids) == 0 {
					for _, loc := range locResp.Locations {
						ids = append(ids, strconv.FormatInt(loc.ID, 10))
					}
				}
				if len(ids) > 0 {
					locIDs = strings.Join(ids, ",")
				}
			}
		}
	}

	if itemIDs == "" && locIDs == "" {
		return mcp.ToolResultError("at least one of 'location_ids' or 'inventory_item_ids' is required by Shopify to list inventory levels (use shopify_list_locations to discover available locations)")
	}

	q := url.Values{}
	if itemIDs != "" {
		q.Set("inventory_item_ids", itemIDs)
	}
	if locIDs != "" {
		q.Set("location_ids", locIDs)
	}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}

	data, _, err := client.Request(c, http.MethodGet, "/inventory_levels.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyAdjustInventoryLevel(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	itemID := intArg(m, "inventory_item_id", 0)
	locID := intArg(m, "location_id", 0)
	adj := intArg(m, "available_adjustment", 0)
	if itemID == 0 || locID == 0 {
		return mcp.ToolResultError("inventory_item_id and location_id are required")
	}

	body := map[string]interface{}{
		"inventory_item_id":    itemID,
		"location_id":          locID,
		"available_adjustment": adj,
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/inventory_levels/adjust.json", url.Values{}, body)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifySetInventoryLevel(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	itemID := intArg(m, "inventory_item_id", 0)
	locID := intArg(m, "location_id", 0)
	avail := intArg(m, "available", 0)
	if itemID == 0 || locID == 0 {
		return mcp.ToolResultError("inventory_item_id and location_id are required")
	}

	body := map[string]interface{}{
		"inventory_item_id": itemID,
		"location_id":       locID,
		"available":         avail,
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/inventory_levels/set.json", url.Values{}, body)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
