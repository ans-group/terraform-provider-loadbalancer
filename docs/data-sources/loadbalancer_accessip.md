# loadbalancer_accessip Data Source

This resource represents a loadbalancer access IP

## Example Usage

```hcl
data "loadbalancer_accessip" "accessip-1" {
  listener_id = 1
  ip = "1.2.3.4"
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener
- `access_ip_id`: ID of access IP
- `ip`: IP address of access IP

## Attributes Reference

- `id` Access IP ID
- `listener_id`: (Required) ID of listener
- `ip`: IP address of access IP