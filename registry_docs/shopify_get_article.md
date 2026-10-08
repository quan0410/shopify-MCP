# shopify_get_article

Retrieve full details of a specific blog article by blog ID and article ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `blog_id` | yes | The numeric ID of the blog |
| `article_id` | yes | The numeric ID of the article |

## Cases

### Typical call

Input:

```json
{
  "blog_id": "123456789",
  "article_id": "123456789"
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

### Missing `blog_id`

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
      "text": "blog_id is required"
    }
  ],
  "isError": true
}
```
