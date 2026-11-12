package device42

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiRackLookupResponse is the shape returned by GET /1.0/racks/?name=
// which is a bare array under the "racks" key (no total_count/offset).
type apiRackLookupResponse struct {
	Racks []apiRackReadResponse `json:"racks"`
}

func datasourceD42Rack() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42RackRead,
		Description: "Read a Device42 rack.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the rack.",
			},
			"rack_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the rack.",
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"building": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"room": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manufacturer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"row": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"orientation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
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

func datasourceD42RackRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	rackID := d.Get("rack_id").(int)

	if name != "" {
		url := fmt.Sprintf("/1.0/racks/?name=%s", url.QueryEscape(name))
		log.Printf("[DEBUG] datasourceD42RackRead - Query racks on URL : %s", url)

		var lookup apiRackLookupResponse
		_, err := client.R().
			SetResult(&lookup).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No rack found: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Racks) == 0 {
			log.Printf("[ERROR] No rack found with name: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Racks) > 1 {
			log.Printf("[ERROR] More than one rack found: %d", len(lookup.Racks))
			d.SetId("")
			return nil
		}

		r := lookup.Racks[0]
		d.SetId(strconv.FormatInt(r.RackID, 10))
		d.Set("rack_id", r.RackID)
		d.Set("name", r.Name)
		d.Set("size", r.Size)
		d.Set("building", r.Building)
		d.Set("room", r.Room)
		d.Set("manufacturer", r.Manufacturer)
		d.Set("notes", r.Notes)
		d.Set("row", r.Row)
		d.Set("orientation", r.Orientation)
		d.Set("tags", joinBuildingList(r.Tags))
		d.Set("custom_fields", flattenCustomFields(r.CustomFields))
		return nil
	}

	if rackID > 0 {
		url := fmt.Sprintf("/1.0/racks/%d/", rackID)
		log.Printf("[DEBUG] datasourceD42RackRead - Query rack on URL : %s", url)

		var resp apiRackReadResponse
		_, err := client.R().
			SetResult(&resp).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No rack found: %d", rackID)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.Itoa(rackID))
		d.Set("rack_id", resp.RackID)
		d.Set("name", resp.Name)
		d.Set("size", resp.Size)
		d.Set("building", resp.Building)
		d.Set("room", resp.Room)
		d.Set("manufacturer", resp.Manufacturer)
		d.Set("notes", resp.Notes)
		d.Set("row", resp.Row)
		d.Set("orientation", resp.Orientation)
		d.Set("tags", joinBuildingList(resp.Tags))
		d.Set("custom_fields", flattenCustomFields(resp.CustomFields))
		return nil
	}

	return fmt.Errorf("either 'name' or 'rack_id' must be set to look up a Device42 rack")
}
