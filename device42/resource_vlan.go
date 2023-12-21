package device42

import (
	"fmt"
	"log"
	"strconv"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiVlanReadReponse struct {
	VlanId      int64         `json:"vlan_id"`
	Number      int64         `json:"number"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        []interface{} `json:"tags"`
}

func resourceD42Vlans() *schema.Resource {
	return &schema.Resource{
		Description: "device42_vlan can be use to manage Business Applications",
		Create:      resourceDevice42VlanCreate,
		Read:        resourceDevice42VlanRead,
		Update:      resourceDevice42VlansUpdate,
		Delete:      resourceDevice42VlanDelete,

		Schema: map[string]*schema.Schema{
			"number": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VLAN ID Number.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If left blank, name will be created as VLANxxxx, e.g. VLAN# 342 will be named VLAN0342",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a VLAN",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes",
			},
		},
	}
}

func resourceDevice42VlanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42VlanRead - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiVlanReadReponse{}).
		Get(fmt.Sprintf("/1.0/vlans/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No vlans found for id %s", d.Id())
		d.SetId("")
		return err
	}

	r := resp.Result().(*apiVlanReadReponse)
	str := fmt.Sprintf("%v", r)
	log.Printf("[DEBUG] resourceDevice42VlanRead - API data %s", str)
	d.Set("number", r.Number)
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("tags", r.Tags)

	return nil
}

func resourceDevice42VlanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting vlan %s (UUID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/vlans/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(url)

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func resourceDevice42VlansUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := map[string]string{
		"id": d.Id(),
	}

	if d.HasChange("name") {
		formData["name"] = d.Get("name").(string)
	}
	if d.HasChange("number") {
		formData["number"] = d.Get("number").(string)
	}
	if d.HasChange("description") {
		formData["description"] = d.Get("description").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}

	if len(formData) > 1 {
		resp, err := client.R().
			SetFormData(formData).
			SetResult(apiResponse{}).
			Put(fmt.Sprintf("/1.0/vlans/%s/", d.Id()))

		if err != nil {
			return err
		}
		r := resp.Result().(*apiResponse)
		log.Printf("[DEBUG] Result: %#v", r)
	}

	return resourceDevice42VlanRead(d, m)
}

func resourceDevice42VlanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	number := d.Get("number").(string)
	description := d.Get("description").(string)

	resp, err := client.R().
		SetFormData(map[string]string{
			"name":        name,
			"number":      number,
			"description": description,
		}).
		SetResult(apiResponse{}).
		Post("/1.0/vlans/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	// Set ID after vlan creation
	d.SetId(strconv.Itoa(id))

	return resourceDevice42VlanRead(d, m)
}
