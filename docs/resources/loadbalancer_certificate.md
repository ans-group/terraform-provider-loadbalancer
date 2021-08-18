# loadbalancer_certificate Resource

This resource is for managing loadbalancer certificates

## Example Usage

```hcl
resource "loadbalancer_certificate" "certificate-1" {
  listener_id = 1
  name        = "somecertificate"
  key         = file("${path.module}/cert.key")
  certificate = file("${path.module}/cert.crt")
  ca_bundle   = file("${path.module}/ca.crt")
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `name`: Name of certificate
- `key`: Private key for certificate
- `certificate`: Certificate contents
- `ca_bundle`: CA bundle contents

## Attributes Reference

- `id`: Certificate ID
- `listener_id`: ID of listener
- `name`: Name of certificate