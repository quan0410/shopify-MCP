# shopify_update_variant

Update a product variant's price, sku, barcode, or inventory policy.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `variant_id` | yes | Variant ID to update |
| `price` | no | New price |
| `sku` | no | New SKU |
| `barcode` | no | New barcode |
| `option1` | no | Option 1 value (e.g. Size) |
| `option2` | no | Option 2 value (e.g. Color) |

## Cases

### Typical call

Input:

```json
{
  "variant_id": "123456789"
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

### Missing `variant_id`

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
      "text": "variant_id is required"
    }
  ],
  "isError": true
}
```
