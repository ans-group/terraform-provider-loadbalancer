package loadbalancer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/ptr"
	loadbalancerservice "github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func resourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTargetGroupCreate,
		Read:   resourceTargetGroupRead,
		Update: resourceTargetGroupUpdate,
		Delete: resourceTargetGroupDelete,
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
			"balance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"close": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sticky": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cookie_opts": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeouts_connect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeouts_server": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"custom_options": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_http_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_expect": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_tcp_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"check_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"send_proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"send_proxy_v2": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sni": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceTargetGroupCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	balance, err := loadbalancerservice.ParseTargetGroupBalance(d.Get("balance").(string))
	if err != nil {
		return err
	}

	mode, err := loadbalancerservice.ParseMode(d.Get("mode").(string))
	if err != nil {
		return err
	}

	var monitorMethod loadbalancerservice.TargetGroupMonitorMethod
	if d.HasChange("monitor_method") {
		monitorMethod, err = loadbalancerservice.ParseTargetGroupMonitorMethod(d.Get("monitor_method").(string))
		if err != nil {
			return err
		}
	}

	createReq := loadbalancerservice.CreateTargetGroupRequest{
		Name:                 d.Get("name").(string),
		ClusterID:            d.Get("cluster_id").(int),
		Balance:              balance,
		Mode:                 mode,
		Close:                d.Get("close").(bool),
		Sticky:               d.Get("sticky").(bool),
		CookieOpts:           d.Get("cookie_opts").(string),
		Source:               d.Get("source").(string),
		TimeoutsConnect:      d.Get("timeouts_connect").(int),
		TimeoutsServer:       d.Get("timeouts_server").(int),
		CustomOptions:        d.Get("custom_options").(string),
		MonitorURL:           d.Get("monitor_url").(string),
		MonitorMethod:        monitorMethod,
		MonitorHost:          d.Get("monitor_host").(string),
		MonitorHTTPVersion:   d.Get("monitor_http_version").(string),
		MonitorExpect:        d.Get("monitor_expect").(string),
		MonitorTCPMonitoring: d.Get("monitor_tcp_monitoring").(bool),
		CheckPort:            d.Get("check_port").(int),
		SendProxy:            d.Get("send_proxy").(bool),
		SendProxyV2:          d.Get("send_proxy_v2").(bool),
		SSL:                  d.Get("ssl").(bool),
		SSLVerify:            d.Get("ssl_verify").(bool),
		SNI:                  d.Get("sni").(bool),
	}
	log.Printf("[DEBUG] Created CreateTargetGroupRequest: %+v", createReq)

	log.Print("[INFO] Creating target group")
	group, err := service.CreateTargetGroup(createReq)
	if err != nil {
		return fmt.Errorf("Error creating target group: %s", err)
	}

	d.SetId(strconv.Itoa(group))

	return resourceTargetGroupRead(d, meta)
}

func resourceTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	groupID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] Retrieving TargetGroup with ID [%d]", groupID)
	group, err := service.GetTargetGroup(groupID)
	if err != nil {
		switch err.(type) {
		case *loadbalancerservice.TargetGroupNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", group.Name)
	d.Set("cluster_id", group.ClusterID)
	d.Set("balance", group.Balance)
	d.Set("mode", group.Mode)
	d.Set("close", group.Close)
	d.Set("sticky", group.Sticky)
	d.Set("cookie_opts", group.CookieOpts)
	d.Set("source", group.Source)
	d.Set("timeouts_connect", group.TimeoutsConnect)
	d.Set("timeouts_server", group.TimeoutsServer)
	d.Set("custom_options", group.CustomOptions)
	d.Set("monitor_url", group.MonitorURL)
	d.Set("monitor_method", group.MonitorMethod)
	d.Set("monitor_host", group.MonitorHost)
	d.Set("monitor_http_version", group.MonitorHTTPVersion)
	d.Set("monitor_expect", group.MonitorExpect)
	d.Set("monitor_tcp_monitoring", group.MonitorTCPMonitoring)
	d.Set("check_port", group.CheckPort)
	d.Set("send_proxy", group.SendProxy)
	d.Set("send_proxy_v2", group.SendProxyV2)
	d.Set("ssl", group.SSL)
	d.Set("ssl_verify", group.SSLVerify)
	d.Set("sni", group.SNI)

	return nil
}

func resourceTargetGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)
	patchReq := loadbalancerservice.PatchTargetGroupRequest{}

	groupID, _ := strconv.Atoi(d.Id())

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("balance") {
		balance, err := loadbalancerservice.ParseTargetGroupBalance(d.Get("balance").(string))
		if err != nil {
			return err
		}

		patchReq.Balance = balance
	}

	if d.HasChange("mode") {
		mode, err := loadbalancerservice.ParseMode(d.Get("mode").(string))
		if err != nil {
			return err
		}

		patchReq.Mode = mode
	}

	if d.HasChange("close") {
		patchReq.Close = ptr.Bool(d.Get("close").(bool))
	}

	if d.HasChange("sticky") {
		patchReq.Sticky = ptr.Bool(d.Get("sticky").(bool))
	}

	if d.HasChange("cookie_opts") {
		patchReq.CookieOpts = d.Get("cookie_opts").(string)
	}

	if d.HasChange("source") {
		patchReq.Source = d.Get("source").(string)
	}

	if d.HasChange("timeouts_connect") {
		patchReq.TimeoutsConnect = d.Get("timeouts_connect").(int)
	}

	if d.HasChange("timeouts_server") {
		patchReq.TimeoutsServer = d.Get("timeouts_server").(int)
	}

	if d.HasChange("custom_options") {
		patchReq.CustomOptions = d.Get("custom_options").(string)
	}

	if d.HasChange("monitor_url") {
		patchReq.MonitorURL = d.Get("monitor_url").(string)
	}

	if d.HasChange("monitor_method") {
		monitorMethod, err := loadbalancerservice.ParseTargetGroupMonitorMethod(d.Get("monitor_method").(string))
		if err != nil {
			return err
		}

		patchReq.MonitorMethod = monitorMethod
	}

	if d.HasChange("monitor_host") {
		patchReq.MonitorHost = d.Get("monitor_host").(string)
	}

	if d.HasChange("monitor_http_version") {
		patchReq.MonitorHTTPVersion = d.Get("monitor_http_version").(string)
	}

	if d.HasChange("monitor_expect") {
		patchReq.MonitorExpect = d.Get("monitor_expect").(string)
	}

	if d.HasChange("monitor_tcp_monitoring") {
		patchReq.MonitorTCPMonitoring = ptr.Bool(d.Get("monitor_tcp_monitoring").(bool))
	}

	if d.HasChange("check_port") {
		patchReq.CheckPort = d.Get("check_port").(int)
	}

	if d.HasChange("send_proxy") {
		patchReq.SendProxy = ptr.Bool(d.Get("send_proxy").(bool))
	}

	if d.HasChange("send_proxy_v2") {
		patchReq.SendProxyV2 = ptr.Bool(d.Get("send_proxy_v2").(bool))
	}

	if d.HasChange("ssl") {
		patchReq.SSL = ptr.Bool(d.Get("ssl").(bool))
	}

	if d.HasChange("ssl_verify") {
		patchReq.SSLVerify = ptr.Bool(d.Get("ssl_verify").(bool))
	}

	if d.HasChange("sni") {
		patchReq.SNI = ptr.Bool(d.Get("sni").(bool))
	}

	log.Printf("[INFO] Updating target group with ID [%d]", groupID)
	err := service.PatchTargetGroup(groupID, patchReq)
	if err != nil {
		return fmt.Errorf("Error updating target group with ID [%d]: %w", groupID, err)
	}

	return resourceTargetGroupRead(d, meta)
}

func resourceTargetGroupDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(loadbalancerservice.LoadBalancerService)

	groupID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Removing target group with ID [%d]", groupID)
	err := service.DeleteTargetGroup(groupID)
	if err != nil {
		return fmt.Errorf("Error removing target group with ID [%d]: %s", groupID, err)
	}

	return nil
}
