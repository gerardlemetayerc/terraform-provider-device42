package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Software() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SoftwareRead,
		Description: "Read Device42 software components, either by name or as a full list.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of a specific software component to look up by exact match.",
			},
			"software_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Device42 ID of the software component matched by `name`.",
			},
			"software": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All software components available in Device42.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the software component.",
						},
						"software_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 ID of the software component.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42SoftwareRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42SoftwareRead - querying software components")

	// GET /1.0/software/ returns ALL components in one response (the ?name= filter
	// is exact-match only), so any filtering is done client-side.
	var resp apiSoftwareListResponse
	_, err := client.R().
		SetResult(&resp).
		Get("/1.0/software/")
	if err != nil {
		return err
	}

	// Always expose the full list so the data source can be used as a list-all
	// reference (TypeList pattern).
	list := make([]interface{}, 0, len(resp.Software))
	for _, s := range resp.Software {
		list = append(list, map[string]interface{}{
			"name":        s.Name,
			"software_id": s.ID,
		})
	}
	d.Set("software", list)

	if name := d.Get("name").(string); name != "" {
		var matches []apiSoftwareReadResponse
		for _, s := range resp.Software {
			if s.Name == name {
				matches = append(matches, s)
			}
		}

		if len(matches) == 0 {
			log.Printf("[WARN] datasourceD42SoftwareRead - no software component found with name: %s", name)
			d.SetId("")
			return nil
		}
		if len(matches) > 1 {
			return fmt.Errorf("datasourceD42SoftwareRead - more than one software component found with name: %s", name)
		}

		d.Set("software_id", matches[0].ID)
		d.SetId(strconv.FormatInt(matches[0].ID, 10))
		return nil
	}

	// List-all mode: no single id, use a stable placeholder.
	d.SetId("all")
	return nil
}
