# shopify_list_price_rules

List discount price rules with optional pagination and status filter.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of price rules to return (default 50, max 250) |
| `since_id` | no | Restrict results to after the specified ID |

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

