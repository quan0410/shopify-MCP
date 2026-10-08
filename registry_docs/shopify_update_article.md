# shopify_update_article

Update an existing blog article's title, body, author, tags, summary, or publication status (is_published).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `article_id` | yes | The numeric ID of the article to update (e.g. '622314520859' or 'articles/622314520859') |
| `blog_id` | no | The numeric ID of the blog (optional) |
| `title` | no | New title for the article |
| `author` | no | Author name |
| `body_html` | no | New HTML or text content |
| `tags` | no | Comma-separated list of tags |
| `summary_html` | no | Article summary/excerpt |
| `is_published` | no | Whether the article is published (true to publish, false to unpublish/draft) |

## Cases

### Typical call

Input:

```json
{
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

### Missing `article_id`

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
      "text": "article_id is required"
    }
  ],
  "isError": true
}
```
