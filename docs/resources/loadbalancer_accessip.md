# loadbalancer_accessip Resource

This resource is for managing loadbalancer access IPs

## Example Usage

```hcl
resource "loadbalancer_accessip" "accessip-1" {
  listener_id = 1
  ip = "1.2.3.4"
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `ip`: (Required) IP address of access IP

## Attributes Reference

- `id`: Access IP ID
- `listener_id`: (Required) ID of listener
- `ip`: IP address of access IP