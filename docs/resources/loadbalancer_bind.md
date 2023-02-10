# loadbalancer_bind Resource

This resource is for managing loadbalancer binds

## Example Usage

```hcl
resource "loadbalancer_bind" "bind-1" {
  listener_id = 1
  ip = "1.2.3.4"
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `vip_id`: (Required) ID of VIP
- `port`: (Required) Port number for bind

## Attributes Reference

- `id`: Bind ID
- `listener_id`: ID of listener
- `vip_id`: ID of VIP
- `port`: Port number for bind

## Import

```
terraform import loadbalancer_bind.example_bind {listener_id}/{bind_id}
```