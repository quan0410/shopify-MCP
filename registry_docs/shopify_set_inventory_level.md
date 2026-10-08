# shopify_set_inventory_level

Set exact available inventory quantity for an inventory item at a location.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `inventory_item_id` | yes | Numeric inventory item ID |
| `location_id` | yes | Numeric location ID |
| `available` | yes | Exact available quantity to set |

## Cases

### Typical call

Input:

```json
{
  "inventory_item_id": 1001,
  "location_id": 1001,
  "available": 1001
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

### Missing `inventory_item_id`

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
      "text": "inventory_item_id is required"
    }
  ],
  "isError": true
}
```
