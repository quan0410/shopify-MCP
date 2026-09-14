package shopify_test

import (
	"os"
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/shopify"
)

func TestParseCredentials_InlineJSON(t *testing.T) {
	raw := `{"shop":"my-store.myshopify.com","token":"shpat_1234567890"}`
	creds, err := shopify.ParseCredentials(raw, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.Shop != "my-store.myshopify.com" {
		t.Fatalf("got shop %q, want my-store.myshopify.com", creds.Shop)
	}
	if creds.Token != "shpat_1234567890" {
		t.Fatalf("got token %q, want shpat_1234567890", creds.Token)
	}
}

func TestParseCredentials_NormalizeShop(t *testing.T) {
	raw := `{"shop":"my-test-store","token":"shpat_test"}`
	creds, err := shopify.ParseCredentials(raw, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.Shop != "my-test-store.myshopify.com" {
		t.Fatalf("got shop %q, want my-test-store.myshopify.com", creds.Shop)
	}
}

func TestParseCredentials_EnvFallback(t *testing.T) {
	t.Setenv("SHOPIFY_ACCESS_TOKEN", "shpat_env_token")
	t.Setenv("SHOPIFY_SHOP", "env-store.myshopify.com")

	creds, err := shopify.ParseCredentials("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.Token != "shpat_env_token" {
		t.Fatalf("got token %q, want shpat_env_token", creds.Token)
	}
	if creds.Shop != "env-store.myshopify.com" {
		t.Fatalf("got shop %q, want env-store.myshopify.com", creds.Shop)
	}
}

func TestParseCredentials_MissingRequired(t *testing.T) {
	_ = os.Unsetenv("SHOPIFY_ACCESS_TOKEN")
	_ = os.Unsetenv("SHOPIFY_TOKEN")
	_ = os.Unsetenv("SHOPIFY_SHOP")
	_ = os.Unsetenv("SHOPIFY_STORE")

	if _, err := shopify.ParseCredentials(`{"shop":"store.myshopify.com"}`, ""); err == nil {
		t.Fatal("expected error when token missing")
	}
	if _, err := shopify.ParseCredentials(`{"token":"shpat_xyz"}`, ""); err == nil {
		t.Fatal("expected error when shop missing")
	}
}
