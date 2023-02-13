package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVipRead,

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

func dataSourceVipRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}

	vips, err := service.GetVIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving vips: %s", err)
	}

	if len(vips) < 1 {
		return errors.New("No vips found with provided arguments")
	}

	if len(vips) > 1 {
		return errors.New("More than 1 vip found with provided arguments")
	}

	d.SetId(strconv.Itoa(vips[0].ID))
	d.Set("cluster_id", vips[0].ClusterID)
	d.Set("internal_cidr", vips[0].InternalCIDR)
	d.Set("external_cidr", vips[0].ExternalCIDR)
	d.Set("mac_address", vips[0].MACAddress)

	return nil
}
