package shopify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// DefaultAPIVersion is the supported Shopify Admin API version.
const DefaultAPIVersion = "2024-04"

// Client calls the Shopify Admin REST API, authenticated with the connecting user's
// Admin API access token.
type Client struct {
	creds      *Credentials
	httpClient *http.Client
	apiVersion string
	apiBase    string // optional override (e.g. for testing)
}

func NewClient(creds *Credentials) *Client {
	timeoutSec := 30
	if v := strings.TrimSpace(os.Getenv("SHOPIFY_HTTP_TIMEOUT_SEC")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutSec = n
		}
	}
	version := strings.TrimSpace(os.Getenv("SHOPIFY_API_VERSION"))
	if version == "" {
		version = DefaultAPIVersion
	}
	base := strings.TrimSpace(os.Getenv("SHOPIFY_API_BASE"))
	return &Client{
		creds:      creds,
		httpClient: &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
		apiVersion: version,
		apiBase:    strings.TrimRight(base, "/"),
	}
}

func (c *Client) adminBase() string {
	if c.apiBase != "" {
		return fmt.Sprintf("%s/admin/api/%s", c.apiBase, c.apiVersion)
	}
	return fmt.Sprintf("https://%s/admin/api/%s", c.creds.Shop, c.apiVersion)
}

// Request executes an authenticated HTTP request against Shopify Admin REST API.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, int, error) {
	if c == nil || c.creds == nil {
		return nil, 0, fmt.Errorf("shopify client not configured")
	}

	base := c.adminBase()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	full := base + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("encode request body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, full, reader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-Shopify-Access-Token", c.creds.Token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("shopify request failed: %w", err)
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, 10<<20) // 10 MiB cap
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := sanitizeProviderError(string(data))
		return data, resp.StatusCode, fmt.Errorf("shopify HTTP %d: %s", resp.StatusCode, msg)
	}
	return data, resp.StatusCode, nil
}

// GraphQL executes an authenticated GraphQL query or mutation against the Shopify Admin API.
func (c *Client) GraphQL(ctx context.Context, query string, variables map[string]interface{}) ([]byte, int, error) {
	if c == nil || c.creds == nil {
		return nil, 0, fmt.Errorf("shopify client not configured")
	}

	payload := map[string]interface{}{
		"query": query,
	}
	if len(variables) > 0 {
		payload["variables"] = variables
	}

	return c.Request(ctx, http.MethodPost, "/graphql.json", nil, payload)
}

// TokenMasked returns a redacted access token preview.
func (c *Client) TokenMasked() string {
	if c == nil || c.creds == nil {
		return "(missing)"
	}
	return MaskToken(c.creds.Token)
}

// Shop returns the connected Shopify store domain.
func (c *Client) Shop() string {
	if c == nil || c.creds == nil {
		return ""
	}
	return c.creds.Shop
}

func sanitizeProviderError(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return "(empty body)"
	}
	if len(body) > 500 {
		body = body[:500] + "…"
	}
	lower := strings.ToLower(body)
	if strings.Contains(lower, "shpat_") || strings.Contains(lower, "token") {
		return "(provider error sanitized)"
	}
	return body
}
