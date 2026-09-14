package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerProductTools(add toolAdder) {
	add(
		"shopify_list_products",
		"List products with optional pagination and filtering by title, status, vendor, or product_type.",
		baseProps(map[string]interface{}{
			"limit":        map[string]interface{}{"type": "integer", "description": "Number of products to return (default 50, max 250)"},
			"since_id":     map[string]interface{}{"type": "string", "description": "Restrict results to after the specified product ID"},
			"title":        map[string]interface{}{"type": "string", "description": "Filter products by title"},
			"vendor":       map[string]interface{}{"type": "string", "description": "Filter products by vendor"},
			"product_type": map[string]interface{}{"type": "string", "description": "Filter products by product type"},
			"status":       map[string]interface{}{"type": "string", "description": "Filter by status: active, archived, draft"},
		}),
		nil,
		handleShopifyListProducts,
	)

	add(
		"shopify_get_product",
		"Retrieve full details of a specific product by product ID.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the product"},
		}),
		[]string{"product_id"},
		handleShopifyGetProduct,
	)

	add(
		"shopify_create_product",
		"Create a new product in Shopify with title, description, vendor, product_type, tags, and optional variants.",
		baseProps(map[string]interface{}{
			"title":        map[string]interface{}{"type": "string", "description": "Product title"},
			"body_html":    map[string]interface{}{"type": "string", "description": "Product description in HTML or plain text"},
			"vendor":       map[string]interface{}{"type": "string", "description": "Product vendor name"},
			"product_type": map[string]interface{}{"type": "string", "description": "Category/type of product"},
			"status":       map[string]interface{}{"type": "string", "description": "active, draft, or archived (default active)"},
			"tags":         map[string]interface{}{"type": "string", "description": "Comma-separated list of tags"},
			"variants":     map[string]interface{}{"type": "array", "description": "List of variant objects (price, sku, title, etc.)"},
		}),
		[]string{"title"},
		handleShopifyCreateProduct,
	)

	add(
		"shopify_update_product",
		"Update an existing product's fields (title, body_html, vendor, status, tags).",
		baseProps(map[string]interface{}{
			"product_id":   map[string]interface{}{"type": "string", "description": "Product ID to update"},
			"title":        map[string]interface{}{"type": "string", "description": "New title"},
			"body_html":    map[string]interface{}{"type": "string", "description": "New description"},
			"vendor":       map[string]interface{}{"type": "string", "description": "New vendor"},
			"product_type": map[string]interface{}{"type": "string", "description": "New product type"},
			"status":       map[string]interface{}{"type": "string", "description": "New status (active, draft, archived)"},
			"tags":         map[string]interface{}{"type": "string", "description": "New tags (comma separated)"},
		}),
		[]string{"product_id"},
		handleShopifyUpdateProduct,
	)

	add(
		"shopify_delete_product",
		"Permanently delete a product by product ID.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{"type": "string", "description": "Product ID to delete"},
		}),
		[]string{"product_id"},
		handleShopifyDeleteProduct,
	)

	add(
		"shopify_list_variants",
		"List all variants for a given product ID.",
		baseProps(map[string]interface{}{
			"product_id": map[string]interface{}{"type": "string", "description": "Product ID to list variants for"},
		}),
		[]string{"product_id"},
		handleShopifyListVariants,
	)

	add(
		"shopify_get_variant",
		"Retrieve a single product variant by variant ID.",
		baseProps(map[string]interface{}{
			"variant_id": map[string]interface{}{"type": "string", "description": "Variant ID"},
		}),
		[]string{"variant_id"},
		handleShopifyGetVariant,
	)

	add(
		"shopify_update_variant",
		"Update a product variant's price, sku, barcode, or inventory policy.",
		baseProps(map[string]interface{}{
			"variant_id": map[string]interface{}{"type": "string", "description": "Variant ID to update"},
			"price":      map[string]interface{}{"type": "string", "description": "New price"},
			"sku":        map[string]interface{}{"type": "string", "description": "New SKU"},
			"barcode":    map[string]interface{}{"type": "string", "description": "New barcode"},
			"option1":    map[string]interface{}{"type": "string", "description": "Option 1 value (e.g. Size)"},
			"option2":    map[string]interface{}{"type": "string", "description": "Option 2 value (e.g. Color)"},
		}),
		[]string{"variant_id"},
		handleShopifyUpdateVariant,
	)
}

