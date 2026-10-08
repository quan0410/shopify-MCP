# shopify_create_product

Create a new product in Shopify with title, description, vendor, product_type, tags, and optional variants.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `title` | yes | Product title |
| `body_html` | no | Product description in HTML or plain text |
| `vendor` | no | Product vendor name |
| `product_type` | no | Category/type of product |
| `status` | no | active, draft, or archived (default active) |
| `tags` | no | Comma-separated list of tags |
| `variants` | no | List of variant objects (price, sku, title, etc.) |

## Cases

### Typical call

Input:

```json
{
  "title": "sample_value"
}
```

Output:

```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"status\": \"success\"}"
    }
  ],
  "isError": false
}
```

### Missing `title`

The tool rejects the call and does not guess the missing value.

Input:

```json
{}
```

Output:

```json
{
  "content": [
    {
      "type": "text",
      "text": "title is required"
    }
  ],
  "isError": true
}
```
