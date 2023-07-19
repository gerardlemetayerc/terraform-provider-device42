package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiReadData struct {
	TotalCount int32               `json:"total_count"`
	Ips        []apiIPReadResponse `json:"ips"`
}

type apiIPReadResponse struct {
	Available    string `json:"available"`
	Id           int32  `json:"id"`
	Ip           string `json:"ip"`
	Label        string `json:"label"`
	Mac_Address  int64  `json:"mac_address"`
	Notes        string `json:"notes"`
	Subnet       string `json:"subnet"`
	Subnet_id    int32  `json:"subnet_id"`
	VrfGroupName string `json:"vrf_group_name"`
}

func resourceD42Ip() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42IpCreate,
		Read:   resourceDevice42IpRead,
		Update: resourceDevice42IpCreate,
		Delete: resourceDevice42IpDelete,

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Network of the subnet. Required for creation, cannot be modified after subnet creation.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet name of the IP.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subnet ID.",
			},
			"available": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet VRF Group",
			},
			"vrf_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subnet VRF Group ID",
			},
		},
	}
}

func resourceDevice42IpRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42IpRead - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiReadData{}).
		Get(fmt.Sprintf("/2.0/ips/?ip_id=%s", d.Id()))

	if err != nil {
		log.Printf("[WARN] No ip found for id %s", d.Id())
		return err
	}

	r := resp.Result().(*apiReadData)
	str := fmt.Sprintf("%v", r)
	log.Printf("[DEBUG] resourceDevice42IpRead - API data %s", str)
	d.Set("available", r.Ips[0].Available)
	d.Set("ip", r.Ips[0].Ip)
	d.Set("subnet", r.Ips[0].Subnet)

	return nil
}

func resourceDevice42IpCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	ip := d.Get("ip").(string)
	subnet := d.Get("subnet").(string)
	available := d.Get("available").(string)
	vrf_group_id := d.Get("vrf_group_id").(int)

	mapData := map[string]string{
		"ipaddress": ip,
	}

	if subnet != "" {
		mapData["subnet"] = subnet
	}

	if available != "" {
		mapData["available"] = available
	}

	if vrf_group_id > 0 {
		mapData["vrf_group_id"] = strconv.Itoa(int(vrf_group_id))
	}

	resp, err := client.R().
		SetFormData(mapData).
		SetResult(apiResponse{}).
		Post("/2.0/ips/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	// Set ID after subnet creation
	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceDevice42IpDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting vlan %s (UUID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/ips/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(url)

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}
