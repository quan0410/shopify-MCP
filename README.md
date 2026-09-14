# shopify-mcp

Shopify MCP tool-server for DatumBridge Agent Kit. Exposes Shopify Admin REST API resources as MCP tools over Streamable HTTP (`POST /mcp`), authenticated by credential bundles injected from `datumbridge-integrations`.

## Features

- **29 MCP Tools** spanning Products, Variants, Orders, Draft Orders, Customers, Inventory Levels, Locations, Collections, and Store Info.
- **Streamable HTTP MCP** transport (protocol version `2024-11-05`) with session management (`Mcp-Session-Id`).
- **Standard Credential Contract**: Consumes injected `credentials_json` containing `shop` and `token`.
- **Zero-Storage Auth**: Never stores tokens or performs OAuth redirects directly.
- **Configurable Tool Filtering**: Enable specific tools via `SHOPIFY_TOOLS` env var or expose full scope by default.

## Running Locally

```bash
cp .env.example .env
# Edit .env to set SHOPIFY_ACCESS_TOKEN and SHOPIFY_SHOP for standalone local testing
go run ./cmd/api
```

## Running Tests

```bash
go test -v ./...
```
