# shopify_create_customer

Create a new customer profile with name, email, phone, and optional tags.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `email` | yes | Customer's email address |
| `first_name` | no | Customer's first name |
| `last_name` | no | Customer's last name |
| `phone` | no | Customer's phone number |
| `tags` | no | Comma separated tags |
| `note` | no | Note about the customer |

## Cases

### Typical call

Input:

```json
{
  "email": "customer@example.com"
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

### Missing `email`

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
      "text": "email is required"
    }
  ],
  "isError": true
}
```
