package loadbalancer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceBindCreate,
		Read:   resourceBindRead,
		Update: resourceBindUpdate,
		Delete: resourceBindDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				ids := strings.Split(d.Id(), "/")
				listenerId, err := strconv.Atoi(ids[0])
				if err != nil {
					return nil, err
				}

				d.SetId(ids[1])
				d.Set("listener_id", listenerId)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vip_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceBindCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	createReq := loadbalancerservice.CreateBindRequest{
		VIPID: d.Get("vip_id").(int),
		Port:  d.Get("port").(int),
	}
	log.Printf("[DEBUG] Created CreateBindRequest: %+v", createReq)

	listenerID := d.Get("listener_id").(int)

	log.Print("[INFO] Creating bind")
	bind, err := service.CreateListenerBind(listenerID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating bind: %s", err)
	}

	d.SetId(strconv.Itoa(bind))

	return resourceBindRead(d, meta)
}

func resourceBindRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	bindID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	log.Printf("[DEBUG] Retrieving Bind with ID [%d]", bindID)
	bind, err := service.GetListenerBind(listenerID, bindID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.BindNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("listener_id", bind.ListenerID)
	d.Set("vip_id", bind.VIPID)
	d.Set("port", bind.Port)

	return nil
}

func resourceBindUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchBindRequest{}

	bindID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	if d.HasChange("vip_id") {
		patchReq.VIPID = d.Get("vip_id").(int)
	}

	if d.HasChange("port") {
		patchReq.Port = d.Get("port").(int)
	}

	log.Printf("[INFO] Updating bind with ID [%d]", bindID)
	err := service.PatchListenerBind(listenerID, bindID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating bind with ID [%d]: %w", bindID, err)
	}

	return resourceBindRead(d, meta)
}

func resourceBindDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	bindID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	log.Printf("[INFO] Removing bind with ID [%d]", bindID)
	err := service.DeleteListenerBind(listenerID, bindID)
	if err != nil {
		return fmt.Errorf("Error removing bind with ID [%d]: %s", bindID, err)
	}

	return nil
}
