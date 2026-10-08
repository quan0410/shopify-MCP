# shopify_create_price_rule

Create a new discount price rule (percentage, fixed_amount, or free_shipping).

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `title` | yes | Title of the price rule |
| `target_type` | yes | line_item or shipping_line |
| `target_selection` | yes | all or entitled |
| `allocation_method` | yes | across or each |
| `value_type` | yes | percentage or fixed_amount |
| `value` | yes | Discount value (negative number, e.g. -15.0) |
| `customer_selection` | yes | all or prerequisite |
| `starts_at` | yes | ISO 8601 start date and time |
| `ends_at` | no | ISO 8601 end date and time (optional) |

## Cases

### Typical call

Input:

```json
{
  "title": "sample_value",
  "target_type": "sample_value",
  "target_selection": "sample_value",
  "allocation_method": "sample_value",
  "value_type": "sample_value",
  "value": "sample_value",
  "customer_selection": "sample_value",
  "starts_at": "sample_value"
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

### Missing `title`

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
      "text": "title is required"
    }
  ],
  "isError": true
}
```
