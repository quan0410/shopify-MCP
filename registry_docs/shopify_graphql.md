# shopify_graphql

Execute an arbitrary GraphQL query or mutation against Shopify Admin API (2024-04). Recommended for fetching orders or customers with specific non-PII fields to bypass Protected Customer Data restrictions.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `query` | yes | The GraphQL query or mutation string |
| `variables` | no | Optional variables map for the GraphQL query |

## Cases

### Typical call

Input:

```json
{
  "query": "{ shop { name } }"
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

### Missing `query`

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
      "text": "query is required"
    }
  ],
  "isError": true
}
```
