package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42SuggestedIpResponse struct {
	Ip string `json:"ip"`
}

func datasourceD42SuggestedIp() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SuggestedIpRead,
		Description: "Read suggested IP.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP as an ID.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the subnet.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The id of the subnet.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Calculated free IP",
			},
		},
	}
}

func datasourceD42SuggestedIpRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("subnet_name").(string)
	subnet_id := d.Get("subnet_id").(int)
	queryString := ""
	separator := ""
	if name != "" {
		queryString = fmt.Sprintf("name=%s", name)
		separator = "&"
		d.Set("name", name)
	}
	if subnet_id > 0 {
		queryString = queryString + separator + fmt.Sprintf("subnet_id=%s", strconv.Itoa(subnet_id))
		d.Set("subnet_id", subnet_id)
	}
	state := d.State()
	log.Printf("[DEBUG] Current ID: %#v", state)
	resp, err := client.R().
		SetResult(datasourceD42SuggestedIpResponse{}).
		Get(fmt.Sprintf("/1.0/suggest_ip/?%s", queryString))
	log.Printf("[DEBUG] targetURl: %s", fmt.Sprintf("/1.0/suggest_ip/?%s", queryString))
	if err != nil {
		log.Printf("[WARN] No ip found: %s", d.Id())

		return nil
	}

	r := resp.Result().(*datasourceD42SuggestedIpResponse)
	log.Printf("[DEBUG] Result: %#v", resp.Result())
	d.Set("ip", r.Ip)
	d.SetId(r.Ip)
	return nil
}
