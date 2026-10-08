package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
	if len(descs) != 51 {
		t.Fatalf("expected full tool scope (51 tools), got %d", len(descs))
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

	for _, want := range wantedTools {
		if _, ok := handlers[want]; !ok {
			t.Fatalf("missing tool handler: %s", want)
		}
	}

	for _, desc := range descs {
		doc, ok := desc.InputSchema["x-datumbridge-docs"].(string)
		if !ok || doc == "" {
			t.Fatalf("missing x-datumbridge-docs in tool descriptor for: %s", desc.Name)
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

func TestHandleListOrdersFilters(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api/2024-04/orders.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("status") != "closed" {
			t.Errorf("expected status 'closed', got %q", q.Get("status"))
		}
		if q.Get("financial_status") != "paid" {
			t.Errorf("expected financial_status 'paid', got %q", q.Get("financial_status"))
		}
		if q.Get("fulfillment_status") != "shipped" {
			t.Errorf("expected fulfillment_status 'shipped', got %q", q.Get("fulfillment_status"))
		}
		if q.Get("created_at_min") != "2024-05-01T00:00:00Z" {
			t.Errorf("expected created_at_min '2024-05-01T00:00:00Z', got %q", q.Get("created_at_min"))
		}
		if q.Get("created_at_max") != "2024-05-31T23:59:59Z" {
			t.Errorf("expected created_at_max '2024-05-31T23:59:59Z', got %q", q.Get("created_at_max"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"orders":[{"id":101},{"id":102}]}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["status"] = "closed"
	m["financial_status"] = "paid"
	m["fulfillment_status"] = "shipped"
	m["created_at_min"] = "2024-05-01T00:00:00Z"
	m["created_at_max"] = "2024-05-31T23:59:59Z"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_list_orders"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleListOrdersDateRangeDefaults(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api/2024-04/orders.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("status") != "any" {
			t.Errorf("expected status 'any', got %q", q.Get("status"))
		}
		if q.Get("created_at_min") != "2024-05-01T00:00:00Z" {
			t.Errorf("expected created_at_min '2024-05-01T00:00:00Z', got %q", q.Get("created_at_min"))
		}
		if q.Get("created_at_max") != "2024-05-31T23:59:59Z" {
			t.Errorf("expected created_at_max '2024-05-31T23:59:59Z', got %q", q.Get("created_at_max"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"orders":[{"id":501}]}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["created_at_min"] = "2024-05-01"
	m["created_at_max"] = "2024-05-31"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_list_orders"](raw)
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

func TestHandleGraphQLSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api/2024-04/graphql.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %q", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"orders":{"edges":[]}}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["query"] = "{ orders(first: 5) { edges { node { id name } } } }"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_graphql"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleListOrdersProtectedCustomerDataFallback(t *testing.T) {
	clearToolFilterEnv(t)
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// First call without fields fails with 403 Forbidden due to Protected Customer Data
			if r.URL.Query().Get("fields") != "" {
				t.Errorf("expected first call to have no fields, got: %q", r.URL.Query().Get("fields"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"errors":"[API] This action requires merchant approval for read_customers or protected customer data"}`))
			return
		}
		// Second call (automatic fallback retry) should include safe fields
		fields := r.URL.Query().Get("fields")
		if fields == "" {
			t.Errorf("expected second call to specify fields")
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
	if callCount != 2 {
		t.Fatalf("expected 2 calls (initial + safe fields fallback), got %d", callCount)
	}
}

func TestHandleListInventoryLevels_AutoLocations(t *testing.T) {
	clearToolFilterEnv(t)
	var locationsCalled, inventoryCalled bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/admin/api/2024-04/locations.json" {
			locationsCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"locations":[{"id":777888,"active":true}]}`))
			return
		}
		if r.URL.Path == "/admin/api/2024-04/inventory_levels.json" {
			inventoryCalled = true
			if r.URL.Query().Get("location_ids") != "777888" {
				t.Errorf("expected location_ids '777888', got: %q", r.URL.Query().Get("location_ids"))
			}
			if r.URL.Query().Get("limit") != "200" {
				t.Errorf("expected limit '200', got: %q", r.URL.Query().Get("limit"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory_levels":[{"inventory_item_id":123,"location_id":777888,"available":10}]}`))
			return
		}
		t.Errorf("unexpected path: %q", r.URL.Path)
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["limit"] = 200
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_list_inventory_levels"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
	if !locationsCalled || !inventoryCalled {
		t.Fatalf("expected both locations and inventory to be called; locations=%v, inventory=%v", locationsCalled, inventoryCalled)
	}
}

func TestHandleListInventoryLevels_ExplicitLocationIDs(t *testing.T) {
	clearToolFilterEnv(t)
	var locationsCalled bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/admin/api/2024-04/locations.json" {
			locationsCalled = true
		}
		if r.URL.Path == "/admin/api/2024-04/inventory_levels.json" {
			if r.URL.Query().Get("location_ids") != "999" {
				t.Errorf("expected location_ids '999', got: %q", r.URL.Query().Get("location_ids"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory_levels":[]}`))
			return
		}
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["location_ids"] = "999"
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_list_inventory_levels"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
	if locationsCalled {
		t.Fatalf("did not expect /locations.json to be called when location_ids is provided explicitly")
	}
}

func TestHandleCreateCustomerMissingEmail(t *testing.T) {
	clearToolFilterEnv(t)
	_, handlers := tools.Register()
	res := handlers["shopify_create_customer"](json.RawMessage(`{"credentials_json":"{\"shop\":\"s.myshopify.com\",\"token\":\"t\"}"}`))
	if res["isError"] != true {
		t.Fatalf("expected error result for missing email, got: %+v", res)
	}
}

func TestHandleUpdateArticleSuccess(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/admin/api/2024-04/blogs/100/articles/200.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		art, ok := body["article"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected article key in body, got: %+v", body)
		}
		if art["published"] != true {
			t.Errorf("expected published true, got %v", art["published"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"article":{"id":200,"published":true}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["blog_id"] = "100"
	m["article_id"] = "200"
	m["is_published"] = true
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_update_article"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleUpdateArticleWithoutBlogID_GraphQL(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/admin/api/2024-04/graphql.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		vars := body["variables"].(map[string]interface{})
		if vars["id"] != "gid://shopify/Article/622314520859" {
			t.Errorf("expected GID for article, got %v", vars["id"])
		}
		art := vars["article"].(map[string]interface{})
		if art["isPublished"] != false {
			t.Errorf("expected isPublished false, got %v", art["isPublished"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"articleUpdate":{"article":{"id":"gid://shopify/Article/622314520859","isPublished":false}}}}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["article_id"] = "articles/622314520859"
	m["is_published"] = false
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_update_article"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}
}

func TestHandleUpdateArticleAllFieldsFakeData(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/admin/api/2024-04/graphql.json" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		vars := body["variables"].(map[string]interface{})
		if vars["id"] != "gid://shopify/Article/622314520859" {
			t.Errorf("expected GID 'gid://shopify/Article/622314520859', got %v", vars["id"])
		}
		art := vars["article"].(map[string]interface{})
		if art["title"] != "Hướng Dẫn Mua Sắm Mùa Hè 2024" {
			t.Errorf("expected title, got %v", art["title"])
		}
		if art["body"] != "<h1>Chào hè rực rỡ</h1><p>Khám phá bộ sưu tập sản phẩm mới nhất cùng nhiều ưu đãi hấp dẫn.</p>" {
			t.Errorf("expected body, got %v", art["body"])
		}
		auth, _ := art["author"].(map[string]interface{})
		if auth["name"] != "DatumBridge Team" {
			t.Errorf("expected author name 'DatumBridge Team', got %v", auth["name"])
		}
		tags, _ := art["tags"].([]interface{})
		if len(tags) != 4 || tags[0] != "summer" || tags[1] != "sale" {
			t.Errorf("expected tags [summer sale huong-dan fashion], got %v", tags)
		}
		if art["summary"] != "<p>Tổng hợp các mẹo mua sắm tiết kiệm và bộ sưu tập hè 2024.</p>" {
			t.Errorf("expected summary, got %v", art["summary"])
		}
		if art["isPublished"] != true {
			t.Errorf("expected isPublished true, got %v", art["isPublished"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"articleUpdate": {
					"article": {
						"id": "gid://shopify/Article/622314520859",
						"title": "Hướng Dẫn Mua Sắm Mùa Hè 2024",
						"isPublished": true,
						"publishedAt": "2024-06-01T10:00:00Z"
					},
					"userErrors": []
				}
			}
		}`))
	}))
	defer srv.Close()

	args := testCreds(t, srv)
	var m map[string]interface{}
	_ = json.Unmarshal(args, &m)
	m["article_id"] = "622314520859"
	m["title"] = "Hướng Dẫn Mua Sắm Mùa Hè 2024"
	m["body_html"] = "<h1>Chào hè rực rỡ</h1><p>Khám phá bộ sưu tập sản phẩm mới nhất cùng nhiều ưu đãi hấp dẫn.</p>"
	m["author"] = "DatumBridge Team"
	m["tags"] = "summer, sale, huong-dan, fashion"
	m["summary_html"] = "<p>Tổng hợp các mẹo mua sắm tiết kiệm và bộ sưu tập hè 2024.</p>"
	m["is_published"] = true
	raw, _ := json.Marshal(m)

	_, handlers := tools.Register()
	res := handlers["shopify_update_article"](raw)
	if res["isError"] == true {
		t.Fatalf("unexpected error result: %+v", res)
	}

	content, ok := res["content"].([]map[string]string)
	if !ok || len(content) == 0 {
		t.Fatalf("expected content in result, got: %+v", res)
	}
	text := content[0]["text"]
	if !strings.Contains(text, "Hướng Dẫn Mua Sắm Mùa Hè 2024") {
		t.Fatalf("expected text to contain title, got: %s", text)
	}
}




