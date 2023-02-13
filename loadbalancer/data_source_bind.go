package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBind() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBindRead,

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bind_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vip_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceBindRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if bindID, ok := d.GetOk("bind_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(bindID.(int))}))
	}
	if vipID, ok := d.GetOk("vip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vip_id", connection.EQOperator, []string{strconv.Itoa(vipID.(int))}))
	}
	if port, ok := d.GetOk("port"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("port", connection.EQOperator, []string{strconv.Itoa(port.(int))}))
	}

	listenerID := d.Get("listener_id")

	binds, err := service.GetListenerBinds(listenerID.(int), params)
	if err != nil {
		return fmt.Errorf("Error retrieving binds: %s", err)
	}

	if len(binds) < 1 {
		return errors.New("No binds found with provided arguments")
	}

	if len(binds) > 1 {
		return errors.New("More than 1 bind found with provided arguments")
	}

	bind := binds[0]

	d.SetId(strconv.Itoa(bind.ID))
	d.Set("listener_id", bind.ListenerID)
	d.Set("vip_id", bind.VIPID)
	d.Set("port", bind.Port)

	return nil
}
