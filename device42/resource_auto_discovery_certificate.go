package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryCertificate struct {
	JobId           int64         `json:"job_id"`
	Name            string        `json:"name"`
	Server          string        `json:"server"`
	StartIpAddress  string        `json:"start_ip_address"`
	EndIpAddress    string        `json:"end_ip_address"`
	ExcludeServers  string        `json:"exclude_servers"`
	Ports           string        `json:"ports"`
	GetCipherSuites string        `json:"get_cipher_suites"`
	FollowChain     string        `json:"follow_chain"`
	DebugLevel      string        `json:"debug_level"`
	ScheduleTime    string        `json:"schedule_time"`
	ScheduleDays    string        `json:"schedule_days"`
	Tags            []interface{} `json:"tags"`
}

type apiAutoDiscoveryCertificateListResponse struct {
	Jobs []apiAutoDiscoveryCertificate `json:"jobs"`
}

func resourceD42AutoDiscoveryCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_certificate can be use to manage a certificate autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryCertificateCreate,
		Read:        resourceDevice42AutoDiscoveryCertificateRead,
		Update:      resourceDevice42AutoDiscoveryCertificateUpdate,
		Delete:      resourceDevice42AutoDiscoveryCertificateDelete,

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
			"start_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_servers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ports": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"get_cipher_suites": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"follow_chain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"debug_level": {
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

func resourceDevice42AutoDiscoveryCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryCertificateFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/certificate/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryCertificateRead(d, m)
}

func resourceDevice42AutoDiscoveryCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryCertificateRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryCertificateListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/certificate/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("server", j.Server)
			d.Set("start_ip_address", j.StartIpAddress)
			d.Set("end_ip_address", j.EndIpAddress)
			d.Set("exclude_servers", j.ExcludeServers)
			d.Set("ports", j.Ports)
			d.Set("get_cipher_suites", j.GetCipherSuites)
			d.Set("follow_chain", j.FollowChain)
			d.Set("debug_level", j.DebugLevel)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryCertificateRead - no certificate job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryCertificateFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/certificate/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryCertificateRead(d, m)
}

func resourceDevice42AutoDiscoveryCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting certificate autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/certificate/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCertificateDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryCertificateFormData builds the create/update payload. The API
// upserts by name: re-posting the same name updates the existing job, so the
// form data is identical for create and update.
func autoDiscoveryCertificateFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name":   d.Get("name").(string),
		"server": d.Get("server").(string),
	}

	if v, ok := d.GetOk("start_ip_address"); ok {
		formData["start_ip_address"] = v.(string)
	}
	if v, ok := d.GetOk("end_ip_address"); ok {
		formData["end_ip_address"] = v.(string)
	}
	if v, ok := d.GetOk("exclude_servers"); ok {
		formData["exclude_servers"] = v.(string)
	}
	if v, ok := d.GetOk("ports"); ok {
		formData["ports"] = v.(string)
	}
	if v, ok := d.GetOk("get_cipher_suites"); ok {
		formData["get_cipher_suites"] = v.(string)
	}
	if v, ok := d.GetOk("follow_chain"); ok {
		formData["follow_chain"] = v.(string)
	}
	if v, ok := d.GetOk("debug_level"); ok {
		formData["debug_level"] = v.(string)
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
