package shopify_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/shopify"
)

func TestClient_RequestSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Shopify-Access-Token") != "test-token" {
			t.Fatalf("expected X-Shopify-Access-Token header, got %q", r.Header.Get("X-Shopify-Access-Token"))
		}
		if r.URL.Path != "/admin/api/2024-04/products.json" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"products":[{"id":123,"title":"T-Shirt"}]}`))
	}))
	defer ts.Close()

	t.Setenv("SHOPIFY_API_BASE", ts.URL)
	creds := &shopify.Credentials{
		Shop:  "my-store.myshopify.com",
		Token: "test-token",
	}
	client := shopify.NewClient(creds)

	data, status, err := client.Request(context.Background(), http.MethodGet, "/products.json", url.Values{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if string(data) != `{"products":[{"id":123,"title":"T-Shirt"}]}` {
		t.Fatalf("unexpected data: %s", string(data))
	}
}

func TestClient_RequestErrorStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"errors":"[API] Invalid API key or access token (unrecognized login or wrong password)"}`))
	}))
	defer ts.Close()

	t.Setenv("SHOPIFY_API_BASE", ts.URL)
	creds := &shopify.Credentials{
		Shop:  "my-store.myshopify.com",
		Token: "invalid-token",
	}
	client := shopify.NewClient(creds)

	_, status, err := client.Request(context.Background(), http.MethodGet, "/products.json", url.Values{}, nil)
	if err == nil {
		t.Fatal("expected error on 401 response")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", status)
	}
}
