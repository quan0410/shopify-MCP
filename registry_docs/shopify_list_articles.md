# shopify_list_articles

List articles/blog posts published in a specific blog.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `blog_id` | yes | The numeric ID of the blog |
| `limit` | no | Number of articles to return (default 50) |
| `since_id` | no | Restrict results to after the specified ID |

## Cases

### Typical call

Input:

```json
{
  "blog_id": "123456789"
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
