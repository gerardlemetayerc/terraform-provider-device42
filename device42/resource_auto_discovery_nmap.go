package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryNmap struct {
	JobId          int64         `json:"job_id"`
	Name           string        `json:"name"`
	Server         string        `json:"server"`
	StartIp        string        `json:"start_ip"`
	EndIp          string        `json:"end_ip"`
	Target         string        `json:"target"`
	ExcludeServers string        `json:"exclude_servers"`
	ScheduleTime   string        `json:"schedule_time"`
	ScheduleDays   string        `json:"schedule_days"`
	Tags           []interface{} `json:"tags"`
}

type apiAutoDiscoveryNmapListResponse struct {
	Jobs []apiAutoDiscoveryNmap `json:"jobs"`
}

func resourceD42AutoDiscoveryNmap() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_nmap can be use to manage an nmap autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryNmapCreate,
		Read:        resourceDevice42AutoDiscoveryNmapRead,
		Update:      resourceDevice42AutoDiscoveryNmapUpdate,
		Delete:      resourceDevice42AutoDiscoveryNmapDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_servers": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceDevice42AutoDiscoveryNmapCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryNmapFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/nmap/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryNmapRead(d, m)
}

func resourceDevice42AutoDiscoveryNmapRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryNmapRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryNmapListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/nmap/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("server", j.Server)
			d.Set("start_ip", j.StartIp)
			d.Set("end_ip", j.EndIp)
			d.Set("target", j.Target)
			d.Set("exclude_servers", j.ExcludeServers)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryNmapRead - no nmap job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryNmapUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryNmapFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/nmap/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryNmapRead(d, m)
}

func resourceDevice42AutoDiscoveryNmapDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting nmap autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/nmap/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryNmapDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryNmapFormData builds the create/update payload. The API upserts by
// name: re-posting the same name updates the existing job, so the form data is
// identical for create and update.
func autoDiscoveryNmapFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("server"); ok {
		formData["server"] = v.(string)
	}
	if v, ok := d.GetOk("start_ip"); ok {
		formData["start_ip"] = v.(string)
	}
	if v, ok := d.GetOk("end_ip"); ok {
		formData["end_ip"] = v.(string)
	}
	if v, ok := d.GetOk("target"); ok {
		formData["target"] = v.(string)
	}
	if v, ok := d.GetOk("exclude_servers"); ok {
		formData["exclude_servers"] = v.(string)
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
