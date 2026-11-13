package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Subnets() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SubnetsRead,
		Description: "Read Device42 subnets, optionally filtered by VLAN, name or VRF group.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter subnets by exact name.",
			},
			"vrf_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter subnets by VRF group id.",
			},
			"vrf_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter subnets by VRF group name.",
			},
			"parent_vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter subnets attached to the given parent VLAN id.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Subnets matching the filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 subnet id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the subnet.",
						},
						"network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network address of the subnet.",
						},
						"mask_bits": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mask bits (CIDR) of the subnet.",
						},
						"gateway": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway of the subnet.",
						},
						"range_begin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "First usable IP of the subnet.",
						},
						"range_end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last usable IP of the subnet.",
						},
						"vrf_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VRF group id of the subnet.",
						},
						"vrf_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VRF group name of the subnet.",
						},
						"parent_vlan_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the parent VLAN the subnet is attached to.",
						},
						"parent_subnet_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the parent subnet, if any.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the subnet.",
						},
						"allocated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Allocation status of the subnet.",
						},
						"customer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Customer attached to the subnet.",
						},
						"custom_fields": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Custom fields of the subnet.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42SubnetsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42SubnetsRead - querying subnets")

	var resp datasourceD42SubnetResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/subnets/")
	if err != nil {
		log.Printf("[WARN] datasourceD42SubnetsRead - error querying subnets: %s", err)
		return err
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("datasourceD42SubnetsRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	name := d.Get("name").(string)
	vrfGroupID := int64(d.Get("vrf_group_id").(int))
	vrfGroupName := d.Get("vrf_group_name").(string)
	parentVlanID := int64(d.Get("parent_vlan_id").(int))

	subnets := make([]interface{}, 0)
	for _, s := range resp.Subnets {
		if name != "" && s.Name != name {
			continue
		}
		if vrfGroupID > 0 && int64(s.VrfGroupId) != vrfGroupID {
			continue
		}
		if vrfGroupName != "" && s.VrfGroupName != vrfGroupName {
			continue
		}
		if parentVlanID > 0 && int64(s.VlanId) != parentVlanID {
			continue
		}

		subnets = append(subnets, map[string]interface{}{
			"subnet_id":        s.Subnet_id,
			"name":             s.Name,
			"network":          s.Network,
			"mask_bits":        s.MaskBits,
			"gateway":          s.Gateway,
			"range_begin":      s.RangeBegin,
			"range_end":        s.RangeEnd,
			"vrf_group_id":     s.VrfGroupId,
			"vrf_group_name":   s.VrfGroupName,
			"parent_vlan_id":   s.VlanId,
			"parent_subnet_id": s.ParentSubnetId,
			"description":      s.Description,
			"allocated":        s.Allocated,
			"customer":         s.Customer,
			"custom_fields":    flattenCustomFields(s.CustomFields),
		})
	}
	d.Set("subnets", subnets)

	// Stable id derived from the filters so distinct queries get distinct ids.
	key := "all"
	if parentVlanID > 0 {
		key = "vlan-" + strconv.FormatInt(parentVlanID, 10)
	}
	if vrfGroupID > 0 {
		key = key + "-vrf-" + strconv.FormatInt(vrfGroupID, 10)
	}
	if vrfGroupName != "" {
		key = key + "-g-" + vrfGroupName
	}
	if name != "" {
		key = key + "-" + name
	}
	d.SetId(key)
	return nil
}
