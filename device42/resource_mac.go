package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiMac struct {
	MacaddressId int64  `json:"macaddress_id"`
	Macaddress   string `json:"macaddress"`
	VlanId       int64  `json:"vlan_id"`
	Port         string `json:"port"`
	PortName     string `json:"port_name"`
}

func resourceD42Mac() *schema.Resource {
	return &schema.Resource{
		Description: "device42_mac can be use to manage MAC addresses.",
		Create:      resourceDevice42MacCreate,
		Read:        resourceDevice42MacRead,
		Update:      resourceDevice42MacUpdate,
		Delete:      resourceDevice42MacDelete,

		Schema: map[string]*schema.Schema{
			"macaddress": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDevice42MacCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := macFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/macs/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42MacCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42MacCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42MacRead(d, m)
}

func resourceDevice42MacRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42MacRead - Starting reading using API for id %s", d.Id())

	var resp apiMac
	apiResp, err := client.R().
		SetResult(&resp).
		Get(fmt.Sprintf("/1.0/macs/%s/", d.Id()))
	if err != nil {
		log.Printf("[WARN] No mac found for id %s", d.Id())
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() == 404 {
		log.Printf("[WARN] No mac found for id %s", d.Id())
		d.SetId("")
		return nil
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42MacRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	mac := resp
	d.Set("macaddress", mac.Macaddress)
	d.Set("vlan_id", strconv.FormatInt(mac.VlanId, 10))
	d.Set("port", mac.Port)

	return nil
}

func resourceDevice42MacUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := macFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/macs/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42MacUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42MacUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42MacRead(d, m)
}

func resourceDevice42MacDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting mac id %s (macaddress: %s)", d.Id(), d.Get("macaddress"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/macs/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42MacDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func macFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"macaddress": d.Get("macaddress").(string),
	}

	if v, ok := d.GetOk("vlan_id"); ok {
		formData["vlan_id"] = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		formData["port"] = v.(string)
	}

	return formData
}
