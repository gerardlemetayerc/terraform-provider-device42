package device42

import (
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Vlans() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42VlansRead,
		Description: "Read multiple Device42 VLANs, optionally filtered by name or id.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter VLANs by exact name.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter VLANs by id.",
			},
			"vlans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "VLANs matching the filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 VLAN id.",
						},
						"number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VLAN number.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VLAN.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the VLAN.",
						},
						"tags": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tags of the VLAN.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42VlansRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42VlansRead - querying VLANs")

	var resp datasourceD42VlanResponse
	_, err := client.R().
		SetResult(&resp).
		Get("/1.0/vlans/")
	if err != nil {
		log.Printf("[WARN] datasourceD42VlansRead - error querying VLANs: %s", err)
		d.SetId("")
		return nil
	}

	name := d.Get("name").(string)
	vlanID := int64(d.Get("vlan_id").(int))

	vlans := make([]interface{}, 0)
	for _, v := range resp.Vlans {
		if name != "" && v.Name != name {
			continue
		}
		if vlanID > 0 && v.VlanId != vlanID {
			continue
		}

		vlans = append(vlans, map[string]interface{}{
			"vlan_id":     v.VlanId,
			"number":      v.Number,
			"name":        v.Name,
			"description": v.Description,
			"tags":        joinBuildingList(v.Tags),
		})
	}
	d.Set("vlans", vlans)

	key := "all"
	if vlanID > 0 {
		key = "vlan-" + strconv.FormatInt(vlanID, 10)
	}
	if name != "" {
		key = key + "-" + name
	}
	d.SetId(key)
	return nil
}
