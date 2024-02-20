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

func resourceACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceACLCreate,
		ReadContext:   resourceACLRead,
		UpdateContext: resourceACLUpdate,
		DeleteContext: resourceACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"listener_id", "target_group_id"},
			},
			"target_group_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"listener_id", "target_group_id"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"argument": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"action": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"argument": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	tflog.Info(ctx, "creating ACL", map[string]any{
		"listener_id":     d.Get("listener_id"),
		"target_group_id": d.Get("target_group_id"),
		"name":            d.Get("name"),
	})

	createReq := loadbalancerservice.CreateACLRequest{
		ListenerID:    d.Get("listener_id").(int),
		TargetGroupID: d.Get("target_group_id").(int),
		Name:          d.Get("name").(string),
		Conditions:    expandACLConditions(d.Get("condition").([]interface{})),
		Actions:       expandACLActions(d.Get("action").([]interface{})),
	}

	tflog.Debug(ctx, "created CreateACLRequest", map[string]any{
		"create_acl_request": createReq,
	})

	acl, err := service.CreateACL(createReq)
	if err != nil {
		return diag.Errorf("Error creating ACL: %s", err)
	}

	d.SetId(strconv.Itoa(acl))

	return resourceACLRead(ctx, d, meta)
}

func resourceACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	aclID, _ := strconv.Atoi(d.Id())

	tflog.Debug(ctx, "retrieving ACL", map[string]any{
		"acl_id": aclID,
	})

	acl, err := service.GetACL(aclID)
	if err != nil {
		var ACLNotFoundError *loadbalancerservice.ACLNotFoundError
		switch {
		case errors.As(err, &ACLNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"listener_id":     acl.ListenerID,
		"target_group_id": acl.TargetGroupID,
		"name":            acl.Name,
		"condition":       flattenACLConditions(acl.Conditions),
		"action":          flattenACLActions(acl.Actions),
	})
}

func resourceACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchACLRequest{}

	aclID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("condition") {
		patchReq.Conditions = expandACLConditions(d.Get("condition").([]interface{}))
	}

	if d.HasChange("action") {
		patchReq.Actions = expandACLActions(d.Get("action").([]interface{}))
	}

	tflog.Info(ctx, "updating ACL", map[string]any{
		"acl_id": aclID,
	})

	err := service.PatchACL(aclID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating ACL with ID [%d]: %s", aclID, err)
	}

	return resourceACLRead(ctx, d, meta)
}

func resourceACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	aclID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "removing ACL", map[string]any{
		"acl_id": aclID,
	})

	err := service.DeleteACL(aclID)
	if err != nil {
		return diag.Errorf("Error removing ACL with ID [%d]: %s", aclID, err)
	}

	return nil
}
