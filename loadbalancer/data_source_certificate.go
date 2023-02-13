package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCertificateRead,

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"certificate_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	params := connection.APIRequestParameters{}

	if certificateID, ok := d.GetOk("certificate_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{strconv.Itoa(certificateID.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	listenerID := d.Get("listener_id")

	certificates, err := service.GetListenerCertificates(listenerID.(int), params)
	if err != nil {
		return fmt.Errorf("Error retrieving certificates: %s", err)
	}

	if len(certificates) < 1 {
		return errors.New("No certificates found with provided arguments")
	}

	if len(certificates) > 1 {
		return errors.New("More than 1 certificate found with provided arguments")
	}

	certificate := certificates[0]

	d.SetId(strconv.Itoa(certificate.ID))
	d.Set("listener_id", certificate.ListenerID)
	d.Set("name", certificate.Name)

	return nil
}
