package loadbalancer

import (
	"context"
	"errors"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/ptr"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerCreate,
		ReadContext:   resourceListenerRead,
		UpdateContext: resourceListenerUpdate,
		DeleteContext: resourceListenerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_target_group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hsts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"hsts_maxage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redirect_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"access_is_allow_list": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_tlsv1": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_tlsv11": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disable_tlsv12": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"custom_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	mode, err := loadbalancerservice.ModeEnum.Parse(d.Get("mode").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "creating listener", map[string]any{
		"name":       d.Get("name"),
		"cluster_id": d.Get("cluster_id"),
	})

	createReq := loadbalancerservice.CreateListenerRequest{
		Name:                 d.Get("name").(string),
		ClusterID:            d.Get("cluster_id").(int),
		Mode:                 mode,
		DefaultTargetGroupID: d.Get("default_target_group_id").(int),
		HSTSMaxAge:           d.Get("hsts_maxage").(int),
		Close:                d.Get("close").(bool),
		RedirectHTTPS:        d.Get("redirect_https").(bool),
		AccessIsAllowList:    d.Get("access_is_allow_list").(bool),
		AllowTLSV1:           d.Get("allow_tlsv1").(bool),
		AllowTLSV11:          d.Get("allow_tlsv11").(bool),
		DisableTLSV12:        d.Get("disable_tlsv12").(bool),
		DisableHTTP2:         d.Get("disable_http2").(bool),
		HTTP2Only:            d.Get("http2_only").(bool),
		CustomCiphers:        d.Get("custom_ciphers").(string),
	}
	tflog.Debug(ctx, "created CreateListenerRequest", map[string]any{
		"request": createReq,
	})

	listener, err := service.CreateListener(createReq)
	if err != nil {
		return diag.Errorf("Error creating listener: %s", err)
	}

	d.SetId(strconv.Itoa(listener))

	return resourceListenerRead(ctx, d, meta)
}

func resourceListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, _ := strconv.Atoi(d.Id())

	tflog.Debug(ctx, "retrieving listener", map[string]any{
		"listener_id": listenerID,
	})

	listener, err := service.GetListener(listenerID)
	if err != nil {
		var listenerNotFoundError *loadbalancerservice.ListenerNotFoundError
		switch {
		case errors.As(err, &listenerNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"name":                    listener.Name,
		"cluster_id":              listener.ClusterID,
		"mode":                    listener.Mode,
		"default_target_group_id": listener.DefaultTargetGroupID,
		"hsts_enabled":            listener.HSTSEnabled,
		"hsts_maxage":             listener.HSTSMaxAge,
		"close":                   listener.Close,
		"redirect_https":          listener.RedirectHTTPS,
		"access_is_allow_list":    listener.AccessIsAllowList,
		"allow_tlsv1":             listener.AllowTLSV1,
		"allow_tlsv11":            listener.AllowTLSV11,
		"disable_tlsv12":          listener.DisableTLSV12,
		"disable_http2":           listener.DisableHTTP2,
		"http2_only":              listener.HTTP2Only,
		"custom_ciphers":          listener.CustomCiphers,
	})
}

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchListenerRequest{}

	listenerID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("mode") {
		mode, err := loadbalancerservice.ModeEnum.Parse(d.Get("mode").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Mode = mode
	}

	if d.HasChange("default_target_group_id") {
		patchReq.DefaultTargetGroupID = d.Get("default_target_group_id").(int)
	}

	if d.HasChange("hsts_enabled") {
		patchReq.HSTSEnabled = ptr.Bool(d.Get("hsts_enabled").(bool))
	}

	if d.HasChange("hsts_maxage") {
		patchReq.HSTSMaxAge = d.Get("hsts_maxage").(int)
	}

	if d.HasChange("close") {
		patchReq.Close = ptr.Bool(d.Get("close").(bool))
	}

	if d.HasChange("redirect_https") {
		patchReq.RedirectHTTPS = ptr.Bool(d.Get("redirect_https").(bool))
	}

	if d.HasChange("access_is_allow_list") {
		patchReq.AccessIsAllowList = ptr.Bool(d.Get("access_is_allow_list").(bool))
	}

	if d.HasChange("allow_tlsv1") {
		patchReq.AllowTLSV1 = ptr.Bool(d.Get("allow_tlsv1").(bool))
	}

	if d.HasChange("allow_tlsv11") {
		patchReq.AllowTLSV11 = ptr.Bool(d.Get("allow_tlsv11").(bool))
	}

	if d.HasChange("disable_tlsv12") {
		patchReq.DisableTLSV12 = ptr.Bool(d.Get("disable_tlsv12").(bool))
	}

	if d.HasChange("disable_http2") {
		patchReq.DisableHTTP2 = ptr.Bool(d.Get("disable_http2").(bool))
	}

	if d.HasChange("http2_only") {
		patchReq.HTTP2Only = ptr.Bool(d.Get("http2_only").(bool))
	}

	if d.HasChange("custom_ciphers") {
		patchReq.CustomCiphers = d.Get("custom_ciphers").(string)
	}

	tflog.Info(ctx, "updating listener", map[string]any{
		"listener_id": listenerID,
	})

	err := service.PatchListener(listenerID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating listener with ID [%d]: %s", listenerID, err)
	}

	return resourceListenerRead(ctx, d, meta)
}

func resourceListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "removing listener", map[string]any{
		"listener_id": listenerID,
	})

	err := service.DeleteListener(listenerID)
	if err != nil {
		return diag.Errorf("Error removing listener with ID [%d]: %s", listenerID, err)
	}

	return nil
}
