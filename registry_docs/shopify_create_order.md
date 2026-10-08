# shopify_create_order

Create a new order in Shopify with line items, customer information, shipping/billing address, and financial status.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `line_items` | yes | array |
| `customer` | no | object |
| `billing_address` | no | Billing address object |
| `shipping_address` | no | Shipping address object |
| `financial_status` | no | paid, pending, authorized, etc. |
| `note` | no | Order notes |
| `tags` | no | Comma separated tags |

## Cases

### Typical call

Input:

```json
{
  "line_items": [
    {"variant_id": 12345, "quantity": 1, "price": "19.99"}
  ]
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

### Missing `line_items`

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
      "text": "line_items is required"
    }
  ],
  "isError": true
}
```
