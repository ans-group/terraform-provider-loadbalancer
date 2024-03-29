package loadbalancer

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTargetRead,

		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_rise": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_fall": {
				Type:     schema.TypeInt,
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
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("target_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if ip, ok := d.GetOk("ip"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("ip", connection.EQOperator, []string{ip.(string)}))
	}
	if port, ok := d.GetOk("port"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("port", connection.EQOperator, []string{strconv.Itoa(port.(int))}))
	}

	targetgroupID := d.Get("target_group_id").(int)

	targets, err := service.GetTargetGroupTargets(targetgroupID, params)
	if err != nil {
		return diag.Errorf("Error retrieving targets: %s", err)
	}

	if len(targets) < 1 {
		return diag.Errorf("No targets found with provided arguments")
	}

	if len(targets) > 1 {
		return diag.Errorf("More than 1 target found with provided arguments")
	}

	d.SetId(strconv.Itoa(targets[0].ID))
	return setKeys(d, map[string]any{
		"target_group_id": targets[0].TargetGroupID,
		"name":            targets[0].Name,
		"ip":              targets[0].IP,
		"port":            targets[0].Port,
		"weight":          targets[0].Weight,
		"backup":          targets[0].Backup,
		"check_interval":  targets[0].CheckInterval,
		"check_ssl":       targets[0].CheckSSL,
		"check_rise":      targets[0].CheckRise,
		"check_fall":      targets[0].CheckFall,
		"disable_http2":   targets[0].DisableHTTP2,
		"http2_only":      targets[0].HTTP2Only,
		"active":          targets[0].Active,
	})
}
