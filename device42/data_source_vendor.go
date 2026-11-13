package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Vendor() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42VendorRead,
		Description: "Read Device42 vendors, either by name or as a full list.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of a specific vendor to look up by exact match.",
			},
			"vendor_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Device42 ID of the vendor matched by `name`.",
			},
			"vendors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All vendors available in Device42.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the vendor.",
						},
						"vendor_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 ID of the vendor.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42VendorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42VendorRead - querying vendors")

	// GET /1.0/vendors/ returns ALL vendors (no filter param is honoured), so any
	// filtering is done client-side.
	var resp apiVendorsResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/vendors/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("datasourceD42VendorRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	// Always expose the full list so the data source can be used as a list-all
	// reference (TypeList pattern).
	list := make([]interface{}, 0, len(resp.Vendors))
	for _, v := range resp.Vendors {
		list = append(list, map[string]interface{}{
			"name":      v.Name,
			"vendor_id": v.VendorID,
		})
	}
	d.Set("vendors", list)

	if name := d.Get("name").(string); name != "" {
		var matches []apiVendorReadResponse
		for _, v := range resp.Vendors {
			if v.Name == name {
				matches = append(matches, v)
			}
		}

		if len(matches) == 0 {
			log.Printf("[WARN] datasourceD42VendorRead - no vendor found with name: %s", name)
			d.SetId("")
			return nil
		}
		if len(matches) > 1 {
			return fmt.Errorf("datasourceD42VendorRead - more than one vendor found with name: %s", name)
		}

		d.Set("vendor_id", matches[0].VendorID)
		d.SetId(strconv.FormatInt(matches[0].VendorID, 10))
		return nil
	}

	// List-all mode: no single id, use a stable placeholder.
	d.SetId("all")
	return nil
}
