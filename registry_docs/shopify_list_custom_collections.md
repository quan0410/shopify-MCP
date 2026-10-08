# shopify_list_custom_collections

List custom (manually-curated) product collections from Shopify.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of collections to return (default 50) |
| `since_id` | no | Restrict results to after specified collection ID |
| `title` | no | Filter collections by title |

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

