package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiServiceLevel is the shape of a single service level object returned by
// GET /1.0/service_level/. The endpoint returns a flat JSON array (no wrapper
// object), and the ?name= query param is ignored so filtering is client-side.
type apiServiceLevel struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func datasourceD42ServiceLevel() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42ServiceLevelRead,
		Description: "Read Device42 service levels (tiers).",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of a specific service level to look up by exact match.",
			},
			"service_level_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Device42 ID of the service level matched by `name`.",
			},
			"levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All service levels available in Device42.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the service level.",
						},
						"service_level_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device42 ID of the service level.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42ServiceLevelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42ServiceLevelRead - querying service levels")

	// GET /1.0/service_level/ returns a flat JSON array of {name, id}. The
	// ?name= query param is ignored by the API, so any filtering is client-side.
	var levels []apiServiceLevel
	_, err := client.R().
		SetResult(&levels).
		Get("/1.0/service_level/")
	if err != nil {
		return err
	}

	// Always expose the full list so the data source can be used as a
	// list-all reference (the gap previously impossible to cover).
	list := make([]interface{}, 0, len(levels))
	for _, l := range levels {
		list = append(list, map[string]interface{}{
			"name":             l.Name,
			"service_level_id": l.ID,
		})
	}
	d.Set("levels", list)

	if name := d.Get("name").(string); name != "" {
		var matches []apiServiceLevel
		for _, l := range levels {
			if l.Name == name {
				matches = append(matches, l)
			}
		}

		if len(matches) == 0 {
			log.Printf("[WARN] datasourceD42ServiceLevelRead - no service level found with name: %s", name)
			d.SetId("")
			return nil
		}
		if len(matches) > 1 {
			return fmt.Errorf("datasourceD42ServiceLevelRead - more than one service level found with name: %s", name)
		}

		d.Set("service_level_id", matches[0].ID)
		d.SetId(strconv.Itoa(matches[0].ID))
		return nil
	}

	// List-all mode: no single id, use a stable placeholder.
	d.SetId("all")
	return nil
}