# loadbalancer_listener Data Source

This resource represents a listener

## Example Usage

```hcl
data "loadbalancer_listener" "listener-1" {
  listener_id = 1
  name = "somelistener"
}
```

## Argument Reference

- `listener_id`: ID of listener
- `name`: Name of listener
- `cluster_id`: ID of cluster

## Attributes Reference

- `id`: Listener ID
- `name`: Name of listener
- `cluster_id`: ID of cluster
- `hsts_enabled`: Specifies HSTS is enabled
- `mode`: Mode of listener
- `hsts_maxage`: HSTS max age
- `close`: Specifies whether keepalive is disabled
- `redirect_https`: Specifies HTTPS redirection is enabled
- `default_target_group_id`: Specifies default target group ID
- `allow_tlsv1`: Specifies TLS 1.0 is enabled
- `allow_tlsv11`: Specifies TLS 1.1 is enabled
- `disable_tlsv12`: Specifies TLS 1.2 is disabled
- `disable_http2`: Specifies HTTP2 is disabled
- `http2_only`: Specifies only HTTP2 is enabled