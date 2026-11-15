package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiVrfGroupListResponse struct {
	VrfGroups []apiVrfGroup `json:"vrfgroup"`
}

func datasourceD42VrfGroup() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42VrfGroupRead,
		Description: "Read a Device42 VRF group.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the VRF group.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the VRF group.",
			},
			"description": {
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

func datasourceD42VrfGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	if name == "" && id == 0 {
		return fmt.Errorf("either 'name' or 'id' must be set to look up a Device42 VRF group")
	}

	var resp apiVrfGroupListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/vrfgroup/")
	if err != nil {
		log.Printf("[WARN] No vrf group found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42VrfGroupRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiVrfGroup
	for i := range resp.VrfGroups {
		g := resp.VrfGroups[i]
		if id > 0 && int(g.Id) == id {
			match = &resp.VrfGroups[i]
			break
		}
		if name != "" && g.Name == name {
			match = &resp.VrfGroups[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No vrf group found with name: %s id: %d", name, id)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.Id, 10))
	d.Set("id", match.Id)
	d.Set("name", match.Name)
	d.Set("description", match.Description)
	d.Set("tags", joinBuildingList(match.Tags))
	d.Set("custom_fields", flattenCustomFields(match.CustomFields))

	return nil
}
