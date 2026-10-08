# shopify_create_draft_order

Create a draft order with line items, customer details, and notes.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `line_items` | yes | array |
| `customer` | no | object |
| `note` | no | Note for draft order |
| `tags` | no | Tags |

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
