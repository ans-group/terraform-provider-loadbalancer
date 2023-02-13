package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTargetGroupRead,

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

func dataSourceTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving target groups: %s", err)
	}

	if len(targetgroups) < 1 {
		return errors.New("No target groups found with provided arguments")
	}

	if len(targetgroups) > 1 {
		return errors.New("More than 1 target group found with provided arguments")
	}

	d.SetId(strconv.Itoa(targetgroups[0].ID))
	d.Set("name", targetgroups[0].Name)
	d.Set("cluster_id", targetgroups[0].ClusterID)
	d.Set("balance", targetgroups[0].Balance)
	d.Set("mode", targetgroups[0].Mode)
	d.Set("close", targetgroups[0].Close)
	d.Set("sticky", targetgroups[0].Sticky)
	d.Set("cookie_opts", targetgroups[0].CookieOpts)
	d.Set("source", targetgroups[0].Source)
	d.Set("timeouts_connect", targetgroups[0].TimeoutsConnect)
	d.Set("timeouts_server", targetgroups[0].TimeoutsServer)
	d.Set("custom_options", targetgroups[0].CustomOptions)
	d.Set("monitor_url", targetgroups[0].MonitorURL)
	d.Set("monitor_method", targetgroups[0].MonitorMethod)
	d.Set("monitor_host", targetgroups[0].MonitorHost)
	d.Set("monitor_http_version", targetgroups[0].MonitorHTTPVersion)
	d.Set("monitor_expect", targetgroups[0].MonitorExpect)
	d.Set("monitor_tcp_monitoring", targetgroups[0].MonitorTCPMonitoring)
	d.Set("check_port", targetgroups[0].CheckPort)
	d.Set("send_proxy", targetgroups[0].SendProxy)
	d.Set("send_proxy_v2", targetgroups[0].SendProxyV2)
	d.Set("ssl", targetgroups[0].SSL)
	d.Set("ssl_verify", targetgroups[0].SSLVerify)
	d.Set("sni", targetgroups[0].SNI)

	return nil
}
