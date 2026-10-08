# shopify_create_article

Create a new article/blog post under a specific blog.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `blog_id` | yes | The numeric ID of the blog |
| `title` | yes | Article title |
| `body_html` | yes | Article HTML or text body |
| `author` | no | Author name |
| `tags` | no | Comma-separated list of tags |
| `summary_html` | no | Article summary/excerpt (optional) |
| `is_published` | no | Whether to publish the article immediately (default true) |

## Cases

### Typical call

Input:

```json
{
  "blog_id": "123456789",
  "title": "sample_value",
  "body_html": "sample_value"
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
