# loadbalancer_listener Resource

This resource is for managing loadbalancer listeners

## Example Usage

```hcl
resource "loadbalancer_listener" "listener-1" {
  cluster_id = 1
  name = "listener1"
  mode = "tcp"
  default_target_group_id = 2
}
```

## Argument Reference

- `name`: (Required) Name of listener
- `cluster_id`: (Required) ID of cluster
- `mode`: (Required) Mode of listener
- `default_target_group_id`: (Required) Specifies default target group ID
- `hsts_enabled`: Specifies HSTS is enabled
- `hsts_maxage`: HSTS max age
- `close`: Specifies whether keepalive is disabled
- `redirect_https`: Specifies HTTPS redirection is enabled
- `allow_tlsv1`: Specifies TLS 1.0 is enabled
- `allow_tlsv11`: Specifies TLS 1.1 is enabled
- `disable_tlsv12`: Specifies TLS 1.2 is disabled
- `disable_http2`: Specifies HTTP2 is disabled
- `http2_only`: Specifies only HTTP2 is enabled

## Attributes Reference

- `id`: Listener ID
- `name`: Name of listener
- `cluster_id`: ID of cluster
- `mode`: Mode of listener
- `default_target_group_id`: Specifies default target group ID
- `hsts_enabled`: Specifies HSTS is enabled
- `hsts_maxage`: HSTS max age
- `close`: Specifies whether keepalive is disabled
- `redirect_https`: Specifies HTTPS redirection is enabled
- `allow_tlsv1`: Specifies TLS 1.0 is enabled
- `allow_tlsv11`: Specifies TLS 1.1 is enabled
- `disable_tlsv12`: Specifies TLS 1.2 is disabled
- `disable_http2`: Specifies HTTP2 is disabled
- `http2_only`: Specifies only HTTP2 is enabled