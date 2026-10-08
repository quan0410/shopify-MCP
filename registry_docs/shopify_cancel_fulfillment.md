# shopify_cancel_fulfillment

Cancel an existing fulfillment by fulfillment ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `fulfillment_id` | yes | The numeric ID of the fulfillment to cancel |

## Cases

### Typical call

Input:

```json
{
  "fulfillment_id": "123456789"
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

### Missing `fulfillment_id`

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
      "text": "fulfillment_id is required"
    }
  ],
  "isError": true
}
```
