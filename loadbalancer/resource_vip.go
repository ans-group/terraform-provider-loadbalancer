package loadbalancer

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	loadbalancerservice "github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func resourceVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessIPCreate,
		Read:   resourceAccessIPRead,
		Update: resourceAccessIPUpdate,
		Delete: resourceAccessIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"internal_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVipCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceClusterRead(d, meta)
}

func resourceVipRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	vipID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving VIP with ID [%d]", vipID)
	vip, err := service.GetVIP(vipID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.VIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("cluster_id", vip.ClusterID)
	d.Set("internal_cidr", vip.InternalCIDR)
	d.Set("external_cidr", vip.ExternalCIDR)
	d.Set("mac_address", vip.MACAddress)

	return nil
}

func resourceVipUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVipDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
