# loadbalancer_targetgroup Resource

This resource is for managing loadbalancer target groups

## Example Usage

```hcl
resource "loadbalancer_targetgroup" "targetgroup-1" {
  cluster_id      = 1
  name            = "group-1"
  balance         = "roundrobin"
  mode            = "tcp"
}
```

## Argument Reference

- `name`: (Required) Name of group
- `cluster_id`: (Required) ID of loadbalancer cluster
- `balance`: (Required) Balance configuration for target group
- `mode`: (Required) Mode configuration for target group
- `close`: Close configuration for target group
- `sticky`: Sticky configuration for target group
- `cookie_opts`: Cookie options for target group
- `source`: Source for target group
- `timeouts_connect`: Connect timeout for target group
- `timeouts_server`: Server timeout for target group
- `custom_options`: Custom options for target group
- `monitor_url`: Monitor URL for target group
- `monitor_method`: Monitor method for target group
- `monitor_host`: Monitor host for target group
- `monitor_http_version`: Monitor HTTP version for target group
- `monitor_expect`: Expected monitor string for target group
- `monitor_tcp_monitoring`: TCP monitoring for target group
- `check_port`: Check port for target group
- `send_proxy`: Specifies proxy protocol should be used for target group
- `send_proxy_v2`: Specifies proxy protocol v2 should be used for target group
- `ssl`: Specifies SSL should be used for target group
- `ssl_verify`: Specifies SSL verifications should be performed for target group
- `sni`: Specifies SNI should be enabled for target group

## Attributes Reference

- `id`: Target group ID
- `name`: Name of group
- `cluster_id`: ID of loadbalancer cluster
- `balance`: Balance configuration for target group
- `mode`: Mode configuration for target group
- `close`: Close configuration for target group
- `sticky`: Sticky configuration for target group
- `cookie_opts`: Cookie options for target group
- `source`: Source for target group
- `timeouts_connect`: Connect timeout for target group
- `timeouts_server`: Server timeout for target group
- `custom_options`: Custom options for target group
- `monitor_url`: Monitor URL for target group
- `monitor_method`: Monitor method for target group
- `monitor_host`: Monitor host for target group
- `monitor_http_version`: Monitor HTTP version for target group
- `monitor_expect`: Expected monitor string for target group
- `monitor_tcp_monitoring`: TCP monitoring for target group
- `check_port`: Check port for target group
- `send_proxy`: Specifies proxy protocol should be used for target group
- `send_proxy_v2`: Specifies proxy protocol v2 should be used for target group
- `ssl`: Specifies SSL should be used for target group
- `ssl_verify`: Specifies SSL verifications should be performed for target group
- `sni`: Specifies SNI should be enabled for target group