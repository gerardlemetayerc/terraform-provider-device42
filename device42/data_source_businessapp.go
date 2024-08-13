package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42BusinessAppApi struct {
	TotalCount   int                           `json:"total_count"`
	Businessapps []apiBusinessAppsReadResponse `json:"businessapps"`
}

func datasourceD42BusinessApp() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42BusinessAppRead,
		Description: "Read device.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This name of the business application.",
			},
			"custom_fields": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Any custom fields that will be used in device42.",
			},
		},
	}
}

func datasourceD42BusinessAppRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42BusinessAppRead - Query: %s", fmt.Sprintf("/1.0/businessapps/?name=%s", d.Get("name").(string)))
	resp, err := client.R().
		SetResult(datasourceD42BusinessAppApi{}).
		Get(fmt.Sprintf("/1.0/businessapps/?name=%s", d.Get("name").(string)))

	if err != nil {
		log.Printf("[WARN] No businessapps found: %s", d.Get("name").(string))
		d.SetId("")
		return nil
	}

	r := resp.Result().(*datasourceD42BusinessAppApi)
	log.Printf("[DEBUG] datasourceD42BusinessAppRead - Result: %#v", resp.Result())
	if len(r.Businessapps) == 1 {
		d.SetId(strconv.Itoa(int((r.Businessapps[0]).BusinessAppId)))
		fields := flattenCustomFields((r.Businessapps[0]).CustomFields)
		d.Set("custom_fields", fields)
	}
	return nil
}
