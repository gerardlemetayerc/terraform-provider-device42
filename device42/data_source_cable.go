package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiCableListResponse struct {
	Cable []apiCable `json:"cable"`
}

func datasourceD42Cable() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42CableRead,
		Description: "Read a Device42 cable.",
		Schema: map[string]*schema.Schema{
			"cable_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of the cable.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the cable.",
			},
			"cable_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cable_length_units": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
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

func datasourceD42CableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	cableID := d.Get("cable_id").(string)
	id := d.Get("id").(int)

	if cableID == "" && id == 0 {
		return fmt.Errorf("either 'cable_id' or 'id' must be set to look up a Device42 cable")
	}

	var resp apiCableListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/cables/")
	if err != nil {
		log.Printf("[WARN] No cable found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42CableRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiCable
	for i := range resp.Cable {
		c := resp.Cable[i]
		if id > 0 && int(c.Id) == id {
			match = &resp.Cable[i]
			break
		}
		if cableID != "" && strings.EqualFold(c.CableId, cableID) {
			match = &resp.Cable[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No cable found with cable_id: %s id: %d", cableID, id)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.Id, 10))
	d.Set("id", match.Id)
	d.Set("cable_id", match.CableId)
	d.Set("cable_length", int64(match.CableLength))
	d.Set("cable_length_units", match.CableLengthUnits)
	d.Set("notes", match.Notes)
	d.Set("vendor", match.Vendor)
	d.Set("tags", joinBuildingList(match.Tags))
	d.Set("custom_fields", flattenCustomFields(match.CustomFields))

	return nil
}
