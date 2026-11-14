package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryPortScan struct {
	JobId                     int64         `json:"job_id"`
	Name                      string        `json:"name"`
	Server                    string        `json:"servers"`
	StartIp                   string        `json:"start_ip"`
	EndIp                     string        `json:"end_ip"`
	CaptureHostsFile          string        `json:"capture_hosts_file"`
	DiscoverCluster           string        `json:"discover_cluster"`
	OverwriteDeviceCategories string        `json:"overwrite_device_categories"`
	ScheduleTime              string        `json:"schedule_time"`
	ScheduleDays              string        `json:"schedule_days"`
	Tags                      []interface{} `json:"tags"`
}

type apiAutoDiscoveryPortScanListResponse struct {
	TotalCount int                        `json:"total_count"`
	Jobs       []apiAutoDiscoveryPortScan `json:"jobs"`
}

func resourceD42AutoDiscoveryPortScan() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_portscan can be use to manage a tcp port scan autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryPortScanCreate,
		Read:        resourceDevice42AutoDiscoveryPortScanRead,
		Update:      resourceDevice42AutoDiscoveryPortScanUpdate,
		Delete:      resourceDevice42AutoDiscoveryPortScanDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"capture_hosts_file": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"discover_cluster": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"overwrite_device_categories": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCommaSeparatedOrderDiff,
			},
		},
	}
}

func resourceDevice42AutoDiscoveryPortScanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryPortScanFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/portscan/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryPortScanRead(d, m)
}

func resourceDevice42AutoDiscoveryPortScanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryPortScanRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryPortScanListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/portscan/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("server", j.Server)
			d.Set("start_ip", j.StartIp)
			d.Set("end_ip", j.EndIp)
			d.Set("capture_hosts_file", j.CaptureHostsFile)
			d.Set("discover_cluster", j.DiscoverCluster)
			d.Set("overwrite_device_categories", j.OverwriteDeviceCategories)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryPortScanRead - no port scan job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryPortScanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryPortScanFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/portscan/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryPortScanRead(d, m)
}

func resourceDevice42AutoDiscoveryPortScanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting port scan autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/portscan/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPortScanDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryPortScanFormData builds the create/update payload. The API
// upserts by name: re-posting the same name updates the existing job, so the
// form data is identical for create and update.
func autoDiscoveryPortScanFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name":   d.Get("name").(string),
		"server": d.Get("server").(string),
	}

	if v, ok := d.GetOk("start_ip"); ok {
		formData["start_ip"] = v.(string)
	}
	if v, ok := d.GetOk("end_ip"); ok {
		formData["end_ip"] = v.(string)
	}
	if v, ok := d.GetOk("capture_hosts_file"); ok {
		formData["capture_hosts_file"] = v.(string)
	}
	if v, ok := d.GetOk("discover_cluster"); ok {
		formData["discover_cluster"] = v.(string)
	}
	if v, ok := d.GetOk("overwrite_device_categories"); ok {
		formData["overwrite_device_categories"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
