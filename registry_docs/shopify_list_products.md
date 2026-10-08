# shopify_list_products

List products with optional pagination and filtering by title, status, vendor, or product_type.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of products to return (default 50, max 250) |
| `since_id` | no | Restrict results to after the specified product ID |
| `title` | no | Filter products by title |
| `vendor` | no | Filter products by vendor |
| `product_type` | no | Filter products by product type |
| `status` | no | Filter by status: active, archived, draft |

## Cases

### Typical call

Input:

```json
{
  "limit": 50,
  "since_id": "example"
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

