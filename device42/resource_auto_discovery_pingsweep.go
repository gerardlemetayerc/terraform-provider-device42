package device42

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryPingSweep struct {
	JobId                     int64           `json:"job_id"`
	Name                      string          `json:"name"`
	Networks                  string          `json:"networks"`
	VrfGroup                  json.RawMessage `json:"vrfgroup"`
	ScheduleTime              string          `json:"schedule_time"`
	SubnetCategory            string          `json:"subnet_category"`
	CreateNewSubnet           string          `json:"create_new_subnet"`
	DebugLevel                string          `json:"debug_level"`
	OverwriteSubnetCategories string          `json:"overwrite_subnet_categories"`
	ScheduleDays              string          `json:"schedule_days"`
	ReverseDns                string          `json:"reverse_dns"`
	StripDomainSuffix         string          `json:"strip_domain_suffix"`
	Type                      string          `json:"type"`
	Schedule                  string          `json:"schedule"`
}

type apiAutoDiscoveryPingSweepListResponse struct {
	Jobs []apiAutoDiscoveryPingSweep `json:"jobs"`
}

func resourceD42AutoDiscoveryPingSweep() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_pingsweep can be use to manage a ping sweep autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryPingSweepCreate,
		Read:        resourceDevice42AutoDiscoveryPingSweepRead,
		Update:      resourceDevice42AutoDiscoveryPingSweepUpdate,
		Delete:      resourceDevice42AutoDiscoveryPingSweepDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"networks": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_new_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"debug_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"overwrite_subnet_categories": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reverse_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"strip_domain_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDevice42AutoDiscoveryPingSweepCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryPingSweepFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/pingsweep/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryPingSweepRead(d, m)
}

func resourceDevice42AutoDiscoveryPingSweepRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryPingSweepRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryPingSweepListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/pingsweep/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("networks", j.Networks)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("subnet_category", j.SubnetCategory)
			d.Set("create_new_subnet", j.CreateNewSubnet)
			d.Set("debug_level", j.DebugLevel)
			d.Set("overwrite_subnet_categories", j.OverwriteSubnetCategories)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("reverse_dns", j.ReverseDns)
			d.Set("strip_domain_suffix", j.StripDomainSuffix)
			d.Set("type", j.Type)
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryPingSweepRead - no pingsweep job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryPingSweepUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryPingSweepFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/pingsweep/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryPingSweepRead(d, m)
}

func resourceDevice42AutoDiscoveryPingSweepDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting pingsweep autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/pingsweep/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryPingSweepDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryPingSweepFormData builds the create/update payload. The API
// upserts by name: re-posting the same name updates the existing job, so the
// form data is identical for create and update.
func autoDiscoveryPingSweepFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("networks"); ok {
		formData["networks"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("subnet_category"); ok {
		formData["subnet_category"] = v.(string)
	}
	if v, ok := d.GetOk("create_new_subnet"); ok {
		formData["create_new_subnet"] = v.(string)
	}
	if v, ok := d.GetOk("debug_level"); ok {
		formData["debug_level"] = v.(string)
	}
	if v, ok := d.GetOk("overwrite_subnet_categories"); ok {
		formData["overwrite_subnet_categories"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}
	if v, ok := d.GetOk("reverse_dns"); ok {
		formData["reverse_dns"] = v.(string)
	}
	if v, ok := d.GetOk("strip_domain_suffix"); ok {
		formData["strip_domain_suffix"] = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		formData["type"] = v.(string)
	}

	return formData
}
