# shopify_create_page

Create a new static store page with title and body HTML.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `title` | yes | Title of the page |
| `body_html` | no | Page content in HTML or plain text |
| `author` | no | Name of the page author (optional) |
| `is_published` | no | Whether the page is published immediately (default true) |

## Cases

### Typical call

Input:

```json
{
  "title": "sample_value"
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

### Missing `title`

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
      "text": "title is required"
    }
  ],
  "isError": true
}
```
