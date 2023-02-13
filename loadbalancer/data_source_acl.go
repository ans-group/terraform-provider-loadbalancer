package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceACL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceACLRead,

		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceACLRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, listenerOk := d.GetOk("listener_id")
	targetGroupID, targetGroupOk := d.GetOk("target_group_id")

	if !listenerOk && !targetGroupOk {
		return errors.New("listener_id must be provided when target_group_id is omitted")
	}

	params := connection.APIRequestParameters{}
	if aclID, ok := d.GetOk("acl_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(aclID.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	var acls []loadbalancer.ACL
	var err error
	if listenerOk {
		acls, err = service.GetListenerACLs(listenerID.(int), params)
		if err != nil {
			return fmt.Errorf("Error retrieving listener ACLs: %s", err)
		}
	} else {
		acls, err = service.GetTargetGroupACLs(targetGroupID.(int), params)
		if err != nil {
			return fmt.Errorf("Error retrieving target group ACLs: %s", err)
		}
	}

	if len(acls) < 1 {
		return errors.New("No ACLs found with provided arguments")
	}

	if len(acls) > 1 {
		return errors.New("More than 1 ACL found with provided arguments")
	}

	acl := acls[0]

	d.SetId(strconv.Itoa(acl.ID))
	d.Set("name", acl.Name)
	d.Set("listener_id", acl.ListenerID)
	d.Set("target_group_id", acl.TargetGroupID)

	return nil
}
