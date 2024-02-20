package loadbalancer

import (
	"context"
	"errors"
	"strconv"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("The loadbalancer_cluster resource can only be imported at this time")
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	clusterID, _ := strconv.Atoi(d.Id())

	tflog.Debug(ctx, "retrieving cluster", map[string]any{
		"cluster_id": clusterID,
	})

	cluster, err := service.GetCluster(clusterID)
	if err != nil {
		var clusterNotFoundError *loadbalancerservice.ClusterNotFoundError
		switch {
		case errors.As(err, &clusterNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"name": cluster.Name,
	})
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchClusterRequest{}

	clusterID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	tflog.Info(ctx, "updating cluster", map[string]any{
		"cluster_id": clusterID,
	})

	err := service.PatchCluster(clusterID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating cluster with ID [%d]: %s", clusterID, err)
	}

	return resourceClusterRead(ctx, d, meta)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
