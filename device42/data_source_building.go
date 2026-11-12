package device42

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiBuildingLookupResponse is the shape returned by GET /1.0/buildings/?name=
// which is a bare array under the "buildings" key (no total_count/offset).
type apiBuildingLookupResponse struct {
	Buildings []apiBuildingReadResponse `json:"buildings"`
}

func datasourceD42Building() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42BuildingRead,
		Description: "Read a Device42 building.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the building.",
			},
			"building_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the building.",
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_phone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func datasourceD42BuildingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	buildingID := d.Get("building_id").(int)

	if name != "" {
		url := fmt.Sprintf("/1.0/buildings/?name=%s", url.QueryEscape(name))
		log.Printf("[DEBUG] datasourceD42BuildingRead - Query buildings on URL : %s", url)

		var lookup apiBuildingLookupResponse
		_, err := client.R().
			SetResult(&lookup).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No building found: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Buildings) == 0 {
			log.Printf("[ERROR] No building found with name: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Buildings) > 1 {
			log.Printf("[ERROR] More than one building found: %d", len(lookup.Buildings))
			d.SetId("")
			return nil
		}

		b := lookup.Buildings[0]
		d.SetId(strconv.FormatInt(b.BuildingID, 10))
		d.Set("building_id", b.BuildingID)
		d.Set("name", b.Name)
		d.Set("address", b.Address)
		d.Set("contact_name", b.ContactName)
		d.Set("contact_phone", b.ContactPhone)
		d.Set("latitude", string(b.Latitude))
		d.Set("longitude", string(b.Longitude))
		d.Set("notes", b.Notes)
		d.Set("tags", joinBuildingList(b.Tags))
		d.Set("groups", joinBuildingList(b.Groups))
		d.Set("custom_fields", flattenCustomFields(b.CustomFields))
		return nil
	}

	if buildingID > 0 {
		url := fmt.Sprintf("/1.0/buildings/%d/", buildingID)
		log.Printf("[DEBUG] datasourceD42BuildingRead - Query building on URL : %s", url)

		var resp apiBuildingReadResponse
		_, err := client.R().
			SetResult(&resp).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No building found: %d", buildingID)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.Itoa(buildingID))
		d.Set("building_id", resp.BuildingID)
		d.Set("name", resp.Name)
		d.Set("address", resp.Address)
		d.Set("contact_name", resp.ContactName)
		d.Set("contact_phone", resp.ContactPhone)
		d.Set("latitude", string(resp.Latitude))
		d.Set("longitude", string(resp.Longitude))
		d.Set("notes", resp.Notes)
		d.Set("tags", joinBuildingList(resp.Tags))
		d.Set("groups", joinBuildingList(resp.Groups))
		d.Set("custom_fields", flattenCustomFields(resp.CustomFields))
		return nil
	}

	return fmt.Errorf("either 'name' or 'building_id' must be set to look up a Device42 building")
}
