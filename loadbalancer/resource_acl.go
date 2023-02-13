package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceACLCreate,
		Read:   resourceACLRead,
		Update: resourceACLUpdate,
		Delete: resourceACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceACLCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	createReq := loadbalancerservice.CreateACLRequest{
		ListenerID:    d.Get("listener_id").(int),
		TargetGroupID: d.Get("target_group_id").(int),
		Name:          d.Get("name").(string),
		Conditions:    expandACLConditions(d.Get("condition").([]interface{})),
		Actions:       expandACLActions(d.Get("action").([]interface{})),
	}
	log.Printf("[DEBUG] Created CreateACLRequest: %+v", createReq)

	//return fmt.Errorf("%+v", createReq)

	log.Print("[INFO] Creating ACL")
	acl, err := service.CreateACL(createReq)
	if err != nil {
		return fmt.Errorf("Error creating ACL: %s", err)
	}

	d.SetId(strconv.Itoa(acl))

	return resourceACLRead(d, meta)
}

func resourceACLRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	aclID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving ACL with ID [%d]", aclID)
	acl, err := service.GetACL(aclID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.ACLNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("listener_id", acl.ListenerID)
	d.Set("target_group_id", acl.TargetGroupID)
	d.Set("name", acl.Name)
	d.Set("condition", flattenACLConditions(acl.Conditions))
	d.Set("action", flattenACLActions(acl.Actions))

	return nil
}

func resourceACLUpdate(d *schema.ResourceData, meta interface{}) error {
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

	log.Printf("[INFO] Updating ACL with ID [%d]", aclID)
	err := service.PatchACL(aclID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating ACL with ID [%d]: %w", aclID, err)
	}

	return resourceACLRead(d, meta)
}

func resourceACLDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	aclID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Removing ACL with ID [%d]", aclID)
	err := service.DeleteACL(aclID)
	if err != nil {
		return fmt.Errorf("Error removing ACL with ID [%d]: %s", aclID, err)
	}

	return nil
}
