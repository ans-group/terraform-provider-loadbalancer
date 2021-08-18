# Load balancer Provider

Official UKFast Loadbalancer Terraform provider

## Example Usage

```hcl
provider "loadbalancer" {
  api_key = "abc"
}

resource "loadbalancer_targetgroup" "targetgroup-1" {
  cluster_id = 123
  name = "testgroup"
}
```

## Argument Reference

* `api_key`: UKFast API key - read/write permissions for `loadbalancer` service required. If omitted, will use `UKF_API_KEY` environment variable value