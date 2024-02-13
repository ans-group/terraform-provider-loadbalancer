# Load balancer Provider

Official ANS Loadbalancer Terraform provider

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

* `api_key`: ANS API key - read/write permissions for `loadbalancer` service required. If omitted, will use `UKF_API_KEY` environment variable value

## Deployments

When you use this Terraform provider, the changes you make to the loadbalancer via Terraform are 'staged', and not deployed to the loadbalancer automatically. This allows you to make changes and switch over to your new configuration atomically. Unfortunately Terraform doesn't support the concept of deployments/commits, so deployments need to be handled externally to Terraform. There are two ways to handle this:

### Manually

To deploy your loadbalancer manually after applying changes via Terraform, you can manually perform the deployment by logging into ANS Glass, clicking Services -> Servers -> Load Balancers -> Deployments -> Deploy Now.

You can also manually deploy using the ANS CLI tool with `ans loadbalancer cluster deploy $CLUSTERID`. Find your cluster ID with `ans loadbalancer cluster list`.

### ANS CLI

The ANS CLI has a wrapper around Terraform/OpenTofu that allows you to apply the Terraform and deploy in a single step. To use this, you must provide an output of the cluster IDs your Terraform is targeting, like so:

```hcl
resource "loadbalancer_cluster" "mycluster" {
  cluster_id = 12345
  name       = "ACME Corp LB 1"
}

output "loadbalancer_cluster_ids" {
  value = [loadbalancer_cluster.mycluster.cluster_id]
}
```

The output name must be `loadbalancer_cluster_ids` and should contain an array of cluster IDs, even if there is only one.

You can then use the ANS CLI tool to plan and apply your Terraform, then deploy your loadbalancer if the Terraform apply was successful.

```shell
$ ans loadbalancer terraform plan -out mycluster.plan
...snip...
$ ans loadbalancer terraform apply mycluster.plan
...snip...

Apply complete! Resources: 0 added, 2 changed, 0 destroyed.

Outputs:

loadbalancer_cluster_ids = [
  12345,
]

Terraform run complete, deploying the configuration to the loadbalancer...
Deploying cluster(s): 12345

Deploying 60604... ok
```
