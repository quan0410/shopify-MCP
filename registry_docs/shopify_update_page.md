# shopify_update_page

Update an existing static page's title, body HTML, or publication status.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `page_id` | yes | The numeric ID of the page to update |
| `title` | no | New title for the page |
| `body_html` | no | New HTML content for the page |
| `is_published` | no | Whether the page is published |

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
