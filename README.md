# terraform-provider-loadbalancer

## Getting Started

This provider is available via the [Terraform Registry](https://registry.terraform.io/providers/ans-group/loadbalancer/latest) with Terraform v0.13+

> :warning: We strongly recommend pinning the provider version to a target major version, as to ensure future breaking changes do not affect workflows and automated CI pipelines

```
terraform {
  required_providers {
    loadbalancer = {
      source  = "ans-group/loadbalancer"
      version = "~> 1.0"
    }
  }
}

provider "loadbalancer" {
  api_key = "abc"
}
```
