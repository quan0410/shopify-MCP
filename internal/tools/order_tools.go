package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

// defaultSafeOrderFields specifies common order fields that exclude Personally Identifiable Information (PII)
// such as customer, email, phone, and addresses, allowing apps without Protected Customer Data (PCD) approval to query orders.
const defaultSafeOrderFields = "id,name,order_number,created_at,financial_status,fulfillment_status,total_price,subtotal_price,currency,line_items,cancel_reason,cancelled_at,closed_at"

func isProtectedDataError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "protected") ||
		strings.Contains(msg, "customer data") ||
		strings.Contains(msg, "read_customers") ||
		strings.Contains(msg, "not permitted") ||
		strings.Contains(msg, "403")
}

func registerOrderTools(add toolAdder) {
	add(
		"shopify_list_orders",
		"List orders from the store with filters for status, financial_status, fulfillment_status, and pagination. To avoid Shopify 403 Forbidden / Protected Customer Data errors when the app lacks PCD approval, specify 'fields' (e.g. 'id,name,created_at,financial_status,fulfillment_status,total_price,currency,line_items') or set safe_mode to true.",
		baseProps(map[string]interface{}{
			"limit":              map[string]interface{}{"type": "integer", "description": "Number of orders to retrieve (default 50, max 250)"},
			"since_id":           map[string]interface{}{"type": "string", "description": "Restrict results to after the specified order ID"},
			"status":             map[string]interface{}{"type": "string", "description": "Filter by status: open, closed, cancelled, any (default open)"},
			"financial_status":   map[string]interface{}{"type": "string", "description": "Filter by financial status: authorized, pending, paid, refunded, voided, any"},
			"fulfillment_status": map[string]interface{}{"type": "string", "description": "Filter by fulfillment status: shipped, partial, unshipped, any"},
			"fields":             map[string]interface{}{"type": "string", "description": "Comma-separated list of fields to retrieve (e.g. 'id,name,order_number,created_at,financial_status,fulfillment_status,total_price,currency,line_items'). Excludes customer/address PII to prevent Shopify 403 Protected Customer Data errors."},
			"safe_mode":          map[string]interface{}{"type": "boolean", "description": "If true, automatically queries only non-PII fields to avoid Protected Customer Data restrictions (default false)."},
		}),
		nil,
		handleShopifyListOrders,
	)

	add(
		"shopify_get_order",
		"Retrieve single order details by order ID. Use 'fields' or 'safe_mode' to exclude customer/address PII and avoid Shopify 403 Protected Customer Data errors.",
		baseProps(map[string]interface{}{
			"order_id":  map[string]interface{}{"type": "string", "description": "Order ID"},
			"fields":    map[string]interface{}{"type": "string", "description": "Comma-separated list of fields to retrieve (e.g. 'id,name,order_number,created_at,financial_status,fulfillment_status,total_price,currency,line_items')."},
			"safe_mode": map[string]interface{}{"type": "boolean", "description": "If true, queries only non-PII fields to avoid Protected Customer Data restrictions (default false)."},
		}),
		[]string{"order_id"},
		handleShopifyGetOrder,
	)

	add(
		"shopify_create_order",
		"Create a new order in Shopify with line items, customer information, shipping/billing address, and financial status.",
		baseProps(map[string]interface{}{
			"line_items":       map[string]interface{}{"type": "array", "description": "List of line items [{variant_id, quantity, price}]"},
			"customer":         map[string]interface{}{"type": "object", "description": "Customer object {id, email, first_name, last_name}"},
			"billing_address":  map[string]interface{}{"type": "object", "description": "Billing address object"},
			"shipping_address": map[string]interface{}{"type": "object", "description": "Shipping address object"},
			"financial_status": map[string]interface{}{"type": "string", "description": "paid, pending, authorized, etc."},
			"note":             map[string]interface{}{"type": "string", "description": "Order notes"},
			"tags":             map[string]interface{}{"type": "string", "description": "Comma separated tags"},
		}),
		[]string{"line_items"},
		handleShopifyCreateOrder,
	)

	add(
		"shopify_cancel_order",
		"Cancel an existing order with optional reason and refund flag.",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{"type": "string", "description": "Order ID to cancel"},
			"reason":   map[string]interface{}{"type": "string", "description": "Reason for cancellation: customer, inventory, fraud, declined, other"},
			"email":    map[string]interface{}{"type": "boolean", "description": "Whether to send an email to the customer (default true)"},
		}),
		[]string{"order_id"},
		handleShopifyCancelOrder,
	)

	add(
		"shopify_close_order",
		"Close an open order that has been completed.",
		baseProps(map[string]interface{}{
			"order_id": map[string]interface{}{"type": "string", "description": "Order ID to close"},
		}),
		[]string{"order_id"},
		handleShopifyCloseOrder,
	)

	add(
		"shopify_list_draft_orders",
		"List draft orders with optional status filter, fields, and pagination.",
		baseProps(map[string]interface{}{
			"limit":     map[string]interface{}{"type": "integer", "description": "Number of draft orders to return"},
			"since_id":  map[string]interface{}{"type": "string", "description": "Filter draft orders after specified ID"},
			"status":    map[string]interface{}{"type": "string", "description": "open, invoice_sent, completed"},
			"fields":    map[string]interface{}{"type": "string", "description": "Comma-separated list of fields to retrieve (e.g. 'id,name,status,total_price,line_items')."},
			"safe_mode": map[string]interface{}{"type": "boolean", "description": "If true, queries only non-PII fields to avoid Protected Customer Data restrictions (default false)."},
		}),
		nil,
		handleShopifyListDraftOrders,
	)

	add(
		"shopify_create_draft_order",
		"Create a draft order with line items, customer details, and notes.",
		baseProps(map[string]interface{}{
			"line_items": map[string]interface{}{"type": "array", "description": "List of line items [{variant_id, quantity, price}]"},
			"customer":   map[string]interface{}{"type": "object", "description": "Customer object {id, email}"},
			"note":       map[string]interface{}{"type": "string", "description": "Note for draft order"},
			"tags":       map[string]interface{}{"type": "string", "description": "Tags"},
		}),
		[]string{"line_items"},
		handleShopifyCreateDraftOrder,
	)
}

