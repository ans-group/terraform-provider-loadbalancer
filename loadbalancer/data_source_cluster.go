package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("cluster_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(id.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if deployed, ok := d.GetOk("deployed"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("deployed", connection.EQOperator, []string{strconv.FormatBool(deployed.(bool))}))
	}

	clusters, err := service.GetClusters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving clusters: %s", err)
	}

	if len(clusters) < 1 {
		return errors.New("No clusters found with provided arguments")
	}

	if len(clusters) > 1 {
		return errors.New("More than 1 cluster found with provided arguments")
	}

	d.SetId(strconv.Itoa(clusters[0].ID))
	d.Set("name", clusters[0].Name)
	d.Set("deployed", strconv.FormatBool(clusters[0].Deployed))

	return nil
}
