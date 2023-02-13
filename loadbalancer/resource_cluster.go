package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(strconv.Itoa(d.Get("cluster_id").(int)))

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	clusterID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving Cluster with ID [%d]", clusterID)
	cluster, err := service.GetCluster(clusterID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.ClusterNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", cluster.Name)

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchClusterRequest{}

	clusterID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	log.Printf("[INFO] Updating cluster with ID [%d]", clusterID)
	err := service.PatchCluster(clusterID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating cluster with ID [%d]: %w", clusterID, err)
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
