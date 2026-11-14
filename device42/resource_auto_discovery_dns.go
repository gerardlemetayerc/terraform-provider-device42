package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryDns struct {
	JobId        int64  `json:"job_id"`
	ZoneName     string `json:"zonename"`
	Nameserver   string `json:"nameserver"`
	Schedule     string `json:"schedule"`
	ScheduleTime string `json:"schedule_time"`
	ScheduleDays string `json:"schedule_days"`
}

type apiAutoDiscoveryDnsListResponse struct {
	Jobs []apiAutoDiscoveryDns `json:"jobs"`
}

func resourceD42AutoDiscoveryDns() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_dns can be use to manage a DNS autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryDnsCreate,
		Read:        resourceDevice42AutoDiscoveryDnsRead,
		Delete:      resourceDevice42AutoDiscoveryDnsDelete,

		Schema: map[string]*schema.Schema{
			"zonename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nameserver": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDevice42AutoDiscoveryDnsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryDnsFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/dns/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryDnsCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryDnsCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryDnsRead(d, m)
}

func resourceDevice42AutoDiscoveryDnsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryDnsRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryDnsListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/dns/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryDnsRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryDnsRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("zonename", j.ZoneName)
			d.Set("nameserver", j.Nameserver)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryDnsRead - no dns job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryDnsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting dns autodiscovery job id %s (zonename: %s)", d.Id(), d.Get("zonename"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/dns/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryDnsDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryDnsFormData builds the create/update payload. The API upserts by
// zonename: re-posting the same zonename updates the existing job, so the form
// data is identical for create and update.
func autoDiscoveryDnsFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"zonename": d.Get("zonename").(string),
	}

	if v, ok := d.GetOk("nameserver"); ok {
		formData["nameserver"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}

	return formData
}
