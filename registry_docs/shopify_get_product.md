# shopify_get_product

Retrieve full details of a specific product by product ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `product_id` | yes | The numeric ID of the product |

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
