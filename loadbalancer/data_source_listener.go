package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceListener() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListenerRead,

		Schema: map[string]*schema.Schema{
			"listener_id": {
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
			"hsts_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hsts_maxage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"redirect_https": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_target_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allow_tlsv1": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_tlsv11": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_tlsv12": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_http2": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"http2_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_ciphers": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("listener_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if clusterID, ok := d.GetOk("cluster_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("cluster_id", connection.EQOperator, []string{strconv.Itoa(clusterID.(int))}))
	}

	listeners, err := service.GetListeners(params)
	if err != nil {
		return diag.Errorf("Error retrieving listeners: %s", err)
	}

	if len(listeners) < 1 {
		return diag.Errorf("No listeners found with provided arguments")
	}

	if len(listeners) > 1 {
		return diag.Errorf("More than 1 listener found with provided arguments")
	}

	d.SetId(strconv.Itoa(listeners[0].ID))

	return setKeys(d, map[string]any{
		"name":                    listeners[0].Name,
		"cluster_id":              listeners[0].ClusterID,
		"hsts_enabled":            listeners[0].HSTSEnabled,
		"mode":                    listeners[0].Mode,
		"hsts_maxage":             listeners[0].HSTSMaxAge,
		"close":                   listeners[0].Close,
		"redirect_https":          listeners[0].RedirectHTTPS,
		"default_target_group_id": listeners[0].DefaultTargetGroupID,
		"allow_tlsv1":             listeners[0].AllowTLSV1,
		"allow_tlsv11":            listeners[0].AllowTLSV11,
		"disable_tlsv12":          listeners[0].DisableTLSV12,
		"disable_http2":           listeners[0].DisableHTTP2,
		"http2_only":              listeners[0].HTTP2Only,
	})
}
