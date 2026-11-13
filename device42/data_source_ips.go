package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Ips() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42IpsRead,
		Description: "Read Device42 IPs of a subnet (or list up to 1000 IPs).",
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter IPs belonging to the given subnet id (server-side filter).",
			},
			"ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IPs matching the filter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 ID of the IP.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"available": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the IP is marked as available.",
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
						"device_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the device attached to the IP, if any.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42IpsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	// /2.0/ips/ is paginated (max 1000 records per response), so filtering must
	// be done server-side via ?subnet_id= rather than client-side. Without a
	// filter only the first 1000 IPs are returned.
	subnetID := d.Get("subnet_id").(int)
	path := "/2.0/ips/"
	key := "all"
	if subnetID > 0 {
		path = fmt.Sprintf("/2.0/ips/?subnet_id=%d", subnetID)
		key = "subnet-" + strconv.Itoa(subnetID)
	}
	log.Printf("[DEBUG] datasourceD42IpsRead - querying %s", path)

	var resp apiReadData
	_, err := client.R().
		SetResult(&resp).
		Get(path)
	if err != nil {
		log.Printf("[WARN] datasourceD42IpsRead - error querying IPs: %s", err)
		d.SetId("")
		return nil
	}

	ips := make([]interface{}, 0, len(resp.Ips))
	for _, ip := range resp.Ips {
		var deviceID int32
		if len(ip.Devices) > 0 {
			deviceID = ip.Devices[0].DeviceID
		}

		ips = append(ips, map[string]interface{}{
			"ip_id":          ip.Id,
			"ip":             ip.Ip,
			"available":      ip.Available,
			"subnet_id":      ip.Subnet_id,
			"subnet":         ip.Subnet,
			"vrf_group_name": ip.VrfGroupName,
			"mac_address":    ip.Mac_Address,
			"label":          ip.Label,
			"notes":          ip.Notes,
			"device_id":      deviceID,
		})
	}
	d.Set("ips", ips)
	d.SetId(key)
	return nil
}
