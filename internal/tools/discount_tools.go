package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerDiscountTools(add toolAdder) {
	add(
		"shopify_list_price_rules",
		"List discount price rules with optional pagination and status filter.",
		baseProps(map[string]interface{}{
			"limit":    map[string]interface{}{"type": "integer", "description": "Number of price rules to return (default 50, max 250)"},
			"since_id": map[string]interface{}{"type": "string", "description": "Restrict results to after the specified ID"},
		}),
		nil,
		handleShopifyListPriceRules,
	)

	add(
		"shopify_get_price_rule",
		"Retrieve details of a specific price rule by ID.",
		baseProps(map[string]interface{}{
			"price_rule_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the price rule"},
		}),
		[]string{"price_rule_id"},
		handleShopifyGetPriceRule,
	)

	add(
		"shopify_create_price_rule",
		"Create a new discount price rule (percentage, fixed_amount, or free_shipping).",
		baseProps(map[string]interface{}{
			"title":              map[string]interface{}{"type": "string", "description": "Title of the price rule"},
			"target_type":        map[string]interface{}{"type": "string", "description": "line_item or shipping_line"},
			"target_selection":   map[string]interface{}{"type": "string", "description": "all or entitled"},
			"allocation_method":  map[string]interface{}{"type": "string", "description": "across or each"},
			"value_type":         map[string]interface{}{"type": "string", "description": "percentage or fixed_amount"},
			"value":              map[string]interface{}{"type": "string", "description": "Discount value (negative number, e.g. -15.0)"},
			"customer_selection": map[string]interface{}{"type": "string", "description": "all or prerequisite"},
			"starts_at":          map[string]interface{}{"type": "string", "description": "ISO 8601 start date and time"},
			"ends_at":            map[string]interface{}{"type": "string", "description": "ISO 8601 end date and time (optional)"},
		}),
		[]string{"title", "target_type", "target_selection", "allocation_method", "value_type", "value", "customer_selection", "starts_at"},
		handleShopifyCreatePriceRule,
	)

	add(
		"shopify_delete_price_rule",
		"Permanently delete a price rule by ID.",
		baseProps(map[string]interface{}{
			"price_rule_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the price rule to delete"},
		}),
		[]string{"price_rule_id"},
		handleShopifyDeletePriceRule,
	)

	add(
		"shopify_list_discount_codes",
		"List discount codes under a specific price rule.",
		baseProps(map[string]interface{}{
			"price_rule_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the price rule"},
			"limit":         map[string]interface{}{"type": "integer", "description": "Number of discount codes to return (default 50)"},
		}),
		[]string{"price_rule_id"},
		handleShopifyListDiscountCodes,
	)

	add(
		"shopify_create_discount_code",
		"Create a new discount code under an existing price rule.",
		baseProps(map[string]interface{}{
			"price_rule_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the price rule"},
			"code":          map[string]interface{}{"type": "string", "description": "Discount code string that customers enter at checkout"},
		}),
		[]string{"price_rule_id", "code"},
		handleShopifyCreateDiscountCode,
	)

	add(
		"shopify_lookup_discount_code",
		"Look up discount code details and its associated price rule by code string.",
		baseProps(map[string]interface{}{
			"code": map[string]interface{}{"type": "string", "description": "The discount code text to search for"},
		}),
		[]string{"code"},
		handleShopifyLookupDiscountCode,
	)
}

func handleShopifyListPriceRules(args json.RawMessage) map[string]interface{} {
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

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/price_rules.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetPriceRule(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ruleID := strArg(m, "price_rule_id")
	if ruleID == "" {
		return mcp.ToolResultError("price_rule_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/price_rules/%s.json", ruleID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreatePriceRule(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	rule := map[string]interface{}{
		"title":              strArg(m, "title"),
		"target_type":        strArg(m, "target_type"),
		"target_selection":   strArg(m, "target_selection"),
		"allocation_method":  strArg(m, "allocation_method"),
		"value_type":         strArg(m, "value_type"),
		"value":              strArg(m, "value"),
		"customer_selection": strArg(m, "customer_selection"),
		"starts_at":          strArg(m, "starts_at"),
	}
	if ends := strArg(m, "ends_at"); ends != "" {
		rule["ends_at"] = ends
	}

	payload := map[string]interface{}{"price_rule": rule}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/price_rules.json", url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyDeletePriceRule(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ruleID := strArg(m, "price_rule_id")
	if ruleID == "" {
		return mcp.ToolResultError("price_rule_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodDelete, fmt.Sprintf("/price_rules/%s.json", ruleID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListDiscountCodes(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ruleID := strArg(m, "price_rule_id")
	if ruleID == "" {
		return mcp.ToolResultError("price_rule_id is required")
	}
	q := url.Values{}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/price_rules/%s/discount_codes.json", ruleID), q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateDiscountCode(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ruleID := strArg(m, "price_rule_id")
	code := strArg(m, "code")
	if ruleID == "" || code == "" {
		return mcp.ToolResultError("price_rule_id and code are required")
	}

	payload := map[string]interface{}{
		"discount_code": map[string]interface{}{"code": code},
	}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, fmt.Sprintf("/price_rules/%s/discount_codes.json", ruleID), url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyLookupDiscountCode(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	code := strArg(m, "code")
	if code == "" {
		return mcp.ToolResultError("code is required")
	}
	q := url.Values{}
	q.Set("code", code)

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/discount_codes/lookup.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
