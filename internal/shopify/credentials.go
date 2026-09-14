package shopify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Credentials matches datumbridge-integrations' ShopifyTokenBundle
// (internal/credentials/vault.go).
type Credentials struct {
	Type         string   `json:"type"`
	Shop         string   `json:"shop"`
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	TokenURI     string   `json:"token_uri,omitempty"`
}

// NormalizeShop normalizes a shop string into standard "example.myshopify.com" host.
func NormalizeShop(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "://") {
		if u, err := url.Parse(s); err == nil && u.Host != "" {
			s = u.Host
		}
	}
	s = strings.TrimRight(s, "/")
	if !strings.Contains(s, ".") {
		s = s + ".myshopify.com"
	}
	return strings.ToLower(s)
}

// ParseCredentials loads Credentials from vault args (credentials_json /
// credentials_path), then falls back to process env for local/dev runs.
// Precedence: injected bundle first, SHOPIFY_ACCESS_TOKEN/SHOPIFY_SHOP env fills blanks.
func ParseCredentials(credentialsJSON, credentialsPath string) (*Credentials, error) {
	c := &Credentials{}
	raw := strings.TrimSpace(credentialsJSON)

	if raw == "" && strings.TrimSpace(credentialsPath) != "" {
		path := filepath.Clean(credentialsPath)
		if strings.Contains(path, "..") {
			return nil, fmt.Errorf("credentials_path must not contain ..")
		}
		allowedRoot := strings.TrimSpace(os.Getenv("SHOPIFY_CREDENTIALS_DIR"))
		if allowedRoot == "" {
			allowedRoot = "/credentials"
		}
		allowedRoot = filepath.Clean(allowedRoot)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		absRoot, err := filepath.Abs(allowedRoot)
		if err != nil {
			return nil, fmt.Errorf("SHOPIFY_CREDENTIALS_DIR: %w", err)
		}
		resolvedRoot, err := filepath.EvalSymlinks(absRoot)
		if err != nil {
			resolvedRoot = absRoot
		}
		resolvedPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		if resolvedPath != resolvedRoot && !strings.HasPrefix(resolvedPath, resolvedRoot+string(os.PathSeparator)) {
			return nil, fmt.Errorf("credentials_path must be under %s", absRoot)
		}
		b, err := os.ReadFile(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("read credentials_path: %w", err)
		}
		raw = strings.TrimSpace(string(b))
	}

	if raw != "" {
		if err := json.Unmarshal([]byte(raw), c); err != nil {
			return nil, fmt.Errorf("invalid credentials_json: %w", err)
		}
	}

	if strings.TrimSpace(c.Token) == "" {
		if envToken := os.Getenv("SHOPIFY_ACCESS_TOKEN"); envToken != "" {
			c.Token = envToken
		} else {
			c.Token = os.Getenv("SHOPIFY_TOKEN")
		}
	}
	if strings.TrimSpace(c.Shop) == "" {
		if envShop := os.Getenv("SHOPIFY_SHOP"); envShop != "" {
			c.Shop = envShop
		} else {
			c.Shop = os.Getenv("SHOPIFY_STORE")
		}
	}

	c.Token = strings.TrimSpace(c.Token)
	c.Shop = NormalizeShop(c.Shop)

	if c.Token == "" {
		return nil, fmt.Errorf("token required (credentials_json.token or SHOPIFY_ACCESS_TOKEN)")
	}
	if c.Shop == "" {
		return nil, fmt.Errorf("shop required (credentials_json.shop or SHOPIFY_SHOP)")
	}
	return c, nil
}

// MaskToken returns a redacted preview for health/debug output.
func MaskToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "(missing)"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "…" + token[len(token)-4:]
}
