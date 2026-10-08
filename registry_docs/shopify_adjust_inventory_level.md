# shopify_adjust_inventory_level

Adjust available inventory quantity for an inventory item at a location by a delta amount (e.g. +5 or -2).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `inventory_item_id` | yes | Numeric inventory item ID |
| `location_id` | yes | Numeric location ID |
| `available_adjustment` | yes | Quantity adjustment delta (positive or negative integer) |

## Cases

### Typical call

Input:

```json
{
  "inventory_item_id": 1001,
  "location_id": 1001,
  "available_adjustment": 1001
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
