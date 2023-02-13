package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/ptr"
	loadbalancerservice "github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerCreate,
		Read:   resourceListenerRead,
		Update: resourceListenerUpdate,
		Delete: resourceListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_target_group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hsts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"hsts_maxage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"redirect_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"access_is_allow_list": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_tlsv1": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_tlsv11": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disable_tlsv12": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disable_http2": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http2_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceListenerCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	mode, err := loadbalancerservice.ParseMode(d.Get("mode").(string))
	if err != nil {
		return err
	}

	createReq := loadbalancerservice.CreateListenerRequest{
		Name:                 d.Get("name").(string),
		ClusterID:            d.Get("cluster_id").(int),
		Mode:                 mode,
		DefaultTargetGroupID: d.Get("default_target_group_id").(int),
		HSTSMaxAge:           d.Get("hsts_maxage").(int),
		Close:                d.Get("close").(bool),
		RedirectHTTPS:        d.Get("redirect_https").(bool),
		AccessIsAllowList:    d.Get("access_is_allow_list").(bool),
		AllowTLSV1:           d.Get("allow_tlsv1").(bool),
		AllowTLSV11:          d.Get("allow_tlsv11").(bool),
		DisableTLSV12:        d.Get("disable_tlsv12").(bool),
		DisableHTTP2:         d.Get("disable_http2").(bool),
		HTTP2Only:            d.Get("http2_only").(bool),
		CustomCiphers:        d.Get("custom_ciphers").(string),
	}
	log.Printf("[DEBUG] Created CreateListenerRequest: %+v", createReq)

	log.Print("[INFO] Creating listener")
	listener, err := service.CreateListener(createReq)
	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
	}

	d.SetId(strconv.Itoa(listener))

	return resourceListenerRead(d, meta)
}

func resourceListenerRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving Listener with ID [%d]", listenerID)
	listener, err := service.GetListener(listenerID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.ListenerNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", listener.Name)
	d.Set("cluster_id", listener.ClusterID)
	d.Set("mode", listener.Mode)
	d.Set("default_target_group_id", listener.DefaultTargetGroupID)
	d.Set("hsts_enabled", listener.HSTSEnabled)
	d.Set("hsts_maxage", listener.HSTSMaxAge)
	d.Set("close", listener.Close)
	d.Set("redirect_https", listener.RedirectHTTPS)
	d.Set("access_is_allow_list", listener.AccessIsAllowList)
	d.Set("allow_tlsv1", listener.AllowTLSV1)
	d.Set("allow_tlsv11", listener.AllowTLSV11)
	d.Set("disable_tlsv12", listener.DisableTLSV12)
	d.Set("disable_http2", listener.DisableHTTP2)
	d.Set("http2_only", listener.HTTP2Only)
	d.Set("custom_ciphers", listener.CustomCiphers)

	return nil
}

func resourceListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchListenerRequest{}

	listenerID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("mode") {
		mode, err := loadbalancerservice.ParseMode(d.Get("mode").(string))
		if err != nil {
			return err
		}

		patchReq.Mode = mode
	}

	if d.HasChange("default_target_group_id") {
		patchReq.DefaultTargetGroupID = d.Get("default_target_group_id").(int)
	}

	if d.HasChange("hsts_enabled") {
		patchReq.HSTSEnabled = ptr.Bool(d.Get("hsts_enabled").(bool))
	}

	if d.HasChange("hsts_maxage") {
		patchReq.HSTSMaxAge = d.Get("hsts_maxage").(int)
	}

	if d.HasChange("close") {
		patchReq.Close = ptr.Bool(d.Get("close").(bool))
	}

	if d.HasChange("redirect_https") {
		patchReq.RedirectHTTPS = ptr.Bool(d.Get("redirect_https").(bool))
	}

	if d.HasChange("access_is_allow_list") {
		patchReq.AccessIsAllowList = ptr.Bool(d.Get("access_is_allow_list").(bool))
	}

	if d.HasChange("allow_tlsv1") {
		patchReq.AllowTLSV1 = ptr.Bool(d.Get("allow_tlsv1").(bool))
	}

	if d.HasChange("allow_tlsv11") {
		patchReq.AllowTLSV11 = ptr.Bool(d.Get("allow_tlsv11").(bool))
	}

	if d.HasChange("disable_tlsv12") {
		patchReq.DisableTLSV12 = ptr.Bool(d.Get("disable_tlsv12").(bool))
	}

	if d.HasChange("disable_http2") {
		patchReq.DisableHTTP2 = ptr.Bool(d.Get("disable_http2").(bool))
	}

	if d.HasChange("http2_only") {
		patchReq.HTTP2Only = ptr.Bool(d.Get("http2_only").(bool))
	}

	if d.HasChange("custom_ciphers") {
		patchReq.CustomCiphers = d.Get("custom_ciphers").(string)
	}

	log.Printf("[INFO] Updating listener with ID [%d]", listenerID)
	err := service.PatchListener(listenerID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating listener with ID [%d]: %w", listenerID, err)
	}

	return resourceListenerRead(d, meta)
}

func resourceListenerDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	listenerID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Removing listener with ID [%d]", listenerID)
	err := service.DeleteListener(listenerID)
	if err != nil {
		return fmt.Errorf("Error removing listener with ID [%d]: %s", listenerID, err)
	}

	return nil
}
