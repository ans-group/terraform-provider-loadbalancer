package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTargetGroupRead,

		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"balance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sticky": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cookie_opts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeouts_connect": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timeouts_server": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom_options": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_http_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_expect": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_tcp_monitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"send_proxy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"send_proxy_v2": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_verify": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sni": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("target_group_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if clusterID, ok := d.GetOk("cluster_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("cluster_id", connection.EQOperator, []string{strconv.Itoa(clusterID.(int))}))
	}

	targetgroups, err := service.GetTargetGroups(params)
	if err != nil {
		return diag.Errorf("Error retrieving target groups: %s", err)
	}

	if len(targetgroups) < 1 {
		return diag.Errorf("No target groups found with provided arguments")
	}

	if len(targetgroups) > 1 {
		return diag.Errorf("More than 1 target group found with provided arguments")
	}

	d.SetId(strconv.Itoa(targetgroups[0].ID))
	return setKeys(d, map[string]any{
		"name":                   targetgroups[0].Name,
		"cluster_id":             targetgroups[0].ClusterID,
		"balance":                targetgroups[0].Balance,
		"mode":                   targetgroups[0].Mode,
		"close":                  targetgroups[0].Close,
		"sticky":                 targetgroups[0].Sticky,
		"cookie_opts":            targetgroups[0].CookieOpts,
		"source":                 targetgroups[0].Source,
		"timeouts_connect":       targetgroups[0].TimeoutsConnect,
		"timeouts_server":        targetgroups[0].TimeoutsServer,
		"custom_options":         targetgroups[0].CustomOptions,
		"monitor_url":            targetgroups[0].MonitorURL,
		"monitor_method":         targetgroups[0].MonitorMethod,
		"monitor_host":           targetgroups[0].MonitorHost,
		"monitor_http_version":   targetgroups[0].MonitorHTTPVersion,
		"monitor_expect":         targetgroups[0].MonitorExpect,
		"monitor_tcp_monitoring": targetgroups[0].MonitorTCPMonitoring,
		"check_port":             targetgroups[0].CheckPort,
		"send_proxy":             targetgroups[0].SendProxy,
		"send_proxy_v2":          targetgroups[0].SendProxyV2,
		"ssl":                    targetgroups[0].SSL,
		"ssl_verify":             targetgroups[0].SSLVerify,
		"sni":                    targetgroups[0].SNI,
	})
}
