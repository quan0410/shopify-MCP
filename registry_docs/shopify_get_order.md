# shopify_get_order

Retrieve single order details by order ID. Use 'fields' or 'safe_mode' to exclude customer/address PII and avoid Shopify 403 Protected Customer Data errors.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `order_id` | yes | Order ID |
| `fields` | no | Comma-separated list of fields to retrieve (e.g. 'id,name,order_number,created_at,financial_status,fulfillment_status,total_price,currency,line_items'). |
| `safe_mode` | no | If true, queries only non-PII fields to avoid Protected Customer Data restrictions (default false). |

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
