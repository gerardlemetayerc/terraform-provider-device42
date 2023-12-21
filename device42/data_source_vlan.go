package device42

import (
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42VlanResponse struct {
	TotalCount int                     `json:"total_count"`
	Vlans      []apiVlanReadResponse `json:"vlans"`
}

func datasourceD42Vlan() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SubnetRead,
		Description: "Read vlan information.",
		Schema: map[string]*schema.Schema{
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vlan",
			},
			"number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "vlan number.",
			},
		},
	}
}

func datasourceD42VlanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	vlan_id := d.Get("vlan_id").(int)
	number  := d.Get("number").(int)
	queryParams := make(url.Values)
  
	if vlan_id > 0 {
		queryParams.Set("vlan_id", vlan_id)
	}
	if number > 0 {
		queryParams.Set("number", strconv.Itoa(number))
	}
  
	client.SetDebug(true)
	resp, err := client.R().
		SetResult(datasourceD42SubnetResponse{}).
		SetHeader("Accept", "application/json").
		SetQueryParamsFromValues(queryParams).
		Get("/1.0/vlans/?" + queryParams.Encode())
	log.Printf("[DEBUG] targetURl: %s", "/1.0/vlans/")
	if err != nil {
		log.Printf("[WARN] No vlans found: %s", d.Id())
		log.Printf("[WARN] No vlans found: %v", err)
		d.SetId("")
		return nil
	}

	r := resp.Result().(*datasourceD42VlanResponse)
	log.Printf("[DEBUG] Result: %#v", resp.Result())
	log.Printf("[DEBUG] Subnets count: %#v", r.Vlans)
	if len(r.Subnets) == 1 {
		d.SetId(strconv.Itoa(int((r.Vlans[0]).VlanId)))
		d.Set("number", (r.Vlans[0]).Number)
		d.Set("vlan_id", strconv.Itoa(int((r.Vlans[0]).VlanId)))
	} else {
		log.Printf("[ERROR] More than one subnet found: %d", len(r.Subnets))
		d.SetId("")
		return nil
	}
	return nil
}
