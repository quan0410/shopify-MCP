# shopify_create_fulfillment

Create a fulfillment for an order or fulfillment order with tracking number and company.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `fulfillment` | yes | Fulfillment payload object (e.g. line_items_by_fulfillment_order, tracking_info, notify_customer) |

## Cases

### Typical call

Input:

```json
{
  "fulfillment": {
    "note": "sample_data"
  }
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

### Missing `fulfillment`

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
      "text": "fulfillment is required"
    }
  ],
  "isError": true
}
```
