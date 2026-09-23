package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerFulfillmentTools(add toolAdder) {
	add(
		"shopify_list_fulfillments",
		"List all fulfillments for a specific order by order ID.",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the order"},
		}),
		[]string{"order_id"},
		handleShopifyListFulfillments,
	)

	add(
		"shopify_get_fulfillment",
		"Retrieve details of a specific fulfillment for an order.",
		baseProps(map[string]interface{}{
			"order_id":       map[string]interface{}{"type": "string", "description": "The numeric ID of the order"},
			"fulfillment_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the fulfillment"},
		}),
		[]string{"order_id", "fulfillment_id"},
		handleShopifyGetFulfillment,
	)

	add(
		"shopify_create_fulfillment",
		"Create a fulfillment for an order or fulfillment order with tracking number and company.",
		baseProps(map[string]interface{}{
			"fulfillment": map[string]interface{}{
				"type":        "object",
				"description": "Fulfillment payload object (e.g. line_items_by_fulfillment_order, tracking_info, notify_customer)",
			},
		}),
		[]string{"fulfillment"},
		handleShopifyCreateFulfillment,
	)

	add(
		"shopify_cancel_fulfillment",
		"Cancel an existing fulfillment by fulfillment ID.",
		baseProps(map[string]interface{}{
			"fulfillment_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the fulfillment to cancel"},
		}),
		[]string{"fulfillment_id"},
		handleShopifyCancelFulfillment,
	)
}

func handleShopifyListFulfillments(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	orderID := strArg(m, "order_id")
	if orderID == "" {
		return mcp.ToolResultError("order_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/orders/%s/fulfillments.json", orderID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetFulfillment(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	orderID := strArg(m, "order_id")
	fulfillmentID := strArg(m, "fulfillment_id")
	if orderID == "" || fulfillmentID == "" {
		return mcp.ToolResultError("order_id and fulfillment_id are required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/orders/%s/fulfillments/%s.json", orderID, fulfillmentID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateFulfillment(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	fulfillment, _ := m["fulfillment"].(map[string]interface{})
	if len(fulfillment) == 0 {
		return mcp.ToolResultError("fulfillment object is required")
	}

	payload := map[string]interface{}{"fulfillment": fulfillment}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/fulfillments.json", url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCancelFulfillment(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	fulfillmentID := strArg(m, "fulfillment_id")
	if fulfillmentID == "" {
		return mcp.ToolResultError("fulfillment_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, fmt.Sprintf("/fulfillments/%s/cancel.json", fulfillmentID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
