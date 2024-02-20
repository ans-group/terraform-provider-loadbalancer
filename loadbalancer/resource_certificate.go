package loadbalancer

import (
	"context"
	"errors"
	"strconv"

	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID := d.Get("listener_id").(int)

	tflog.Info(ctx, "creating certificate", map[string]any{
		"name":        d.Get("name"),
		"listener_id": d.Get("listener_id"),
	})

	createReq := loadbalancerservice.CreateCertificateRequest{
		Name:        d.Get("name").(string),
		Key:         d.Get("key").(string),
		Certificate: d.Get("certificate").(string),
		CABundle:    d.Get("ca_bundle").(string),
	}

	tflog.Debug(ctx, "created CreateCertificateRequest", map[string]any{
		"request": createReq,
	})

	certificate, err := service.CreateListenerCertificate(listenerID, createReq)
	if err != nil {
		return diag.Errorf("Error creating certificate: %s", err)
	}

	d.SetId(strconv.Itoa(certificate))

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	certificateID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	tflog.Debug(ctx, "retrieving certificate", map[string]any{
		"certificate_id": certificateID,
		"listener_id":    listenerID,
	})

	certificate, err := service.GetListenerCertificate(listenerID, certificateID)
	if err != nil {
		var certificateNotFoundError *loadbalancerservice.CertificateNotFoundError
		switch {
		case errors.As(err, &certificateNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	return setKeys(d, map[string]any{
		"listener_id": certificate.ListenerID,
		"name":        certificate.Name,
	})
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "updating certificate", map[string]any{
		"certificate_id": certificateID,
		"listener_id":    listenerID,
		"name":           d.Get("name"),
	})

	err := service.PatchListenerCertificate(listenerID, certificateID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating certificate with ID [%d]: %s", certificateID, err)
	}

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(loadbalancerservice.LoadBalancerService)

	certificateID, _ := strconv.Atoi(d.Id())
	listenerID := d.Get("listener_id").(int)

	tflog.Info(ctx, "removing certificate", map[string]any{
		"certificate_id": certificateID,
		"listener_id":    listenerID,
	})

	err := service.DeleteListenerCertificate(listenerID, certificateID)
	if err != nil {
		return diag.Errorf("Error removing certificate with ID [%d]: %s", certificateID, err)
	}

	return nil
}
