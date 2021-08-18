# loadbalancer_cluster Data Source

This resource represents a loadbalancer cluster

## Example Usage

```hcl
data "loadbalancer_cluster" "cluster-1" {
  name = "my-cluster"
}
```

## Argument Reference

- `cluster_id`: ID of loadbalancer cluster
- `name`: Name of loadbalancer cluster
- `deployed`: Deployment status loadbalancer cluster

## Attributes Reference

- `id`: Cluster ID
- `name`: Name of loadbalancer cluster
- `deployed`: Deployment status loadbalancer cluster