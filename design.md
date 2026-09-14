# Design — shopify-mcp

## Class

**Tool-server** (Streamable HTTP MCP), matching `Bright-Data-MCP` and `jira-mcp`'s shape.

## Language

Go 1.23, matching the platform standard established in `doc/MCP_TOOL_REFERENCE_ARCHITECTURE.md` (BR1).

## Upstream

- Shopify Admin REST API (`https://{shop}/admin/api/2024-04`)
- `{shop}` comes directly from the injected credential bundle (`Credentials.Shop`), never from tool arguments (preventing multi-tenant confusion).

## Auth

- Production: Vault `credentials_json` injected per tool call via `POST /internal/v1/inject` on `datumbridge-integrations`.
- Fallback: `SHOPIFY_ACCESS_TOKEN` and `SHOPIFY_SHOP` env vars (local dev only).
- `credentials_path` under `SHOPIFY_CREDENTIALS_DIR` supported as a path-jailed alternative.
- Refresh token handling: Managed upstream by `datumbridge-integrations`. The tool-server uses the token in `X-Shopify-Access-Token` directly.

## Tools

29 tools across 6 categories:
1. **Products & Variants** (`internal/tools/product_tools.go`): 8 tools
2. **Orders & Draft Orders** (`internal/tools/order_tools.go`): 7 tools
3. **Customers** (`internal/tools/customer_tools.go`): 5 tools
4. **Inventory & Locations** (`internal/tools/inventory_tools.go`): 4 tools
5. **Collections** (`internal/tools/collection_tools.go`): 4 tools
6. **Shop** (`internal/tools/shop_tools.go`): 1 tool

## Security Baseline

- Never log secrets or tokens.
- Mask token in logs and health/status checks.
- Enforce shop domain pinning from bundle.
- 10 MiB response limit on upstream reads.
