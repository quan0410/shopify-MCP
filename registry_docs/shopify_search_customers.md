# shopify_search_customers

Search for customers matching a query string (e.g. name, email, phone number, address).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `query` | yes | Search query text |
| `limit` | no | Maximum number of results to return |

## Cases

### Typical call

Input:

```json
{
  "query": "{ shop { name } }"
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

### Missing `query`

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
      "text": "query is required"
    }
  ],
  "isError": true
}
```
