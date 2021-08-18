package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"
	loadbalancerservice "github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func resourceTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceTargetCreate,
		Read:   resourceTargetRead,
		Update: resourceTargetUpdate,
		Delete: resourceTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"check_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"check_rise": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"check_fall": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disable_http2": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http2_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceTargetCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	createReq := loadbalancerservice.CreateTargetRequest{
		Name:          d.Get("name").(string),
		IP:            connection.IPAddress(d.Get("ip").(string)),
		Port:          d.Get("port").(int),
		Weight:        d.Get("weight").(int),
		Backup:        d.Get("backup").(bool),
		CheckInterval: d.Get("check_interval").(int),
		CheckSSL:      d.Get("check_ssl").(bool),
		CheckRise:     d.Get("check_rise").(int),
		CheckFall:     d.Get("check_fall").(int),
		DisableHTTP2:  d.Get("disable_http2").(bool),
		HTTP2Only:     d.Get("http2_only").(bool),
		Active:        d.Get("active").(bool),
	}
	log.Printf("[DEBUG] Created CreateTargetRequest: %+v", createReq)

	targetGroupID := d.Get("target_group_id").(int)

	log.Print("[INFO] Creating target")
	target, err := service.CreateTargetGroupTarget(targetGroupID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating target: %s", err)
	}

	d.SetId(strconv.Itoa(target))

	return resourceTargetRead(d, meta)
}

func resourceTargetRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	targetID, _ := strconv.Atoi(d.Id())
	targetGroupID := d.Get("target_group_id").(int)

	log.Printf("[DEBUG] Retrieving Target with ID [%d]", targetID)
	target, err := service.GetTargetGroupTarget(targetGroupID, targetID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.TargetNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", target.Name)
	d.Set("target_group_id", target.TargetGroupID)
	d.Set("ip", target.IP)
	d.Set("port", target.Port)
	d.Set("weight", target.Weight)
	d.Set("backup", target.Backup)
	d.Set("check_interval", target.CheckInterval)
	d.Set("check_ssl", target.CheckSSL)
	d.Set("check_rise", target.CheckRise)
	d.Set("check_fall", target.CheckFall)
	d.Set("disable_http2", target.DisableHTTP2)
	d.Set("http2_only", target.HTTP2Only)
	d.Set("active", target.Active)

	return nil
}

func resourceTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchTargetRequest{}

	targetID, _ := strconv.Atoi(d.Id())
	targetGroupID := d.Get("target_group_id").(int)

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("ip") {
		patchReq.IP = connection.IPAddress(d.Get("ip").(string))
	}

	if d.HasChange("port") {
		patchReq.Port = d.Get("port").(int)
	}

	if d.HasChange("weight") {
		patchReq.Weight = d.Get("weight").(int)
	}

	if d.HasChange("backup") {
		patchReq.Backup = ptr.Bool(d.Get("backup").(bool))
	}

	if d.HasChange("check_interval") {
		patchReq.CheckInterval = d.Get("check_interval").(int)
	}

	if d.HasChange("check_ssl") {
		patchReq.CheckSSL = ptr.Bool(d.Get("check_ssl").(bool))
	}

	if d.HasChange("check_rise") {
		patchReq.CheckRise = d.Get("check_rise").(int)
	}

	if d.HasChange("check_fall") {
		patchReq.CheckFall = d.Get("check_fall").(int)
	}

	if d.HasChange("disable_http2") {
		patchReq.DisableHTTP2 = ptr.Bool(d.Get("disable_http2").(bool))
	}

	if d.HasChange("http2_only") {
		patchReq.HTTP2Only = ptr.Bool(d.Get("http2_only").(bool))
	}

	if d.HasChange("active") {
		patchReq.Active = ptr.Bool(d.Get("active").(bool))
	}

	log.Printf("[INFO] Updating target with ID [%d]", targetID)
	err := service.PatchTargetGroupTarget(targetGroupID, targetID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating target with ID [%d]: %w", targetID, err)
	}

	return resourceTargetRead(d, meta)
}

func resourceTargetDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	targetID, _ := strconv.Atoi(d.Id())
	targetGroupID := d.Get("target_group_id").(int)

	log.Printf("[INFO] Removing target with ID [%d]", targetID)
	err := service.DeleteTargetGroupTarget(targetGroupID, targetID)
	if err != nil {
		return fmt.Errorf("Error removing target with ID [%d]: %s", targetID, err)
	}

	return nil
}
