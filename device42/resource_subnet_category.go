package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiSubnetCategory struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Groups      string `json:"groups"`
}

type apiSubnetCategoryListResponse struct {
	Categories []apiSubnetCategory `json:"Categories"`
}

func resourceD42SubnetCategory() *schema.Resource {
	return &schema.Resource{
		Description: "device42_subnet_category can be use to manage subnet categories.",
		Create:      resourceDevice42SubnetCategoryCreate,
		Read:        resourceDevice42SubnetCategoryRead,
		Update:      resourceDevice42SubnetCategoryUpdate,
		Delete:      resourceDevice42SubnetCategoryDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCommaSeparatedOrderDiff,
			},
		},
	}
}

func resourceDevice42SubnetCategoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := subnetCategoryFormData(d, "")

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/subnet_category/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SubnetCategoryCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42SubnetCategoryCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42SubnetCategoryRead(d, m)
}

func resourceDevice42SubnetCategoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42SubnetCategoryRead - Starting reading using API for id %s", d.Id())

	var resp apiSubnetCategoryListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/subnet_category/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42SubnetCategoryRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42SubnetCategoryRead - invalid id %q: %s", d.Id(), err)
	}

	for _, c := range resp.Categories {
		if int(c.Id) == id {
			d.Set("name", c.Name)
			d.Set("description", c.Description)
			d.Set("groups", c.Groups)
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42SubnetCategoryRead - no subnet category found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42SubnetCategoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := subnetCategoryFormData(d, d.Id())

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/subnet_category/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SubnetCategoryUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42SubnetCategoryRead(d, m)
}

func resourceDevice42SubnetCategoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting subnet category id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/subnet_category/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42SubnetCategoryDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func subnetCategoryFormData(d *schema.ResourceData, id string) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if id != "" {
		formData["id"] = id
	}

	if v, ok := d.GetOk("description"); ok {
		formData["description"] = v.(string)
	}
	if v, ok := d.GetOk("groups"); ok {
		formData["groups"] = v.(string)
	}

	return formData
}
