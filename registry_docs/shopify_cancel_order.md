# shopify_cancel_order

Cancel an existing order with optional reason and refund flag.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `order_id` | yes | Order ID to cancel |
| `reason` | no | Reason for cancellation: customer, inventory, fraud, declined, other |
| `email` | no | Whether to send an email to the customer (default true) |

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
