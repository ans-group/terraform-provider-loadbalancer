# loadbalancer_bind Data Source

This resource represents a bind

## Example Usage

```hcl
data "loadbalancer_bind" "bind-1" {
  listener_id = 1
  port = 80
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `bind_id`: ID of bind
- `vip_id`: ID of VIP
- `port`: Port number for bind

## Attributes Reference

- `id`: Bind ID
- `listener_id`: ID of listener
- `vip_id`: ID of VIP
- `port`: Port number for bind