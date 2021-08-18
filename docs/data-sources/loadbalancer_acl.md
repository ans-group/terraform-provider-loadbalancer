# loadbalancer_acl Data Source

This resource represents an ACL

## Example Usage

```hcl
data "loadbalancer_acl" "acl-1" {
  listener_id = 1
  name        = "testacl"
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener. Mutually exclusive with `target_group_id`
- `target_group_id`: (Required) ID of target group. Mutually exclusive with `listener_id`
- `acl_id`: ID of ACL
- `name`: Name of ACL

## Attributes Reference

- `id`: ACL ID
- `name`: Name of ACL
- `listener_id`: ID of listener
- `target_group_id`: ID of target group