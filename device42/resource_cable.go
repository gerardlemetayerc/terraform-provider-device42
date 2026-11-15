package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiCableResponse struct {
	Cable []apiCable `json:"cable"`
}

type apiCable struct {
	Id                int64         `json:"id"`
	CableId           string        `json:"cable_id"`
	OriginContentType string        `json:"origin_content_type"`
	OriginCableId     string        `json:"origin_cable_id"`
	OriginNetportId   int64         `json:"origin_netport_id"`
	OriginNetportName string        `json:"origin_netport_name"`
	OriginPatchPanel  string        `json:"origin_patch_panel"`
	CableLength       float64       `json:"cable_length"`
	CableLengthUnits  string        `json:"cable_length_units"`
	Notes             string        `json:"notes"`
	Vendor            string        `json:"vendor"`
	Tags              []interface{} `json:"tags"`
	CustomFields      []customField `json:"custom_fields"`
}

type apiCableSingleResponse struct {
	Cable apiCable `json:"cable"`
}

func resourceD42Cable() *schema.Resource {
	return &schema.Resource{
		Description: "device42_cable can be use to manage cables connecting two endpoints (patch panel ports, netports, circuits).",
		Create:      resourceDevice42CableCreate,
		Read:        resourceDevice42CableRead,
		Update:      resourceDevice42CableUpdate,
		Delete:      resourceDevice42CableDelete,

		Schema: map[string]*schema.Schema{
			"cable_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"origin_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"origin_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cable_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cable_length_units": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCommaSeparatedOrderDiff,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceDevice42CableCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := cableFormData(d, "")

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/cables/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CableCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42CableCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42CableRead(d, m)
}

func resourceDevice42CableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42CableRead - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiCableSingleResponse{}).
		Get(fmt.Sprintf("/1.0/cables/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No cable found for id %s", d.Id())
		d.SetId("")
		return err
	}
	if resp.StatusCode() == 404 {
		log.Printf("[WARN] No cable found for id %s", d.Id())
		d.SetId("")
		return nil
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42CableRead - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiCableSingleResponse).Cable
	d.Set("cable_id", r.CableId)
	d.Set("cable_length", int64(r.CableLength))
	d.Set("cable_length_units", r.CableLengthUnits)
	d.Set("notes", r.Notes)
	d.Set("vendor", r.Vendor)
	d.Set("tags", joinBuildingList(r.Tags))
	d.Set("custom_fields", flattenCustomFields(r.CustomFields))

	return nil
}

func resourceDevice42CableUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := cableFormData(d, d.Id())

	var changed bool
	if d.HasChange("cable_length") {
		changed = true
	}
	if d.HasChange("cable_length_units") {
		changed = true
	}
	if d.HasChange("notes") {
		changed = true
	}
	if d.HasChange("vendor") {
		changed = true
	}
	if d.HasChange("tags") {
		changed = true
	}

	if !changed {
		return resourceDevice42CableRead(d, m)
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/cables/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CableUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42CableRead(d, m)
}

func resourceDevice42CableDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting cable id %s (cable_id: %s)", d.Id(), d.Get("cable_id"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/cables/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42CableDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func cableFormData(d *schema.ResourceData, id string) map[string]string {
	formData := map[string]string{
		"cable_id": d.Get("cable_id").(string),
	}

	if id != "" {
		formData["id"] = id
	}

	if v, ok := d.GetOk("origin_type"); ok {
		formData["origin_type"] = v.(string)
	}
	if v, ok := d.GetOk("origin_id"); ok {
		formData["origin_id"] = v.(string)
	}
	if v, ok := d.GetOk("end_type"); ok {
		formData["end_type"] = v.(string)
	}
	if v, ok := d.GetOk("end_id"); ok {
		formData["end_id"] = v.(string)
	}
	if v, ok := d.GetOk("cable_length"); ok {
		formData["cable_length"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("cable_length_units"); ok {
		formData["cable_length_units"] = v.(string)
	}
	if v, ok := d.GetOk("notes"); ok {
		formData["notes"] = v.(string)
	}
	if v, ok := d.GetOk("vendor"); ok {
		formData["vendor"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
