package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryBlade struct {
	JobId           int64         `json:"job_id"`
	Name            string        `json:"name"`
	Server          string        `json:"server"`
	StartIp         string        `json:"start_ip"`
	EndIp           string        `json:"end_ip"`
	ExcludeServers  string        `json:"exclude_servers"`
	SnmpRetries     int           `json:"snmp_retries"`
	SnmpTimeout     int           `json:"snmp_timeout"`
	SnmpMinParallel int           `json:"snmp_min_parallel"`
	DoPortCheck     interface{}   `json:"do_port_check"`
	ClearCustomer   interface{}   `json:"clear_customer"`
	ScheduleTime    string        `json:"schedule_time"`
	ScheduleDays    string        `json:"schedule_days"`
	SnmpV3          string        `json:"snmp_v3"`
	Tags            []interface{} `json:"tags"`
}

type apiAutoDiscoveryBladeListResponse struct {
	Jobs []apiAutoDiscoveryBlade `json:"jobs"`
}

func resourceD42AutoDiscoveryBlade() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_blade can be use to manage a blade autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryBladeCreate,
		Read:        resourceDevice42AutoDiscoveryBladeRead,
		Update:      resourceDevice42AutoDiscoveryBladeUpdate,
		Delete:      resourceDevice42AutoDiscoveryBladeDelete,

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
			"exclude_servers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snmp_retries": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snmp_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snmp_min_parallel": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"do_port_check": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clear_customer": {
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
			"snmp_v3": {
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

func resourceDevice42AutoDiscoveryBladeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryBladeFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/blade/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryBladeRead(d, m)
}

func resourceDevice42AutoDiscoveryBladeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryBladeRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryBladeListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/blade/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("server", j.Server)
			d.Set("start_ip", j.StartIp)
			d.Set("end_ip", j.EndIp)
			d.Set("exclude_servers", j.ExcludeServers)
			d.Set("snmp_retries", j.SnmpRetries)
			d.Set("snmp_timeout", j.SnmpTimeout)
			d.Set("snmp_min_parallel", j.SnmpMinParallel)
			d.Set("do_port_check", fmt.Sprintf("%v", j.DoPortCheck))
			d.Set("clear_customer", fmt.Sprintf("%v", j.ClearCustomer))
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("snmp_v3", j.SnmpV3)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryBladeRead - no blade job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryBladeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryBladeFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/blade/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryBladeRead(d, m)
}

func resourceDevice42AutoDiscoveryBladeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting blade autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/blade/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryBladeDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryBladeFormData builds the create/update payload. The API upserts
// by name: re-posting the same name updates the existing job, so the form data
// is identical for create and update.
func autoDiscoveryBladeFormData(d *schema.ResourceData) map[string]string {
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
	if v, ok := d.GetOk("exclude_servers"); ok {
		formData["exclude_servers"] = v.(string)
	}
	if v, ok := d.GetOk("snmp_retries"); ok {
		formData["snmp_retries"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("snmp_timeout"); ok {
		formData["snmp_timeout"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("snmp_min_parallel"); ok {
		formData["snmp_min_parallel"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("do_port_check"); ok {
		formData["do_port_check"] = v.(string)
	}
	if v, ok := d.GetOk("clear_customer"); ok {
		formData["clear_customer"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}
	if v, ok := d.GetOk("snmp_v3"); ok {
		formData["snmp_v3"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
