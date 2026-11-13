package device42

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiSuggestSubnetResponse is the flat object returned by
// GET /1.0/suggest_subnet/{id}/. The swagger spec documents only `ip` and
// `mask`, so unknown fields are ignored.
type apiSuggestSubnetResponse struct {
	IP   string `json:"ip"`
	Mask int    `json:"mask"`
}

func datasourceD42SubnetSuggested() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SubnetSuggestedRead,
		Description: "Suggest the next available sub-subnet of a parent subnet.",
		Schema: map[string]*schema.Schema{
			"parent_subnet_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the parent subnet to suggest a child from.",
			},
			"mask_bits": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Mask bits of the suggested sub-subnet.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Suggested sub-subnet network address.",
			},
		},
	}
}

func datasourceD42SubnetSuggestedRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	parentSubnetID := d.Get("parent_subnet_id").(int)
	maskBits := d.Get("mask_bits").(int)

	path := fmt.Sprintf("/1.0/suggest_subnet/%d/?mask_bits=%d", parentSubnetID, maskBits)
	log.Printf("[DEBUG] datasourceD42SubnetSuggestedRead - querying %s", path)

	var result apiSuggestSubnetResponse
	resp, err := client.R().
		SetResult(&result).
		Get(path)
	if err != nil {
		log.Printf("[WARN] datasourceD42SubnetSuggestedRead - error: %s", err)
		d.SetId("")
		return nil
	}

	// The API returns HTTP 500 when no free sub-subnet is available under the
	// parent; resty does not surface that as an error, so check the status.
	if resp.StatusCode() != http.StatusOK || result.IP == "" {
		log.Printf("[WARN] datasourceD42SubnetSuggestedRead - no free sub-subnet available under parent %d", parentSubnetID)
		d.SetId("")
		return nil
	}

	d.Set("ip", result.IP)
	d.Set("mask_bits", result.Mask)
	d.SetId(result.IP)
	return nil
}
