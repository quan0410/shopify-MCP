# shopify_lookup_discount_code

Look up discount code details and its associated price rule by code string.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `code` | yes | The discount code text to search for |

## Cases

### Typical call

Input:

```json
{
  "code": "SUMMER2024"
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

### Missing `code`

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
      "text": "code is required"
    }
  ],
  "isError": true
}
```
