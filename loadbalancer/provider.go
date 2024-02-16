package loadbalancer

import (
	"errors"
	"os"

	"github.com/ans-group/sdk-go/pkg/client"
	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: func() (interface{}, error) {
					key := os.Getenv("ANS_API_KEY")
					if key != "" {
						return key, nil
					}

					return "", errors.New("api_key required")
				},
				Description: "API token required to authenticate with ANS APIs. See https://developers.ukfast.io for more details",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"loadbalancer_accessip":    dataSourceAccessIP(),
			"loadbalancer_acl":         dataSourceACL(),
			"loadbalancer_bind":        dataSourceBind(),
			"loadbalancer_certificate": dataSourceCertificate(),
			"loadbalancer_cluster":     dataSourceCluster(),
			"loadbalancer_listener":    dataSourceListener(),
			"loadbalancer_target":      dataSourceTarget(),
			"loadbalancer_targetgroup": dataSourceTargetGroup(),
			"loadbalancer_vip":         dataSourceVip(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"loadbalancer_accessip":    resourceAccessIP(),
			"loadbalancer_acl":         resourceACL(),
			"loadbalancer_bind":        resourceBind(),
			"loadbalancer_certificate": resourceCertificate(),
			"loadbalancer_cluster":     resourceCluster(),
			"loadbalancer_listener":    resourceListener(),
			"loadbalancer_target":      resourceTarget(),
			"loadbalancer_targetgroup": resourceTargetGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return getService(d.Get("api_key").(string)), nil
}

func getClient(apiKey string) client.Client {
	return client.NewClient(connection.NewAPIKeyCredentialsAPIConnection(apiKey))
}

func getService(apiKey string) loadbalancerservice.LoadBalancerService {
	return getClient(apiKey).LoadBalancerService()
}

func setKeys(d *schema.ResourceData, kv map[string]any) diag.Diagnostics {
	for k, v := range kv {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
