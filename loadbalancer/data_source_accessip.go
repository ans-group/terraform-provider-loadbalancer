package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccessIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessIPRead,

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"access_ip_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAccessIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	listenerID, listenerOk := d.GetOk("listener_id")
	accessIPID, accessOk := d.GetOk("access_ip_id")

	if !listenerOk && !accessOk {
		return diag.Errorf("listener_id must be provided when access_ip_id is omitted")
	}

	var accessIP loadbalancer.AccessIP
	var err error
	if listenerOk {
		if accessOk {
			params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(accessIPID.(int))}))
		}
		if ip, ok := d.GetOk("ip"); ok {
			params.WithFilter(*connection.NewAPIRequestFiltering("ip", connection.EQOperator, []string{ip.(string)}))
		}

		accessIPs, err := service.GetListenerAccessIPs(listenerID.(int), params)
		if err != nil {
			return diag.Errorf("Error retrieving access IPs: %s", err)
		}

		if len(accessIPs) < 1 {
			return diag.Errorf("No access IPs found with provided arguments")
		}

		if len(accessIPs) > 1 {
			return diag.Errorf("More than 1 access IP found with provided arguments")
		}

		accessIP = accessIPs[0]
	} else {
		accessIP, err = service.GetAccessIP(accessIPID.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(accessIP.ID))
	return setKeys(d, map[string]any{
		"ip": accessIP.IP,
	})
}
