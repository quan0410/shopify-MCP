# shopify_update_product

Update an existing product's fields (title, body_html, vendor, status, tags).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `product_id` | yes | Product ID to update |
| `title` | no | New title |
| `body_html` | no | New description |
| `vendor` | no | New vendor |
| `product_type` | no | New product type |
| `status` | no | New status (active, draft, archived) |
| `tags` | no | New tags (comma separated) |

## Cases

### Typical call

Input:

```json
{
  "product_id": "123456789"
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

### Missing `product_id`

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
      "text": "product_id is required"
    }
  ],
  "isError": true
}
```
