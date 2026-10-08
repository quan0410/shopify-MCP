# shopify_list_inventory_levels

Retrieve inventory levels for specific inventory_item_ids or location_ids. Shopify requires at least one of location_ids or inventory_item_ids; if neither is provided, active store locations are automatically queried.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `location_ids` | no | Comma-separated list or array of location IDs. Required by Shopify unless inventory_item_ids is provided (auto-detected from active store locations if omitted). |
| `inventory_item_ids` | no | Comma-separated list or array of inventory item IDs. Required by Shopify unless location_ids is provided. |
| `limit` | no | Number of inventory levels to return (default 50, max 250) |

## Cases

### Typical call

Input:

```json
{
  "location_ids": "example",
  "inventory_item_ids": "example"
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

