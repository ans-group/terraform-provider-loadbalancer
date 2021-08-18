# loadbalancer_target Resource

This resource is for managing loadbalancer targets

## Example Usage

```hcl
resource "loadbalancer_target" "target-1" {
  target_group_id = 1
  name            = "target-1"
  ip              = "1.2.3.4"
  port            = 80
}
```

## Argument Reference

- `name`: (Required) Name of target
- `target_group_id`: (Required) ID of target group
- `ip`: (Required) IP address of target
- `port`: (Required) Port number of target
- `weight`: Weight of target
- `backup`: Specifies target is a backup
- `check_interval`: Check interval for target
- `check_ssl`: Specifies SSL should be used for checks
- `check_rise`: Check rise value for target
- `check_fall`: Check fall value for target
- `disable_http2`: Specifies HTTP2 is disabled for target
- `http2_only`: HTTP2 only is enabled for target
- `active`: Active status of target

## Attributes Reference

- `id`: Target ID
- `name`: Name of target
- `target_group_id`: ID of target group
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