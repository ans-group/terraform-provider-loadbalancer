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

func dataSourceACL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceACLRead,

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

func dataSourceACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, listenerOk := d.GetOk("listener_id")
	targetGroupID, targetGroupOk := d.GetOk("target_group_id")

	if !listenerOk && !targetGroupOk {
		return diag.Errorf("listener_id must be provided when target_group_id is omitted")
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
			return diag.Errorf("Error retrieving listener ACLs: %s", err)
		}
	} else {
		acls, err = service.GetTargetGroupACLs(targetGroupID.(int), params)
		if err != nil {
			return diag.Errorf("Error retrieving target group ACLs: %s", err)
		}
	}

	if len(acls) < 1 {
		return diag.Errorf("No ACLs found with provided arguments")
	}

	if len(acls) > 1 {
		return diag.Errorf("More than 1 ACL found with provided arguments")
	}

	acl := acls[0]

	d.SetId(strconv.Itoa(acl.ID))
	return setKeys(d, map[string]any{
		"name":            acl.Name,
		"listener_id":     acl.ListenerID,
		"target_group_id": acl.TargetGroupID,
	})
}
