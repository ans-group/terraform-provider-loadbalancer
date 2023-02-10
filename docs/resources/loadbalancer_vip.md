# loadbalancer_vip Resource

This resource is for managing loadbalancer VIPs

Note - it is not currently possible to create/update/delete VIPs via this resource. An existing `cluster_id` must be provided

## Example Usage

```hcl
resource "loadbalancer_vip" "vip-1" {
  cluster_id = 1
  internal_cidr = "10.0.0.5/30"
  external_cidr = "10.0.0.5/30"
}
```

## Argument Reference

- `cluster_id`: (Required) ID of cluster
- `internal_cidr`: (Required) Internal CIDR
- `external_cidr`: (Required) External CIDR

## Attributes Reference

- `id`: Cluster ID
- `cluster_id`: (Required) ID of cluster
- `internal_cidr`: Internal CIDR
- `external_cidr`: External CIDR
- `mac_address`: MAC Address of VIP