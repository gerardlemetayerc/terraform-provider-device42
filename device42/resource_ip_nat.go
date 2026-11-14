package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiIpNat struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	IpAddressFrom    string `json:"ip_address_from"`
	IpAddressTo      string `json:"ip_address_to"`
	Notes            string `json:"notes"`
	Protocol         string `json:"protocol"`
	TwoWayRelation   string `json:"two_way_relation"`
	SourcePortStart  string `json:"source_port_start"`
	SourcePortEnd    string `json:"source_port_end"`
	TargetPortStart  string `json:"target_port_start"`
	TargetPortEnd    string `json:"target_port_end"`
	IpAddressFromEnd string `json:"ip_address_from_end"`
	IpAddressToEnd   string `json:"ip_address_to_end"`
}

type apiIpNatListResponse struct {
	Ipnats []apiIpNat `json:"ipnats"`
}

func resourceD42IpNat() *schema.Resource {
	return &schema.Resource{
		Description: "device42_ip_nat can be use to manage IP NAT records.",
		Create:      resourceDevice42IpNatCreate,
		Read:        resourceDevice42IpNatRead,
		Update:      resourceDevice42IpNatUpdate,
		Delete:      resourceDevice42IpNatDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address_from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"two_way_relation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_port_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_port_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_port_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_from_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_to_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDevice42IpNatCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := ipNatFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/ipnat/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42IpNatCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42IpNatCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42IpNatRead(d, m)
}

func resourceDevice42IpNatRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42IpNatRead - Starting reading using API for id %s", d.Id())

	var resp apiIpNatListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/ipnat/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42IpNatRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42IpNatRead - invalid id %q: %s", d.Id(), err)
	}

	for _, n := range resp.Ipnats {
		if int(n.Id) == id {
			d.Set("name", n.Name)
			d.Set("ip_address_from", n.IpAddressFrom)
			d.Set("ip_address_to", n.IpAddressTo)
			d.Set("notes", n.Notes)
			d.Set("protocol", n.Protocol)
			d.Set("two_way_relation", n.TwoWayRelation)
			d.Set("source_port_start", n.SourcePortStart)
			d.Set("source_port_end", n.SourcePortEnd)
			d.Set("target_port_start", n.TargetPortStart)
			d.Set("target_port_end", n.TargetPortEnd)
			d.Set("ip_address_from_end", n.IpAddressFromEnd)
			d.Set("ip_address_to_end", n.IpAddressToEnd)
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42IpNatRead - no ip nat found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42IpNatUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	payload := map[string]string{
		"name":            d.Get("name").(string),
		"ip_address_from": d.Get("ip_address_from").(string),
		"ip_address_to":   d.Get("ip_address_to").(string),
	}
	if v, ok := d.GetOk("two_way_relation"); ok {
		payload["two_way_relation"] = v.(string)
	} else {
		payload["two_way_relation"] = "False"
	}
	if v, ok := d.GetOk("notes"); ok {
		payload["notes"] = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		payload["protocol"] = v.(string)
	}
	if v, ok := d.GetOk("source_port_start"); ok {
		payload["source_port_start"] = v.(string)
	}
	if v, ok := d.GetOk("source_port_end"); ok {
		payload["source_port_end"] = v.(string)
	}
	if v, ok := d.GetOk("target_port_start"); ok {
		payload["target_port_start"] = v.(string)
	}
	if v, ok := d.GetOk("target_port_end"); ok {
		payload["target_port_end"] = v.(string)
	}
	if v, ok := d.GetOk("ip_address_from_end"); ok {
		payload["ip_address_from_end"] = v.(string)
	}
	if v, ok := d.GetOk("ip_address_to_end"); ok {
		payload["ip_address_to_end"] = v.(string)
	}

	resp, err := client.R().
		SetResult(apiResponse{}).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put(fmt.Sprintf("/1.0/ipnat/%s/", d.Id()))

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42IpNatUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42IpNatRead(d, m)
}

func resourceDevice42IpNatDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting ip nat id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/ipnat/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42IpNatDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func ipNatFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name":            d.Get("name").(string),
		"ip_address_from": d.Get("ip_address_from").(string),
		"ip_address_to":   d.Get("ip_address_to").(string),
	}

	if v, ok := d.GetOk("notes"); ok {
		formData["notes"] = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		formData["protocol"] = v.(string)
	}
	if v, ok := d.GetOk("two_way_relation"); ok {
		formData["two_way_relation"] = v.(string)
	}
	if v, ok := d.GetOk("source_port_start"); ok {
		formData["source_port_start"] = v.(string)
	}
	if v, ok := d.GetOk("source_port_end"); ok {
		formData["source_port_end"] = v.(string)
	}
	if v, ok := d.GetOk("target_port_start"); ok {
		formData["target_port_start"] = v.(string)
	}
	if v, ok := d.GetOk("target_port_end"); ok {
		formData["target_port_end"] = v.(string)
	}
	if v, ok := d.GetOk("ip_address_from_end"); ok {
		formData["ip_address_from_end"] = v.(string)
	}
	if v, ok := d.GetOk("ip_address_to_end"); ok {
		formData["ip_address_to_end"] = v.(string)
	}

	return formData
}
