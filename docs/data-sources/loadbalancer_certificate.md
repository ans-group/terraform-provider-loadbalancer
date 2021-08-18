# loadbalancer_certificate Data Source

This resource represents a certificate

## Example Usage

```hcl
data "loadbalancer_certificate" "certificate-1" {
  listener_id 1
  name = "my-certificate"
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `certificate_id`: ID of certificate
- `name`: Name of certificate

## Attributes Reference

- `id`: Certificate ID
- `listener_id`: ID of listener
- `name`: Name of certificate