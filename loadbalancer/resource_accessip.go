package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccessIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessIPCreate,
		Read:   resourceAccessIPRead,
		Update: resourceAccessIPUpdate,
		Delete: resourceAccessIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccessIPCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	createReq := loadbalancerservice.CreateAccessIPRequest{
		IP: connection.IPAddress(d.Get("ip").(string)),
	}
	log.Printf("[DEBUG] Created CreateAccessIPRequest: %+v", createReq)

	listenerID := d.Get("listener_id").(int)

	log.Print("[INFO] Creating access IP")
	accessIP, err := service.CreateListenerAccessIP(listenerID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating access IP: %s", err)
	}

	d.SetId(strconv.Itoa(accessIP))

	return resourceAccessIPRead(d, meta)
}

func resourceAccessIPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	accessIPID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving AccessIP with ID [%d]", accessIPID)
	accessIP, err := service.GetAccessIP(accessIPID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.AccessIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("ip", accessIP.IP)

	return nil
}

func resourceAccessIPUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchAccessIPRequest{}

	accessIPID, _ := strconv.Atoi(d.Id())

	if d.HasChange("ip") {
		patchReq.IP = connection.IPAddress(d.Get("ip").(string))
	}

	log.Printf("[INFO] Updating access IP with ID [%d]", accessIPID)
	err := service.PatchAccessIP(accessIPID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating access IP with ID [%d]: %w", accessIPID, err)
	}

	return resourceAccessIPRead(d, meta)
}

func resourceAccessIPDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	accessIPID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Removing access IP with ID [%d]", accessIPID)
	err := service.DeleteAccessIP(accessIPID)
	if err != nil {
		return fmt.Errorf("Error removing access IP with ID [%d]: %s", accessIPID, err)
	}

	return nil
}
