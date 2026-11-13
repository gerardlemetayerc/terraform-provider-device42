package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func resourceD42SubnetSuggested() *schema.Resource {
	return &schema.Resource{
		Create: resourceD42SubnetSuggestedCreate,
		Read:   resourceD42SubnetSuggestedRead,
		Delete: resourceD42SubnetSuggestedDelete,

		Description: "Suggests and reserves the next available sub-subnet of a parent subnet. The suggestion is consumed by creating the subnet in Device42, so repeated applies stay stable.",
		Schema: map[string]*schema.Schema{
			"parent_subnet_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the parent subnet to suggest a child from.",
			},
			"mask_bits": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Mask bits of the suggested sub-subnet.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name to give the created sub-subnet.",
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network address of the suggested (and created) sub-subnet.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the created sub-subnet.",
			},
		},
	}
}

func resourceD42SubnetSuggestedCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	parentSubnetID := d.Get("parent_subnet_id").(int)
	maskBits := d.Get("mask_bits").(int)

	// 1. Ask Device42 for the next available sub-subnet under the parent.
	suggestPath := fmt.Sprintf("/1.0/suggest_subnet/%d/?mask_bits=%d", parentSubnetID, maskBits)
	log.Printf("[DEBUG] resourceD42SubnetSuggestedCreate - suggesting via %s", suggestPath)

	var suggested apiSuggestSubnetResponse
	resp, err := client.R().
		SetResult(&suggested).
		Get(suggestPath)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK || suggested.IP == "" {
		return fmt.Errorf("no free sub-subnet available under parent subnet %d with mask_bits %d", parentSubnetID, maskBits)
	}

	// 2. Reserve the suggestion by creating the subnet, so the value is consumed.
	createForm := map[string]string{
		"network":          suggested.IP,
		"mask_bits":        strconv.Itoa(suggested.Mask),
		"parent_subnet_id": strconv.Itoa(parentSubnetID),
	}
	if name := d.Get("name").(string); name != "" {
		createForm["name"] = name
	}

	log.Printf("[DEBUG] resourceD42SubnetSuggestedCreate - creating subnet network=%s mask=%d", suggested.IP, suggested.Mask)
	createResp, err := client.R().
		SetFormData(createForm).
		SetResult(apiResponse{}).
		Post("/1.0/subnets/")
	if err != nil {
		return err
	}

	r := createResp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceD42SubnetSuggestedCreate - api returned code %d", r.Code)
	}
	if len(r.Msg) < 2 {
		return fmt.Errorf("incorrect response to subnet creation")
	}

	id := int(r.Msg[1].(float64))
	d.SetId(strconv.Itoa(id))
	return resourceD42SubnetSuggestedRead(d, m)
}

func resourceD42SubnetSuggestedRead(d *schema.ResourceData, m interface{}) error {
	readResponse, err := apiDevice42Get(m.(*resty.Client), fmt.Sprintf("/1.0/subnets/%s/", d.Id()), apiSubnetReadResponse{})
	if err != nil {
		log.Printf("[WARN] No subnet found: %s", d.Id())
		d.SetId("")
		return nil
	}

	r := readResponse.Result().(*apiSubnetReadResponse)
	log.Printf("[DEBUG] resourceD42SubnetSuggestedRead - API data %#v", r)

	d.Set("network", r.Network)
	d.Set("mask_bits", r.MaskBits)
	d.Set("subnet_id", r.Subnet_id)
	if r.Name != "" {
		d.Set("name", r.Name)
	}
	return nil
}

func resourceD42SubnetSuggestedDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting suggested subnet %s (id: %s)", d.Get("name"), d.Id())

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/subnets/%s/", d.Id()))
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] resourceD42SubnetSuggestedDelete - Result: %#v", r)
	return nil
}
