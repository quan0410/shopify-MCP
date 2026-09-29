package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerCustomerTools(add toolAdder) {
	add(
		"shopify_list_customers",
		"List customers registered on the Shopify store with optional pagination.",
		baseProps(map[string]interface{}{
			"limit":    map[string]interface{}{"type": "integer", "description": "Number of customers to return (default 50, max 250)"},
			"since_id": map[string]interface{}{"type": "string", "description": "Restrict results to after the specified customer ID"},
		}),
		nil,
		handleShopifyListCustomers,
	)

	add(
		"shopify_get_customer",
		"Retrieve detailed information for a single customer by customer ID.",
		baseProps(map[string]interface{}{
			"customer_id": map[string]interface{}{"type": "string", "description": "Customer ID"},
		}),
		[]string{"customer_id"},
		handleShopifyGetCustomer,
	)

	add(
		"shopify_search_customers",
		"Search for customers matching a query string (e.g. name, email, phone number, address).",
		baseProps(map[string]interface{}{
			"query": map[string]interface{}{"type": "string", "description": "Search query text"},
			"limit": map[string]interface{}{"type": "integer", "description": "Maximum number of results to return"},
		}),
		[]string{"query"},
		handleShopifySearchCustomers,
	)

	add(
		"shopify_create_customer",
		"Create a new customer profile with name, email, phone, and optional tags.",
		baseProps(map[string]interface{}{
			"email":      map[string]interface{}{"type": "string", "description": "Customer's email address"},
			"first_name": map[string]interface{}{"type": "string", "description": "Customer's first name"},
			"last_name":  map[string]interface{}{"type": "string", "description": "Customer's last name"},
			"phone":      map[string]interface{}{"type": "string", "description": "Customer's phone number"},
			"tags":       map[string]interface{}{"type": "string", "description": "Comma separated tags"},
			"note":       map[string]interface{}{"type": "string", "description": "Note about the customer"},
		}),
		[]string{"email"},
		handleShopifyCreateCustomer,
	)

	add(
		"shopify_update_customer",
		"Update customer details (name, email, phone, tags, notes) by customer ID.",
		baseProps(map[string]interface{}{
			"customer_id": map[string]interface{}{"type": "string", "description": "Customer ID to update"},
			"first_name":  map[string]interface{}{"type": "string", "description": "Updated first name"},
			"last_name":   map[string]interface{}{"type": "string", "description": "Updated last name"},
			"email":       map[string]interface{}{"type": "string", "description": "Updated email"},
			"phone":       map[string]interface{}{"type": "string", "description": "Updated phone number"},
			"tags":        map[string]interface{}{"type": "string", "description": "Updated tags"},
			"note":        map[string]interface{}{"type": "string", "description": "Updated notes"},
		}),
		[]string{"customer_id"},
		handleShopifyUpdateCustomer,
	)
}

func handleShopifyListCustomers(args json.RawMessage) map[string]interface{} {
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
	data, _, err := client.Request(c, http.MethodGet, "/customers.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetCustomer(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	customerID := strArg(m, "customer_id")
	if customerID == "" {
		return mcp.ToolResultError("customer_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/customers/%s.json", customerID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifySearchCustomers(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	query := strArg(m, "query")
	if query == "" {
		return mcp.ToolResultError("query is required")
	}
	q := url.Values{}
	q.Set("query", query)
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/customers/search.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateCustomer(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	em := strArg(m, "email")
	ph := strArg(m, "phone")
	if em == "" && ph == "" {
		return mcp.ToolResultError("email is required to create a customer (or phone)")
	}

	customer := map[string]interface{}{}
	if fn := strArg(m, "first_name"); fn != "" {
		customer["first_name"] = fn
	}
	if ln := strArg(m, "last_name"); ln != "" {
		customer["last_name"] = ln
	}
	if em != "" {
		customer["email"] = em
	}
	if ph != "" {
		customer["phone"] = ph
	}
	if tags := strArg(m, "tags"); tags != "" {
		customer["tags"] = tags
	}
	if note := strArg(m, "note"); note != "" {
		customer["note"] = note
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/customers.json", url.Values{}, map[string]interface{}{"customer": customer})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyUpdateCustomer(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	customerID := strArg(m, "customer_id")
	if customerID == "" {
		return mcp.ToolResultError("customer_id is required")
	}

	customer := map[string]interface{}{}
	if fn := strArg(m, "first_name"); fn != "" {
		customer["first_name"] = fn
	}
	if ln := strArg(m, "last_name"); ln != "" {
		customer["last_name"] = ln
	}
	if em := strArg(m, "email"); em != "" {
		customer["email"] = em
	}
	if ph := strArg(m, "phone"); ph != "" {
		customer["phone"] = ph
	}
	if tags := strArg(m, "tags"); tags != "" {
		customer["tags"] = tags
	}
	if note := strArg(m, "note"); note != "" {
		customer["note"] = note
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPut, fmt.Sprintf("/customers/%s.json", customerID), url.Values{}, map[string]interface{}{"customer": customer})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
