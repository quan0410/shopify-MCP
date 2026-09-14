package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/tools"
)

func clearToolFilterEnv(t *testing.T) {
	t.Helper()
	t.Setenv("SHOPIFY_TOOLS", "")
}

func TestRegisterFullScopeByDefault(t *testing.T) {
	clearToolFilterEnv(t)
	descs, handlers := tools.Register()
	if len(descs) != 29 {
		t.Fatalf("expected full tool scope (29 tools), got %d", len(descs))
	}

	wantedTools := []string{
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
	}

	for _, want := range wantedTools {
		if _, ok := handlers[want]; !ok {
			t.Fatalf("missing tool handler: %s", want)
		}
	}
}

func TestRegisterCustomAllowlist(t *testing.T) {
	clearToolFilterEnv(t)
	t.Setenv("SHOPIFY_TOOLS", "shopify_get_product,shopify_list_orders")
	descs, handlers := tools.Register()
	if len(descs) != 2 {
		t.Fatalf("expected 2 allowlisted tools, got %d", len(descs))
	}
	if _, ok := handlers["shopify_get_product"]; !ok {
		t.Fatal("shopify_get_product should be enabled")
	}
	if _, ok := handlers["shopify_list_orders"]; !ok {
		t.Fatal("shopify_list_orders should be enabled")
	}
	if _, ok := handlers["shopify_create_product"]; ok {
		t.Fatal("shopify_create_product should be filtered out")
	}
}

func testCreds(t *testing.T, srv *httptest.Server) json.RawMessage {
	t.Helper()
	t.Setenv("SHOPIFY_API_BASE", srv.URL)
	creds := `{"shop":"test-shop.myshopify.com","token":"shpat_mock"}`
	raw, _ := json.Marshal(map[string]string{"credentials_json": creds})
	return raw
}

func TestHandleGetProductSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Shopify-Access-Token") != "shpat_mock" {
			t.Errorf("token header missing or invalid: %q", r.Header.Get("X-Shopify-Access-Token"))
		}
		if r.URL.Path != "/admin/api/2024-04/products/12345.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"product":{"id":12345,"title":"Premium Silk Shirt"}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["product_id"] = "12345"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_get_product"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleGetProductMissingID(t *testing.T) {
	clearToolFilterEnv(t)
	_, handlers := tools.Register()
	res := handlers["shopify_get_product"](json.RawMessage(`{"credentials_json":"{\"shop\":\"s.myshopify.com\",\"token\":\"t\"}"}`))
	if res["isError"] != true {
		t.Fatalf("expected error result for missing product_id, got: %+v", res)
	}
}

func TestHandleCreateProductSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		prod := body["product"].(map[string]interface{})
		if prod["title"] != "New Sneaker" {
			t.Errorf("expected title 'New Sneaker', got %v", prod["title"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"product":{"id":999,"title":"New Sneaker"}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["title"] = "New Sneaker"
	m["vendor"] = "Nike"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_create_product"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleListOrdersSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api/2024-04/orders.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"orders":[{"id":1001,"financial_status":"paid"}]}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	_, handlers := tools.Register()
	res := handlers["shopify_list_orders"](args)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleGetShopSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api/2024-04/shop.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"shop":{"name":"Demo Shop","domain":"test-shop.myshopify.com"}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	_, handlers := tools.Register()
	res := handlers["shopify_get_shop"](args)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}
