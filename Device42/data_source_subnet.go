package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42SubnetResponse struct {
	TotalCount int                     `json:"total_count"`
	Subnets    []apiSubnetReadResponse `json:"subnets"`
}

func datasourceD42Subnet() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SubnetRead,
		Description: "Read device.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the subnet.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The subnet id.",
			},
			"range_begin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"range_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42SubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	subnet_id := d.Get("subnet_id").(int)
	queryString := ""
	separator := ""
	if name != "" {
		queryString = fmt.Sprintf("name=%s", name)
		separator = "&"
	}
	if subnet_id > 0 {
		queryString = queryString + separator + fmt.Sprintf("subnet_id=%s", strconv.Itoa(subnet_id))
	}

	resp, err := client.R().
		SetResult(datasourceD42SubnetResponse{}).
		Get(fmt.Sprintf("/1.0/subnets/?%s", queryString))
	log.Printf("[DEBUG] targetURl: %s", fmt.Sprintf("/2.0/devices/?name=%s", d.Get("name").(string)))
	if err != nil {
		log.Printf("[WARN] No subnet found: %s", d.Id())
		d.SetId("")
		return nil
	}

	r := resp.Result().(*datasourceD42SubnetResponse)
	log.Printf("[DEBUG] Result: %#v", resp.Result())
	if len(r.Subnets) == 1 {
		d.SetId(strconv.Itoa(int((r.Subnets[0]).Subnet_id)))
		d.Set("subnet_id", (r.Subnets[0]).Subnet_id)
		d.Set("name", r.Subnets[0].Name)
		d.Set("range_begin", r.Subnets[0].RangeBegin)
		d.Set("range_end", r.Subnets[0].RangeEnd)
	}
	return nil
}