func handleShopifyListProducts(args json.RawMessage) map[string]interface{} {
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
	if t := strArg(m, "title"); t != "" {
		q.Set("title", t)
	}
	if v := strArg(m, "vendor"); v != "" {
		q.Set("vendor", v)
	}
	if pt := strArg(m, "product_type"); pt != "" {
		q.Set("product_type", pt)
	}
	if st := strArg(m, "status"); st != "" {
		q.Set("status", st)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/products.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetProduct(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	productID := strArg(m, "product_id")
	if productID == "" {
		return mcp.ToolResultError("product_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/products/%s.json", productID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateProduct(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	title := strArg(m, "title")
	if title == "" {
		return mcp.ToolResultError("title is required")
	}

	product := map[string]interface{}{
		"title": title,
	}
	if v := strArg(m, "body_html"); v != "" {
		product["body_html"] = v
	}
	if v := strArg(m, "vendor"); v != "" {
		product["vendor"] = v
	}
	if v := strArg(m, "product_type"); v != "" {
		product["product_type"] = v
	}
	if v := strArg(m, "status"); v != "" {
		product["status"] = v
	}
	if v := strArg(m, "tags"); v != "" {
		product["tags"] = v
	}
	if variants := sliceMapArg(m, "variants"); len(variants) > 0 {
		product["variants"] = variants
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/products.json", url.Values{}, map[string]interface{}{"product": product})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyUpdateProduct(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	productID := strArg(m, "product_id")
	if productID == "" {
		return mcp.ToolResultError("product_id is required")
	}

	product := map[string]interface{}{}
	if v := strArg(m, "title"); v != "" {
		product["title"] = v
	}
	if v := strArg(m, "body_html"); v != "" {
		product["body_html"] = v
	}
	if v := strArg(m, "vendor"); v != "" {
		product["vendor"] = v
	}
	if v := strArg(m, "product_type"); v != "" {
		product["product_type"] = v
	}
	if v := strArg(m, "status"); v != "" {
		product["status"] = v
	}
	if v := strArg(m, "tags"); v != "" {
		product["tags"] = v
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPut, fmt.Sprintf("/products/%s.json", productID), url.Values{}, map[string]interface{}{"product": product})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyDeleteProduct(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	productID := strArg(m, "product_id")
	if productID == "" {
		return mcp.ToolResultError("product_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodDelete, fmt.Sprintf("/products/%s.json", productID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if len(data) == 0 {
		return jsonResult(map[string]string{"status": "deleted", "product_id": productID})
	}
	return rawResult(data)
}

func handleShopifyListVariants(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	productID := strArg(m, "product_id")
	if productID == "" {
		return mcp.ToolResultError("product_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/products/%s/variants.json", productID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetVariant(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	variantID := strArg(m, "variant_id")
	if variantID == "" {
		return mcp.ToolResultError("variant_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/variants/%s.json", variantID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyUpdateVariant(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	variantID := strArg(m, "variant_id")
	if variantID == "" {
		return mcp.ToolResultError("variant_id is required")
	}

	variant := map[string]interface{}{}
	if v := strArg(m, "price"); v != "" {
		variant["price"] = v
	}
	if v := strArg(m, "sku"); v != "" {
		variant["sku"] = v
	}
	if v := strArg(m, "barcode"); v != "" {
		variant["barcode"] = v
	}
	if v := strArg(m, "option1"); v != "" {
		variant["option1"] = v
	}
	if v := strArg(m, "option2"); v != "" {
		variant["option2"] = v
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPut, fmt.Sprintf("/variants/%s.json", variantID), url.Values{}, map[string]interface{}{"variant": variant})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}
