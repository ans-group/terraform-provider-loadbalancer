package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("cluster_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if deployed, ok := d.GetOk("deployed"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("deployed", connection.EQOperator, []string{strconv.FormatBool(deployed.(bool))}))
	}

	clusters, err := service.GetClusters(params)
	if err != nil {
		return diag.Errorf("Error retrieving clusters: %s", err)
	}

	if len(clusters) < 1 {
		return diag.Errorf("No clusters found with provided arguments")
	}

	if len(clusters) > 1 {
		return diag.Errorf("More than 1 cluster found with provided arguments")
	}

	d.SetId(strconv.Itoa(clusters[0].ID))
	return setKeys(d, map[string]any{
		"name":     clusters[0].Name,
		"deployed": clusters[0].Deployed,
	})
}
