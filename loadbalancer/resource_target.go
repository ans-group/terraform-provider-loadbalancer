package loadbalancer

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTargetCreate,
		ReadContext:   resourceTargetRead,
		UpdateContext: resourceTargetUpdate,
		DeleteContext: resourceTargetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				ids := strings.Split(d.Id(), "/")
				targetGroupId, err := strconv.Atoi(ids[0])
				if err != nil {
					return nil, err
				}

				d.SetId(ids[1])
				err = d.Set("target_group_id", targetGroupId)

				return []*schema.ResourceData{d}, err
			},
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
				Default:  true,
			},
		},
	}
}

func resourceTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	targetGroupID := d.Get("target_group_id").(int)

	tflog.Info(ctx, "creating target", map[string]any{
		"target_group_id": targetGroupID,
		"name":            d.Get("name"),
		"ip":              d.Get("ip"),
		"port":            d.Get("port"),
	})

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

	tflog.Debug(ctx, "created CreateTargetRequest", map[string]any{
		"request": createReq,
	})

	target, err := service.CreateTargetGroupTarget(targetGroupID, createReq)
	if err != nil {
		return diag.Errorf("Error creating target: %s", err)
	}

	d.SetId(strconv.Itoa(target))

	return resourceTargetRead(ctx, d, meta)
}

func resourceTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	targetID, _ := strconv.Atoi(d.Id())
	targetGroupID := d.Get("target_group_id").(int)

	tflog.Debug(ctx, "retrieving target", map[string]any{
		"target_id":       targetID,
		"target_group_id": targetGroupID,
	})

	target, err := service.GetTargetGroupTarget(targetGroupID, targetID)
	if err != nil {
		var targetNotFoundError *loadbalancerservice.TargetNotFoundError
		switch {
		case errors.As(err, &targetNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"name":            target.Name,
		"target_group_id": target.TargetGroupID,
		"ip":              target.IP,
		"port":            target.Port,
		"weight":          target.Weight,
		"backup":          target.Backup,
		"check_interval":  target.CheckInterval,
		"check_ssl":       target.CheckSSL,
		"check_rise":      target.CheckRise,
		"check_fall":      target.CheckFall,
		"disable_http2":   target.DisableHTTP2,
		"http2_only":      target.HTTP2Only,
		"active":          target.Active,
	})
}

func resourceTargetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "updating target", map[string]any{
		"target_id":       targetID,
		"target_group_id": targetGroupID,
	})

	err := service.PatchTargetGroupTarget(targetGroupID, targetID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating target with ID [%d]: %s", targetID, err)
	}

	return resourceTargetRead(ctx, d, meta)
}

func resourceTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	targetID, _ := strconv.Atoi(d.Id())
	targetGroupID := d.Get("target_group_id").(int)

	tflog.Info(ctx, "removing target", map[string]any{
		"target_id":       targetID,
		"target_group_id": targetGroupID,
	})

	err := service.DeleteTargetGroupTarget(targetGroupID, targetID)
	if err != nil {
		return diag.Errorf("Error removing target with ID [%d]: %s", targetID, err)
	}

	return nil
}
