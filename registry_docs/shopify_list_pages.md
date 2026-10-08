# shopify_list_pages

List static store pages (e.g. About Us, Contact, FAQ, Terms) with optional filtering.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `limit` | no | Number of pages to return (default 50, max 250) |
| `since_id` | no | Restrict results to after the specified ID |
| `title` | no | Filter pages by title |

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

