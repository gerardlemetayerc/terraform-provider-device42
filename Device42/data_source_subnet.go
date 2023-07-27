package device42

import (
	"log"
	"net/url"
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
				Description: "The name of the subnet.",
			},
			"mask_bits": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Mask CIDR of the subnet.",
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
			"vrf_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceD42SubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	vrf_group_name := d.Get("vrf_group_name").(string)
	subnet_id := d.Get("subnet_id").(int)
	mask_bits := d.Get("mask_bits").(int)
	network := d.Get("network").(string)
	queryParams := make(url.Values)
	if name != "" {
		queryParams.Set("name", name)
	}
	if subnet_id > 0 {
		queryParams.Set("subnet_id", strconv.Itoa(subnet_id))
	}
	if mask_bits > 0 {
		queryParams.Set("mask_bits", strconv.Itoa(mask_bits))
	}
	if vrf_group_name != "" {
		queryParams.Set("vrf_group", vrf_group_name)
	}
	if network != "" {
		queryParams.Set("network", network)

	}
	client.SetDebug(true)
	resp, err := client.R().
		SetResult(datasourceD42SubnetResponse{}).
		SetHeader("Accept", "application/json").
		SetQueryParamsFromValues(queryParams).
		Get("/1.0/subnets/?" + queryParams.Encode())
	log.Printf("[DEBUG] targetURl: %s", "/1.0/subnets/")
	if err != nil {
		log.Printf("[WARN] No subnet found: %s", d.Id())
		log.Printf("[WARN] No subnet found: %v", err)
		d.SetId("")
		return nil
	}

	r := resp.Result().(*datasourceD42SubnetResponse)
	log.Printf("[DEBUG] Result: %#v", resp.Result())
	log.Printf("[DEBUG] Subnets count: %#v", r.Subnets)
	if len(r.Subnets) == 1 {
		d.SetId(strconv.Itoa(int((r.Subnets[0]).Subnet_id)))
		d.Set("subnet_id", (r.Subnets[0]).Subnet_id)
		d.Set("name", r.Subnets[0].Name)
		d.Set("range_begin", r.Subnets[0].RangeBegin)
		d.Set("range_end", r.Subnets[0].RangeEnd)
		d.Set("vrf_group_id", r.Subnets[0].VrfGroupId)
		d.Set("vrf_group_name", r.Subnets[0].VrfGroupName)
	} else {
		log.Printf("[ERROR] More than one subnet found: %d", len(r.Subnets))
		d.SetId("")
		return nil
	}
	return nil
}
