package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42SubnetCategory() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SubnetCategoryRead,
		Description: "Read a Device42 subnet category.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the subnet category.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the subnet category.",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42SubnetCategoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	if name == "" && id == 0 {
		return fmt.Errorf("either 'name' or 'id' must be set to look up a Device42 subnet category")
	}

	var resp apiSubnetCategoryListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/subnet_category/")
	if err != nil {
		log.Printf("[WARN] No subnet category found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42SubnetCategoryRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiSubnetCategory
	for i := range resp.Categories {
		c := resp.Categories[i]
		if id > 0 && int(c.Id) == id {
			match = &resp.Categories[i]
			break
		}
		if name != "" && c.Name == name {
			match = &resp.Categories[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No subnet category found with name: %s id: %d", name, id)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.Id, 10))
	d.Set("id", match.Id)
	d.Set("name", match.Name)
	d.Set("description", match.Description)
	d.Set("groups", match.Groups)

	return nil
}
