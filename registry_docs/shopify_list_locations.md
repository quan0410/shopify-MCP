# shopify_list_locations

List all inventory locations configured on the Shopify store.

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

