# shopify_list_orders

List orders from the store with filters for status, financial_status, fulfillment_status, created_at range, and pagination. To avoid Shopify 403 Forbidden / Protected Customer Data errors when the app lacks PCD approval, specify 'fields' (e.g. 'id,name,created_at,financial_status,fulfillment_status,total_price,currency,line_items') or set safe_mode to true.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of orders to retrieve (default 50, max 250) |
| `since_id` | no | Restrict results to after the specified order ID |
| `status` | no | Filter by status: open, closed, cancelled, any (default open) |
| `financial_status` | no | Filter by financial status: authorized, pending, paid, refunded, voided, any |
| `fulfillment_status` | no | Filter by fulfillment status: shipped, partial, unshipped, any |
| `created_at_min` | no | Filter orders created at or after date-time (ISO 8601, e.g. '2024-05-01T00:00:00Z') |
| `created_at_max` | no | Filter orders created at or before date-time (ISO 8601, e.g. '2024-05-31T23:59:59Z') |
| `fields` | no | Comma-separated list of fields to retrieve (e.g. 'id,name,order_number,created_at,financial_status,fulfillment_status,total_price,currency,line_items'). Excludes customer/address PII to prevent Shopify 403 Protected Customer Data errors. |
| `safe_mode` | no | If true, automatically queries only non-PII fields to avoid Protected Customer Data restrictions (default false). |

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

