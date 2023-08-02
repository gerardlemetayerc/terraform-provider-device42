package device42

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiReadData struct {
	TotalCount int32               `json:"total_count"`
	Ips        []apiIPReadResponse `json:"ips"`
}

type apiIPReadDevice struct {
	DeviceID int32  `json:"device_id"`
	Name     string `json:"name"`
}

type apiIPReadResponse struct {
	Available    string            `json:"available"`
	Id           int32             `json:"id"`
	Ip           string            `json:"ip"`
	Label        string            `json:"label"`
	Mac_Address  int64             `json:"mac_address"`
	Notes        string            `json:"notes"`
	Subnet       string            `json:"subnet"`
	Subnet_id    int32             `json:"subnet_id"`
	VrfGroupName string            `json:"vrf_group_name"`
	Devices      []apiIPReadDevice `json:"devices"`
}

func IsIPAddress(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if net.ParseIP(v) == nil {
		errs = append(errs, fmt.Errorf("%q is not a valid IP address: %s", key, v))
	}
	return
}

func resourceD42Ip() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42IpCreate,
		Read:   resourceDevice42IpRead,
		Update: resourceDevice42IpCreate,
		Delete: resourceDevice42IpDelete,

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				ValidateFunc: IsIPAddress,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				Description:  "Network of the subnet. Required for creation, cannot be modified after subnet creation.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subnet ID.",
			},
			"available": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Set",
			},
			"vrf_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subnet VRF Group ID",
			},
			"device_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of device to attach to network",
			},
			"suggest_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use device42 IP suggestion functionnality to retrieve available IP.",
				Default:     false,
			},
		},
	}
}

func resourceDevice42IpRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] resourceDevice42IpRead - Starting reading using API for id %s", d.Id())

	resp, err := apiDevice42Get(m.(*resty.Client), fmt.Sprintf("/2.0/ips/?ip_id=%s", d.Id()), apiReadData{})

	if err != nil {
		return err
	}

	r := resp.Result().(*apiReadData)
	str := fmt.Sprintf("%v", r)
	log.Printf("[DEBUG] resourceDevice42IpRead - API data %s", str)
	if len(r.Ips) > 0 {
		if len(r.Ips[0].Devices) > 0 {
			d.Set("device_id", r.Ips[0].Devices[0].DeviceID)
		}
		d.Set("available", r.Ips[0].Available)
		d.Set("ip", r.Ips[0].Ip)
	}
	return nil
}

func resourceDevice42IpCreate(d *schema.ResourceData, m interface{}) error {
	ip := d.Get("ip").(string)
	available := d.Get("available").(string)
	vrf_group_id := d.Get("vrf_group_id").(int)
	device_id := d.Get("device_id").(int)

	if d.Get("suggest_ip").(bool) && d.Get("ip").(string) == "" {
		// Call the function to get a suggested IP (this function would be similar to the datasourceD42SuggestedIpRead)
		suggestedIp, err := getSuggestedIP(d, m)
		if err != nil {
			return err
		}
		ip = suggestedIp
		d.Set("ip", ip)
	}

	mapData := map[string]string{
		"ipaddress": ip,
	}

	if available != "" {
		mapData["available"] = available
	}

	if vrf_group_id > 0 {
		mapData["vrf_group_id"] = strconv.Itoa(int(vrf_group_id))
	}

	if device_id > 0 {
		mapData["device_ids"] = strconv.Itoa(int(device_id))
	}

	resp, err := apiDevice42Post(m.(*resty.Client), "/2.0/ips/", mapData, apiResponse{})

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42IpCreate - API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42IpCreate - Result: %#v", r)

	// Check if Device42 has answer correct response length
	if len(r.Msg) > 1 {
		id := int(r.Msg[1].(float64))
		// Set ID after subnet creation
		d.SetId(strconv.Itoa(id))
	} else {
		return fmt.Errorf("resourceDevice42IpCreate - API response incorrect: %v", r.Msg)
	}

	return nil
}

func resourceDevice42IpDelete(d *schema.ResourceData, m interface{}) error {
	resp, err := apiDevice42Delete(m.(*resty.Client), fmt.Sprintf("/1.0/ips/%s/", d.Id()), apiResponse{})

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// This function will query Device42 API to retrieve available IP.
func getSuggestedIP(d *schema.ResourceData, m interface{}) (string, error) {
	client := m.(*resty.Client)

	subnet_id := d.Get("subnet_id").(int)
	queryString := ""

	if subnet_id > 0 {
		queryString = fmt.Sprintf("subnet_id=%s", strconv.Itoa(subnet_id))
	}

	resp, err := client.R().
		SetResult(datasourceD42SuggestedIpResponse{}).
		Get(fmt.Sprintf("/1.0/suggest_ip/?%s", queryString))

	log.Printf("[DEBUG] getSuggestedIP - Query : %s", fmt.Sprintf("/1.0/suggest_ip/?%s", queryString))
	if err != nil {
		return "", fmt.Errorf("getSuggestedIP - failed to fetch suggested IP: %s", err)
	}

	r := resp.Result().(*datasourceD42SuggestedIpResponse)
	log.Printf("[DEBUG] getSuggestedIP - Result: %#v", r)
	if r.Ip == "" {
		return "", fmt.Errorf("no suggested IP returned from Device42 API")
	}

	return r.Ip, nil
}
