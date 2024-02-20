package loadbalancer

import (
	"context"
	"errors"
	"strconv"
	"strings"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBind() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBindCreate,
		ReadContext:   resourceBindRead,
		UpdateContext: resourceBindUpdate,
		DeleteContext: resourceBindDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				ids := strings.Split(d.Id(), "/")
				listenerId, err := strconv.Atoi(ids[0])
				if err != nil {
					return nil, err
				}

				d.SetId(ids[1])
				err = d.Set("listener_id", listenerId)

				return []*schema.ResourceData{d}, err
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

func resourceBindCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID := d.Get("listener_id").(int)

	tflog.Info(ctx, "creating bind", map[string]any{
		"vip_id":      d.Get("vip_id"),
		"port":        d.Get("port"),
		"listener_id": d.Get("listener_id"),
	})

	createReq := loadbalancerservice.CreateBindRequest{
		VIPID: d.Get("vip_id").(int),
		Port:  d.Get("port").(int),
	}
	tflog.Debug(ctx, "created CreateBindRequest", map[string]any{
		"create_bind_request": createReq,
	})

	bind, err := service.CreateListenerBind(listenerID, createReq)
	if err != nil {
		return diag.Errorf("Error creating bind: %s", err)
	}

	d.SetId(strconv.Itoa(bind))

	return resourceBindRead(ctx, d, meta)
}

func resourceBindRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	bindID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	tflog.Debug(ctx, "retrieving bind", map[string]any{
		"bind_id":     bindID,
		"listener_id": listenerID,
	})

	bind, err := service.GetListenerBind(listenerID, bindID)
	if err != nil {
		var bindNotFoundError *loadbalancerservice.BindNotFoundError
		switch {
		case errors.As(err, &bindNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"listener_id": bind.ListenerID,
		"vip_id":      bind.VIPID,
		"port":        bind.Port,
	})
}

func resourceBindUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "updating bind", map[string]any{
		"bind_id":     bindID,
		"listener_id": listenerID,
	})

	err := service.PatchListenerBind(listenerID, bindID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating bind with ID [%d]: %s", bindID, err)
	}

	return resourceBindRead(ctx, d, meta)
}

func resourceBindDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	bindID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	tflog.Info(ctx, "removing bind", map[string]any{
		"bind_id":     bindID,
		"listener_id": listenerID,
	})

	err := service.DeleteListenerBind(listenerID, bindID)
	if err != nil {
		return diag.Errorf("Error removing bind with ID [%d]: %s", bindID, err)
	}

	return nil
}
