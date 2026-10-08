# shopify_create_custom_collection

Create a new custom collection with title, description, and optional image.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `title` | yes | Collection title |
| `body_html` | no | Collection description |
| `published` | no | Whether collection is published (default true) |

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
