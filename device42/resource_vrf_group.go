package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiVrfGroup struct {
	Id           int64         `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Tags         []interface{} `json:"tags"`
	CustomFields []customField `json:"custom_fields"`
}

type apiVrfGroupSingleResponse struct {
	VrfGroup apiVrfGroup `json:"vrfgroup"`
}

func resourceD42VrfGroup() *schema.Resource {
	return &schema.Resource{
		Description: "device42_vrf_group can be use to manage VRF groups.",
		Create:      resourceDevice42VrfGroupCreate,
		Read:        resourceDevice42VrfGroupRead,
		Update:      resourceDevice42VrfGroupUpdate,
		Delete:      resourceDevice42VrfGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressTagsAppendDiff,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceDevice42VrfGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := vrfGroupFormData(d, "")

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/vrfgroup/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42VrfGroupCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42VrfGroupCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42VrfGroupRead(d, m)
}

func resourceDevice42VrfGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42VrfGroupRead - Starting reading using API for id %s", d.Id())

	var resp apiVrfGroupSingleResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get(fmt.Sprintf("/1.0/vrfgroup/%s/", d.Id()))
	if err != nil {
		log.Printf("[WARN] No vrf group found for id %s", d.Id())
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() == 404 {
		log.Printf("[WARN] No vrf group found for id %s", d.Id())
		d.SetId("")
		return nil
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42VrfGroupRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	g := resp.VrfGroup
	d.Set("name", g.Name)
	d.Set("description", g.Description)
	d.Set("tags", joinBuildingList(g.Tags))
	d.Set("custom_fields", flattenCustomFields(g.CustomFields))

	return nil
}

func resourceDevice42VrfGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	payload := map[string]string{
		"name": d.Get("name").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		payload["description"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		payload["tags"] = v.(string)
	}

	resp, err := client.R().
		SetResult(apiResponse{}).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put(fmt.Sprintf("/1.0/vrfgroup/%s/", d.Id()))

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42VrfGroupUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42VrfGroupRead(d, m)
}

func resourceDevice42VrfGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting vrf group id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/vrfgroup/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42VrfGroupDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func vrfGroupFormData(d *schema.ResourceData, id string) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if id != "" {
		formData["id"] = id
	}

	if v, ok := d.GetOk("description"); ok {
		formData["description"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