func handleShopifyListOrders(args json.RawMessage) map[string]interface{} {
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
	if st := strArg(m, "status"); st != "" {
		q.Set("status", st)
	}
	if fs := strArg(m, "financial_status"); fs != "" {
		q.Set("financial_status", fs)
	}
	if ful := strArg(m, "fulfillment_status"); ful != "" {
		q.Set("fulfillment_status", ful)
	}
	fields := strArg(m, "fields")
	if fields == "" && boolArg(m, "safe_mode", false) {
		fields = defaultSafeOrderFields
	}
	if fields != "" {
		q.Set("fields", fields)
	}

	c, cancel := ctx()
	defer cancel()
	data, code, err := client.Request(c, http.MethodGet, "/orders.json", q, nil)
	if err != nil {
		// Fallback: If failed due to Protected Customer Data and fields wasn't explicitly specified, retry with safe non-PII fields
		if (code == http.StatusForbidden || isProtectedDataError(err)) && fields == "" {
			q.Set("fields", defaultSafeOrderFields)
			retryData, _, retryErr := client.Request(c, http.MethodGet, "/orders.json", q, nil)
			if retryErr == nil {
				return rawResult(retryData)
			}
		}
		if code == http.StatusForbidden || isProtectedDataError(err) {
			return mcp.ToolResultError(fmt.Sprintf("%s. Hint: App may lack approval for Protected Customer Data (PCD). Pass 'fields' with non-PII fields (e.g. '%s') or use shopify_graphql.", err.Error(), defaultSafeOrderFields))
		}
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetOrder(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	orderID := strArg(m, "order_id")
	if orderID == "" {
		return mcp.ToolResultError("order_id is required")
	}

	q := url.Values{}
	fields := strArg(m, "fields")
	if fields == "" && boolArg(m, "safe_mode", false) {
		fields = defaultSafeOrderFields
	}
	if fields != "" {
		q.Set("fields", fields)
	}

	c, cancel := ctx()
	defer cancel()
	data, code, err := client.Request(c, http.MethodGet, fmt.Sprintf("/orders/%s.json", orderID), q, nil)
	if err != nil {
		// Fallback: If failed due to Protected Customer Data and fields wasn't explicitly specified, retry with safe non-PII fields
		if (code == http.StatusForbidden || isProtectedDataError(err)) && fields == "" {
			q.Set("fields", defaultSafeOrderFields)
			retryData, _, retryErr := client.Request(c, http.MethodGet, fmt.Sprintf("/orders/%s.json", orderID), q, nil)
			if retryErr == nil {
				return rawResult(retryData)
			}
		}
		if code == http.StatusForbidden || isProtectedDataError(err) {
			return mcp.ToolResultError(fmt.Sprintf("%s. Hint: App may lack approval for Protected Customer Data (PCD). Pass 'fields' with non-PII fields (e.g. '%s') or use shopify_graphql.", err.Error(), defaultSafeOrderFields))
		}
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateOrder(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	lineItems := sliceMapArg(m, "line_items")
	if len(lineItems) == 0 {
		return mcp.ToolResultError("line_items is required and must not be empty")
	}

	order := map[string]interface{}{
		"line_items": lineItems,
	}
	if cust := mapArg(m, "customer"); len(cust) > 0 {
		order["customer"] = cust
	}
	if bill := mapArg(m, "billing_address"); len(bill) > 0 {
		order["billing_address"] = bill
	}
	if ship := mapArg(m, "shipping_address"); len(ship) > 0 {
		order["shipping_address"] = ship
	}
	if fs := strArg(m, "financial_status"); fs != "" {
		order["financial_status"] = fs
	}
	if note := strArg(m, "note"); note != "" {
		order["note"] = note
	}
	if tags := strArg(m, "tags"); tags != "" {
		order["tags"] = tags
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/orders.json", url.Values{}, map[string]interface{}{"order": order})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCancelOrder(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	orderID := strArg(m, "order_id")
	if orderID == "" {
		return mcp.ToolResultError("order_id is required")
	}

	body := map[string]interface{}{}
	if r := strArg(m, "reason"); r != "" {
		body["reason"] = r
	}
	if b := boolArg(m, "email", true); !b {
		body["email"] = false
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, fmt.Sprintf("/orders/%s/cancel.json", orderID), url.Values{}, body)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCloseOrder(args json.RawMessage) map[string]interface{} {
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
	data, _, err := client.Request(c, http.MethodPost, fmt.Sprintf("/orders/%s/close.json", orderID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListDraftOrders(args json.RawMessage) map[string]interface{} {
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
	if st := strArg(m, "status"); st != "" {
		q.Set("status", st)
	}
	fields := strArg(m, "fields")
	if fields == "" && boolArg(m, "safe_mode", false) {
		fields = "id,name,status,total_price,subtotal_price,currency,line_items,created_at,updated_at"
	}
	if fields != "" {
		q.Set("fields", fields)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/draft_orders.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateDraftOrder(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	lineItems := sliceMapArg(m, "line_items")
	if len(lineItems) == 0 {
		return mcp.ToolResultError("line_items is required and must not be empty")
	}

	draft := map[string]interface{}{
		"line_items": lineItems,
	}
	if cust := mapArg(m, "customer"); len(cust) > 0 {
		draft["customer"] = cust
	}
	if note := strArg(m, "note"); note != "" {
		draft["note"] = note
	}
	if tags := strArg(m, "tags"); tags != "" {
		draft["tags"] = tags
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/draft_orders.json", url.Values{}, map[string]interface{}{"draft_order": draft})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
