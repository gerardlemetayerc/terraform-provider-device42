package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42DeviceResponse struct {
	TotalCount int                     `json:"total_count"`
	Devices    []apiDeviceReadResponse `json:"devices"`
}

func datasourceD42Device() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42DeviceRead,
		Description: "Read device.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the device.",
			},
			"device_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The hostname of the device.",
			},
		},
	}
}

func datasourceD42DeviceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetResult(datasourceD42DeviceResponse{}).
		Get(fmt.Sprintf("/2.0/devices/?name=%s", d.Get("name").(string)))
	log.Printf("[DEBUG] targetURl: %s", fmt.Sprintf("/2.0/devices/?name=%s", d.Get("name").(string)))
	if err != nil {
		log.Printf("[WARN] No device found: %s", d.Id())
		d.SetId("")
		return nil
	}

	r := resp.Result().(*datasourceD42DeviceResponse)
	log.Printf("[DEBUG] Result: %#v", resp.Result())
	if r.TotalCount == 1 {
		d.SetId(strconv.Itoa(int((r.Devices[0]).DeviceID)))
		d.Set("device_id", (r.Devices[0]).DeviceID)
		d.Set("name", r.Devices[0].Name)
	}
	return nil
}
