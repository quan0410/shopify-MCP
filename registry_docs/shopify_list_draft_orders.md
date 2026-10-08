# shopify_list_draft_orders

List draft orders with optional status filter, fields, and pagination.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of draft orders to return |
| `since_id` | no | Filter draft orders after specified ID |
| `status` | no | open, invoice_sent, completed |
| `fields` | no | Comma-separated list of fields to retrieve (e.g. 'id,name,status,total_price,line_items'). |
| `safe_mode` | no | If true, queries only non-PII fields to avoid Protected Customer Data restrictions (default false). |

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

