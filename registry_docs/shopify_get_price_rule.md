# shopify_get_price_rule

Retrieve details of a specific price rule by ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `price_rule_id` | yes | The numeric ID of the price rule |

## Cases

### Typical call

Input:

```json
{
  "price_rule_id": "123456789"
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

### Missing `price_rule_id`

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
      "text": "price_rule_id is required"
    }
  ],
  "isError": true
}
```
