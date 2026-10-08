# shopify_update_customer

Update customer details (name, email, phone, tags, notes) by customer ID.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `customer_id` | yes | Customer ID to update |
| `first_name` | no | Updated first name |
| `last_name` | no | Updated last name |
| `email` | no | Updated email |
| `phone` | no | Updated phone number |
| `tags` | no | Updated tags |
| `note` | no | Updated notes |

## Cases

### Typical call

Input:

```json
{
  "customer_id": "123456789"
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

### Missing `customer_id`

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
      "text": "customer_id is required"
    }
  ],
  "isError": true
}
```
