package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
)

func dataSourceAccessIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessIPRead,

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

func dataSourceAccessIPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	listenerID, listenerOk := d.GetOk("listener_id")
	accessIPID, accessOk := d.GetOk("access_ip_id")

	if !listenerOk && !accessOk {
		return errors.New("listener_id must be provided when access_ip_id is omitted")
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
			return fmt.Errorf("Error retrieving access IPs: %s", err)
		}

		if len(accessIPs) < 1 {
			return errors.New("No access IPs found with provided arguments")
		}

		if len(accessIPs) > 1 {
			return errors.New("More than 1 access IP found with provided arguments")
		}

		accessIP = accessIPs[0]
	} else {
		accessIP, err = service.GetAccessIP(accessIPID.(int))
		if err != nil {
			return err
		}
	}

	d.SetId(strconv.Itoa(accessIP.ID))
	d.Set("ip", accessIP.IP)

	return nil
}
