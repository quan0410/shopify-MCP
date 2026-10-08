# shopify_get_page

Retrieve details and HTML content of a specific static page by ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `page_id` | yes | The numeric ID of the page |

## Cases

### Typical call

Input:

```json
{
  "page_id": "123456789"
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

### Missing `page_id`

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
      "text": "page_id is required"
    }
  ],
  "isError": true
}
```
