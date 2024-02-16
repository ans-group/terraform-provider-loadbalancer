package loadbalancer

import (
	"context"
	"errors"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccessIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessIPCreate,
		ReadContext:   resourceAccessIPRead,
		UpdateContext: resourceAccessIPUpdate,
		DeleteContext: resourceAccessIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceAccessIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID := d.Get("listener_id").(int)

	tflog.Info(ctx, "creating access IP", map[string]any{
		"ip":          d.Get("ip"),
		"listener_id": d.Get("listener_id"),
	})

	createReq := loadbalancerservice.CreateAccessIPRequest{
		IP: connection.IPAddress(d.Get("ip").(string)),
	}
	tflog.Debug(ctx, "created CreatedAccessIPRequest", map[string]any{
		"request": createReq,
	})

	accessIP, err := service.CreateListenerAccessIP(listenerID, createReq)
	if err != nil {
		return diag.Errorf("Error creating access IP: %s", err)
	}

	d.SetId(strconv.Itoa(accessIP))

	return resourceAccessIPRead(ctx, d, meta)
}

func resourceAccessIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	accessIPID, _ := strconv.Atoi(d.Id())

	tflog.Debug(ctx, "retrieving AccessIP", map[string]any{
		"access_ip_id": accessIPID,
	})
	accessIP, err := service.GetAccessIP(accessIPID)
	if err != nil {
		var accessIPNotFoundError *loadbalancerservice.AccessIPNotFoundError
		switch {
		case errors.As(err, &accessIPNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"ip": accessIP.IP,
	})
}

func resourceAccessIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchAccessIPRequest{}

	accessIPID, _ := strconv.Atoi(d.Id())

	if d.HasChange("ip") {
		patchReq.IP = connection.IPAddress(d.Get("ip").(string))
	}

	tflog.Info(ctx, "updating access IP", map[string]any{
		"access_ip_id": accessIPID,
	})
	err := service.PatchAccessIP(accessIPID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating access IP with ID [%d]: %s", accessIPID, err)
	}

	return resourceAccessIPRead(ctx, d, meta)
}

func resourceAccessIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	accessIPID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "removing access IP", map[string]any{
		"access_ip_id": accessIPID,
	})

	err := service.DeleteAccessIP(accessIPID)
	if err != nil {
		return diag.Errorf("Error removing access IP with ID [%d]: %s", accessIPID, err)
	}

	return nil
}
