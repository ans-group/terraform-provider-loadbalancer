# loadbalancer_acl Resource

This resource is for managing loadbalancer ACLs

## Example Usage

```hcl
resource "loadbalancer_acl" "acl-1" {
  listener_id = 1
  name        = "acl-1"
  condition {
    name = "header_matches"
    argument {
      name  = "header"
      value = "host"
    }
    argument {
      name  = "value"
      value = "ukfast.co.uk"
    }
  }
  action {
    name = "redirect"
    argument {
      name  = "location"
      value = "developers.ukfast.io"
    }
    argument {
      name  = "status"
      value = "302"
    }
  }
}
```

## Argument Reference

- `listener_id`: (Required) ID of listener. Mutually exclusive with `target_group_id`
- `target_group_id`: (Required) ID of target group. Mutually exclusive with `listener_id`
- `name`: Name of ACL
- `condition`: List of conditions
  - `name`: (Required) Name of condition
  - `argument`: List of arguments
    - `name`: (Required) Name of argument
    - `value`: (Required) Value of argument
- `action`: List of actions
  - `name`: (Required) Name of action
  - `argument`: List of arguments
    - `name`: (Required) Name of argument
    - `value`: (Required) Value of argument

## Attributes Reference

- `id`: ACL ID
- `listener_id`: ID of listener
- `target_group_id`: ID of target group
- `name`: Name of ACL
- `condition`: List of conditions
  - `name`: Name of condition
  - `argument`: List of arguments
    - `name`: Name of argument
    - `value`: Value of argument
- `action`: List of actions
  - `name`: Name of action
  - `argument`: List of arguments
    - `name`: Name of argument
    - `value`: Value of argument