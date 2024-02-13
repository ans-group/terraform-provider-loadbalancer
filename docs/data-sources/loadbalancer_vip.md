# loadbalancer_vip Data Source

This resource represents a loadbalancer vip

## Example Usage

```hcl
data "loadbalancer_vip" "vip-1" {
  vip_id = 1
}
```

## Argument Reference

- `vip_id`: ID of loadbalancer vip
- `internal_cidr`: Interal CIDR of loadbalancer vip
- `external_cidr`: External CIDR of loadbalancer vip
- `mac_address`: MAC address of loadbalancer vip

## Attributes Reference

- `id`: ID of loadbalancer vip
- `internal_cidr`: Interal CIDR of loadbalancer vip
- `external_cidr`: External CIDR of loadbalancer vip
- `mac_address`: MAC address of loadbalancer vip