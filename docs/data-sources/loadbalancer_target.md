# loadbalancer_target Data Source

This resource represents a target

## Example Usage

```hcl
data "loadbalancer_target" "target-1" {
  target_group_id = 123
  name            = "my-target"
}
```

## Argument Reference

- `target_group_id`: (Required) ID of target group
- `target_id`: ID of target
- `name`: Name of target
- `ip`: IP address of target
- `port`: Port number of target

## Attributes Reference

- `id`: Target ID
- `target_group_id`: ID of target group
- `name`: Name of target
- `ip`: IP address of target
- `port`: Port number of target
- `weight`: Weight of target
- `backup`: Specifies target is a backup
- `check_interval`: Check interval for target
- `check_ssl`: Specifies SSL should be used for checks
- `check_rise`: Check rise value for target
- `check_fall`: Check fall value for target
- `disable_http2`: Specifies HTTP2 is disabled for target
- `http2_only`: HTTP2 only is enabled for target
- `active`: Active status of target