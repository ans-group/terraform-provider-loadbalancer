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

func resourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTargetGroupCreate,
		ReadContext:   resourceTargetGroupRead,
		UpdateContext: resourceTargetGroupUpdate,
		DeleteContext: resourceTargetGroupDelete,
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
			"balance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sticky": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cookie_opts": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeouts_connect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeouts_server": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"custom_options": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_http_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_expect": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_tcp_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"check_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"send_proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"send_proxy_v2": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sni": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceTargetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	balance, err := loadbalancerservice.TargetGroupBalanceEnum.Parse(d.Get("balance").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	mode, err := loadbalancerservice.ModeEnum.Parse(d.Get("mode").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var monitorMethod loadbalancerservice.TargetGroupMonitorMethod
	if d.HasChange("monitor_method") {
		monitorMethod, err = loadbalancerservice.TargetGroupMonitorMethodEnum.Parse(d.Get("monitor_method").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	tflog.Info(ctx, "creating target group", map[string]any{
		"name":       d.Get("name"),
		"cluster_id": d.Get("cluster_id"),
		"balance":    d.Get("balance"),
		"mode":       d.Get("mode"),
	})

	createReq := loadbalancerservice.CreateTargetGroupRequest{
		Name:                 d.Get("name").(string),
		ClusterID:            d.Get("cluster_id").(int),
		Balance:              balance,
		Mode:                 mode,
		Close:                d.Get("close").(bool),
		Sticky:               d.Get("sticky").(bool),
		CookieOpts:           d.Get("cookie_opts").(string),
		Source:               d.Get("source").(string),
		TimeoutsConnect:      d.Get("timeouts_connect").(int),
		TimeoutsServer:       d.Get("timeouts_server").(int),
		CustomOptions:        d.Get("custom_options").(string),
		MonitorURL:           d.Get("monitor_url").(string),
		MonitorMethod:        monitorMethod,
		MonitorHost:          d.Get("monitor_host").(string),
		MonitorHTTPVersion:   d.Get("monitor_http_version").(string),
		MonitorExpect:        d.Get("monitor_expect").(string),
		MonitorTCPMonitoring: d.Get("monitor_tcp_monitoring").(bool),
		CheckPort:            d.Get("check_port").(int),
		SendProxy:            d.Get("send_proxy").(bool),
		SendProxyV2:          d.Get("send_proxy_v2").(bool),
		SSL:                  d.Get("ssl").(bool),
		SSLVerify:            d.Get("ssl_verify").(bool),
		SNI:                  d.Get("sni").(bool),
	}
	tflog.Debug(ctx, "created CreateTargetGroupRequest", map[string]any{
		"request": createReq,
	})

	group, err := service.CreateTargetGroup(createReq)
	if err != nil {
		return diag.Errorf("Error creating target group: %s", err)
	}

	d.SetId(strconv.Itoa(group))

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	groupID, _ := strconv.Atoi(d.Id())

	tflog.Debug(ctx, "retrieving target group", map[string]any{
		"target_group_id": groupID,
	})

	group, err := service.GetTargetGroup(groupID)
	if err != nil {
		var targetGroupNotFoundError *loadbalancerservice.TargetGroupNotFoundError
		switch {
		case errors.As(err, &targetGroupNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"name":                   group.Name,
		"cluster_id":             group.ClusterID,
		"balance":                group.Balance,
		"mode":                   group.Mode,
		"close":                  group.Close,
		"sticky":                 group.Sticky,
		"cookie_opts":            group.CookieOpts,
		"source":                 group.Source,
		"timeouts_connect":       group.TimeoutsConnect,
		"timeouts_server":        group.TimeoutsServer,
		"custom_options":         group.CustomOptions,
		"monitor_url":            group.MonitorURL,
		"monitor_method":         group.MonitorMethod,
		"monitor_host":           group.MonitorHost,
		"monitor_http_version":   group.MonitorHTTPVersion,
		"monitor_expect":         group.MonitorExpect,
		"monitor_tcp_monitoring": group.MonitorTCPMonitoring,
		"check_port":             group.CheckPort,
		"send_proxy":             group.SendProxy,
		"send_proxy_v2":          group.SendProxyV2,
		"ssl":                    group.SSL,
		"ssl_verify":             group.SSLVerify,
		"sni":                    group.SNI,
	})
}

func resourceTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchTargetGroupRequest{}

	groupID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("balance") {
		balance, err := loadbalancerservice.TargetGroupBalanceEnum.Parse(d.Get("balance").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Balance = balance
	}

	if d.HasChange("mode") {
		mode, err := loadbalancerservice.ModeEnum.Parse(d.Get("mode").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Mode = mode
	}

	if d.HasChange("close") {
		patchReq.Close = ptr.Bool(d.Get("close").(bool))
	}

	if d.HasChange("sticky") {
		patchReq.Sticky = ptr.Bool(d.Get("sticky").(bool))
	}

	if d.HasChange("cookie_opts") {
		patchReq.CookieOpts = d.Get("cookie_opts").(string)
	}

	if d.HasChange("source") {
		patchReq.Source = d.Get("source").(string)
	}

	if d.HasChange("timeouts_connect") {
		patchReq.TimeoutsConnect = d.Get("timeouts_connect").(int)
	}

	if d.HasChange("timeouts_server") {
		patchReq.TimeoutsServer = d.Get("timeouts_server").(int)
	}

	if d.HasChange("custom_options") {
		patchReq.CustomOptions = d.Get("custom_options").(string)
	}

	if d.HasChange("monitor_url") {
		patchReq.MonitorURL = d.Get("monitor_url").(string)
	}

	if d.HasChange("monitor_method") {
		monitorMethod, err := loadbalancerservice.TargetGroupMonitorMethodEnum.Parse(d.Get("monitor_method").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.MonitorMethod = monitorMethod
	}

	if d.HasChange("monitor_host") {
		patchReq.MonitorHost = d.Get("monitor_host").(string)
	}

	if d.HasChange("monitor_http_version") {
		patchReq.MonitorHTTPVersion = d.Get("monitor_http_version").(string)
	}

	if d.HasChange("monitor_expect") {
		patchReq.MonitorExpect = d.Get("monitor_expect").(string)
	}

	if d.HasChange("monitor_tcp_monitoring") {
		patchReq.MonitorTCPMonitoring = ptr.Bool(d.Get("monitor_tcp_monitoring").(bool))
	}

	if d.HasChange("check_port") {
		patchReq.CheckPort = d.Get("check_port").(int)
	}

	if d.HasChange("send_proxy") {
		patchReq.SendProxy = ptr.Bool(d.Get("send_proxy").(bool))
	}

	if d.HasChange("send_proxy_v2") {
		patchReq.SendProxyV2 = ptr.Bool(d.Get("send_proxy_v2").(bool))
	}

	if d.HasChange("ssl") {
		patchReq.SSL = ptr.Bool(d.Get("ssl").(bool))
	}

	if d.HasChange("ssl_verify") {
		patchReq.SSLVerify = ptr.Bool(d.Get("ssl_verify").(bool))
	}

	if d.HasChange("sni") {
		patchReq.SNI = ptr.Bool(d.Get("sni").(bool))
	}

	tflog.Info(ctx, "updating target group", map[string]any{
		"target_group_id": groupID,
	})

	err := service.PatchTargetGroup(groupID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating target group with ID [%d]: %s", groupID, err)
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	groupID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "removing target group", map[string]any{
		"target_group_id": groupID,
	})

	err := service.DeleteTargetGroup(groupID)
	if err != nil {
		return diag.Errorf("Error removing target group with ID [%d]: %s", groupID, err)
	}

	return nil
}
