package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	loadbalancerservice "github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ca_bundle": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	createReq := loadbalancerservice.CreateCertificateRequest{
		Name:        d.Get("name").(string),
		Key:         d.Get("key").(string),
		Certificate: d.Get("certificate").(string),
		CABundle:    d.Get("ca_bundle").(string),
	}
	log.Printf("[DEBUG] Created CreateCertificateRequest: %+v", createReq)

	listenerID := d.Get("listener_id").(int)

	log.Print("[INFO] Creating certificate")
	certificate, err := service.CreateListenerCertificate(listenerID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating certificate: %s", err)
	}

	d.SetId(strconv.Itoa(certificate))

	return resourceCertificateRead(d, meta)
}

func resourceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	certificateID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	log.Printf("[DEBUG] Retrieving Certificate with ID [%d]", certificateID)
	certificate, err := service.GetListenerCertificate(listenerID, certificateID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.CertificateNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("listener_id", certificate.ListenerID)
	d.Set("name", certificate.Name)

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchCertificateRequest{}

	certificateID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("key") {
		patchReq.Key = d.Get("key").(string)
	}

	if d.HasChange("certificate") {
		patchReq.Certificate = d.Get("certificate").(string)
	}

	if d.HasChange("ca_bundle") {
		patchReq.CABundle = d.Get("ca_bundle").(string)
	}

	log.Printf("[INFO] Updating certificate with ID [%d]", certificateID)
	err := service.PatchListenerCertificate(listenerID, certificateID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating certificate with ID [%d]: %w", certificateID, err)
	}

	return resourceCertificateRead(d, meta)
}

func resourceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	certificateID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	log.Printf("[INFO] Removing certificate with ID [%d]", certificateID)
	err := service.DeleteListenerCertificate(listenerID, certificateID)
	if err != nil {
		return fmt.Errorf("Error removing certificate with ID [%d]: %s", certificateID, err)
	}

	return nil
}
