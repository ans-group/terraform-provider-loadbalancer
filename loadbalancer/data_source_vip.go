package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVipRead,

		Schema: map[string]*schema.Schema{
			"vip_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"internal_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}

	vips, err := service.GetVIPs(params)
	if err != nil {
		return diag.Errorf("Error retrieving vips: %s", err)
	}

	if len(vips) < 1 {
		return diag.Errorf("No vips found with provided arguments")
	}

	if len(vips) > 1 {
		return diag.Errorf("More than 1 vip found with provided arguments")
	}

	d.SetId(strconv.Itoa(vips[0].ID))
	return setKeys(d, map[string]any{
		"cluster_id":    vips[0].ClusterID,
		"internal_cidr": vips[0].InternalCIDR,
		"external_cidr": vips[0].ExternalCIDR,
		"mac_address":   vips[0].MACAddress,
	})
}
