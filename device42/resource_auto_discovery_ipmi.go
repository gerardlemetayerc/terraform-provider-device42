package device42

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryIpmi struct {
	JobId                     int64           `json:"job_id"`
	Name                      string          `json:"name"`
	Server                    string          `json:"server"`
	VrfGroup                  json.RawMessage `json:"vrfgroup"`
	ScheduleTime              string          `json:"schedule_time"`
	ObjectCategory            string          `json:"object_category"`
	UpdateModelIfFound        string          `json:"update_model_if_found"`
	OverwriteObjectCategories string          `json:"overwrite_object_categories"`
	IpStart                   string          `json:"ip_start"`
	IpEnd                     string          `json:"ip_end"`
	OverwriteDeviceName       string          `json:"overwrite_device_name"`
	HostnameToUse             int64           `json:"hostname_to_use"`
	DebugLevel                string          `json:"debug_level"`
	DiscoveryType             string          `json:"discovery_type"`
	RunAsOperator             string          `json:"run_as_operator"`
	ScheduleDays              string          `json:"schedule_days"`
	ExcludeServers            string          `json:"exclude_servers"`
	Tags                      []interface{}   `json:"tags"`
}

type apiAutoDiscoveryIpmiListResponse struct {
	Jobs []apiAutoDiscoveryIpmi `json:"jobs"`
}

func resourceD42AutoDiscoveryIpmi() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_ipmi can be use to manage an IPMI/Redfish autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryIpmiCreate,
		Read:        resourceDevice42AutoDiscoveryIpmiRead,
		Update:      resourceDevice42AutoDiscoveryIpmiUpdate,
		Delete:      resourceDevice42AutoDiscoveryIpmiDelete,

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
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_model_if_found": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite_object_categories": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite_device_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname_to_use": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"debug_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"discovery_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"run_as_operator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_servers": {
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

func resourceDevice42AutoDiscoveryIpmiCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryIpmiFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/ipmi/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryIpmiRead(d, m)
}

func resourceDevice42AutoDiscoveryIpmiRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryIpmiRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryIpmiListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/ipmi/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("server", j.Server)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("object_category", j.ObjectCategory)
			d.Set("update_model_if_found", j.UpdateModelIfFound)
			d.Set("overwrite_object_categories", j.OverwriteObjectCategories)
			d.Set("ip_start", j.IpStart)
			d.Set("ip_end", j.IpEnd)
			d.Set("overwrite_device_name", j.OverwriteDeviceName)
			d.Set("hostname_to_use", j.HostnameToUse)
			d.Set("debug_level", j.DebugLevel)
			d.Set("discovery_type", j.DiscoveryType)
			d.Set("run_as_operator", j.RunAsOperator)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("exclude_servers", j.ExcludeServers)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryIpmiRead - no ipmi job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryIpmiUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryIpmiFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/ipmi/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryIpmiRead(d, m)
}

func resourceDevice42AutoDiscoveryIpmiDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting ipmi autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/ipmi/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryIpmiDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryIpmiFormData builds the create/update payload. The API upserts
// by name: re-posting the same name updates the existing job, so the form data
// is identical for create and update.
func autoDiscoveryIpmiFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("server"); ok {
		formData["server"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("object_category"); ok {
		formData["object_category"] = v.(string)
	}
	if v, ok := d.GetOk("update_model_if_found"); ok {
		formData["update_model_if_found"] = v.(string)
	}
	if v, ok := d.GetOk("overwrite_object_categories"); ok {
		formData["overwrite_object_categories"] = v.(string)
	}
	if v, ok := d.GetOk("ip_start"); ok {
		formData["ip_start"] = v.(string)
	}
	if v, ok := d.GetOk("ip_end"); ok {
		formData["ip_end"] = v.(string)
	}
	if v, ok := d.GetOk("overwrite_device_name"); ok {
		formData["overwrite_device_name"] = v.(string)
	}
	if v, ok := d.GetOk("hostname_to_use"); ok {
		formData["hostname_to_use"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("debug_level"); ok {
		formData["debug_level"] = v.(string)
	}
	if v, ok := d.GetOk("discovery_type"); ok {
		formData["discovery_type"] = v.(string)
	}
	if v, ok := d.GetOk("run_as_operator"); ok {
		formData["run_as_operator"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}
	if v, ok := d.GetOk("exclude_servers"); ok {
		formData["exclude_servers"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
