# shopify_get_shop

Retrieve configuration and metadata for the authenticated Shopify store (name, domain, currency, timezone, email).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| *(none)* | - | This tool does not require additional arguments. |

## Cases

### Typical call

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
      "text": "{\"status\": \"success\"}"
    }
  ],
  "isError": false
}
```

