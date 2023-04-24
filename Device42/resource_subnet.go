package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceD42Subnet() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The hostname of the device.",
			},
			"vrf_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "d42null",
				Description: "Subnet VRF Group",
			},
			"parent_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent subnet id",
			},
			"customer": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Any custom fields that will be used in device42.",
			},
			"mask_bits": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"gateway": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"service_level": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"vlan": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceDevice42SubnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	deviceType := d.Get("type").(string)
	serviceLevel := d.Get("service_level").(string)

	resp, err := client.R().
		SetFormData(map[string]string{
			"name":          name,
			"type":          deviceType,
			"service_level": serviceLevel,
		}).
		SetResult(apiResponse{}).
		Post("/1.0/subnets/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	d.SetId(strconv.Itoa(id))

	return resourceDevice42DeviceRead(d, m)
}
