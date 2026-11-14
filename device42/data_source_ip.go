package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Ip() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42IpRead,
		Description: "Read details of an existing Device42 IP by address.",
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: IsIPAddress,
				Description:  "The IP address to look up.",
			},
			"ip_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Device42 ID of the IP.",
			},
			"available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the IP is marked as available.",
			},
			"device_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the device attached to the IP.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Subnet ID the IP belongs to.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet network the IP belongs to.",
			},
			"vrf_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VRF group name of the subnet.",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address associated with the IP.",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Label of the IP.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes for the IP.",
			},
		},
	}
}

func datasourceD42IpRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	ip := d.Get("ip").(string)

	log.Printf("[DEBUG] datasourceD42IpRead - Querying IP on address %s", ip)

	// GET /2.0/ips/ supports filtering by the exact IP address via ?ip=.
	path := fmt.Sprintf("/2.0/ips/?ip=%s", ip)
	resp, err := apiDevice42Get(client, path, apiReadData{})
	if err != nil {
		log.Printf("[WARN] datasourceD42IpRead - no IP found for %s: %s", ip, err)
		d.SetId("")
		return nil
	}

	r := resp.Result().(*apiReadData)
	if len(r.Ips) == 0 {
		log.Printf("[WARN] datasourceD42IpRead - no IP found for %s", ip)
		d.SetId("")
		return nil
	}

	// The API filters by exact address so at most one match is expected.
	if len(r.Ips) > 1 {
		return fmt.Errorf("datasourceD42IpRead - more than one IP found for %s", ip)
	}

	ipData := r.Ips[0]
	d.SetId(strconv.Itoa(int(ipData.Id)))
	d.Set("ip_id", ipData.Id)
	d.Set("available", ipData.Available)
	d.Set("subnet_id", ipData.Subnet_id)
	d.Set("subnet", ipData.Subnet)
	d.Set("vrf_group_name", ipData.VrfGroupName)
	d.Set("mac_address", ipData.Mac_Address)
	d.Set("label", ipData.Label)
	d.Set("notes", ipData.Notes)
	if len(ipData.Devices) > 0 {
		d.Set("device_id", ipData.Devices[0].DeviceID)
	}

	return nil
}
