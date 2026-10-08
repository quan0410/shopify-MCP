# shopify_list_fulfillments

List all fulfillments for a specific order by order ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `order_id` | yes | The numeric ID of the order |

## Cases

### Typical call

Input:

```json
{
  "order_id": "123456789"
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

### Missing `order_id`

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
      "text": "order_id is required"
    }
  ],
  "isError": true
}
```
