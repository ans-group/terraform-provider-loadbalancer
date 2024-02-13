# loadbalancer_cluster Resource

This resource is for managing loadbalancer clusters

Note - it is not currently possible to create clusters via this resource. An existing `cluster_id` must be provided. We recommend that you use the `loadbalancer_cluster` datasource instead.

## Example Usage

```hcl
resource "loadbalancer_cluster" "cluster-1" {
  cluster_id = 1
  name = "somecluster"
}
```

## Argument Reference

- `cluster_id`: (Required) ID of cluster
- `name`: (Required) Name of cluster

## Attributes Reference

- `id`: Cluster ID
- `name`: Name of loadbalancer cluster