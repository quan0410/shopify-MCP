# shopify_list_smart_collections

List automated (smart) product collections with rule conditions.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of collections to return (default 50) |
| `title` | no | Filter smart collections by title |

## Cases

### Typical call

Input:

```json
{
  "limit": 50,
  "title": "example"
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

