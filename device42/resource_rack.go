package device42

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiRackReadResponse struct {
	RackID       int64         `json:"rack_id"`
	Name         string        `json:"name"`
	Size         int           `json:"size"`
	Building     string        `json:"building"`
	Room         string        `json:"room"`
	Manufacturer string        `json:"manufacturer"`
	Notes        string        `json:"notes"`
	Row          string        `json:"row"`
	Orientation  string        `json:"orientation"`
	Tags         []interface{} `json:"tags"`
	CustomFields []customField `json:"custom_fields"`
}

func resourceD42Rack() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42RackCreate,
		Read:   resourceDevice42RackRead,
		Update: resourceDevice42RackUpdate,
		Delete: resourceDevice42RackDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the rack.",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Size of the rack (number of U).",
			},
			"rack_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the rack.",
			},
			"building": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the building that contains this rack.",
			},
			"room": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the room that contains this rack.",
			},
			"room_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the room that contains this rack (used when room name is not unique).",
			},
			"manufacturer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rack manufacturer.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"row": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Row name of the rack.",
			},
			"orientation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Orientation of the rack: left, right, up or down.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a rack (comma separated).",
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
		},

		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*resty.Client)
				name := d.Id()

				rackId, err := getRackIdFromName(client, name)
				if err != nil {
					return nil, err
				}

				d.SetId(rackId)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42RackCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name
	formData["size"] = strconv.Itoa(d.Get("size").(int))

	if v := d.Get("building").(string); v != "" {
		formData["building"] = v
	}
	if v := d.Get("room").(string); v != "" {
		formData["room"] = v
	}
	if v := d.Get("room_id").(int); v > 0 {
		formData["room_id"] = strconv.Itoa(v)
	}
	if v := d.Get("manufacturer").(string); v != "" {
		formData["manufacturer"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("row").(string); v != "" {
		formData["row"] = v
	}
	if v := d.Get("orientation").(string); v != "" {
		formData["orientation"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/racks/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42RackCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42RackCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42RackCreate - unable to extract rack id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("rack_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42RackCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/rack/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42RackCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42RackRead(d, m)
}

func resourceDevice42RackRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/racks/%s/", d.Id())
	log.Printf("[DEBUG] resourceDevice42RackRead - reading rack on target URL: %s", d42Url)

	var resp apiRackReadResponse
	_, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No rack found: %s", d.Id())
		d.SetId("")
		return err
	}

	d.Set("name", resp.Name)
	d.Set("size", resp.Size)
	d.Set("rack_id", resp.RackID)
	d.Set("building", resp.Building)
	d.Set("room", resp.Room)
	d.Set("manufacturer", resp.Manufacturer)
	d.Set("notes", resp.Notes)
	d.Set("row", resp.Row)
	d.Set("orientation", resp.Orientation)
	d.Set("tags", joinBuildingList(resp.Tags))
	d.Set("custom_fields", flattenCustomFields(resp.CustomFields))

	return nil
}

func resourceDevice42RackUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name
	formData["size"] = strconv.Itoa(d.Get("size").(int))

	if d.HasChange("name") {
		// name is ForceNew so this should not occur, but guard anyway.
		formData["new_name"] = name
	}
	if d.HasChange("size") {
		formData["size"] = strconv.Itoa(d.Get("size").(int))
	}
	if d.HasChange("building") {
		formData["building"] = d.Get("building").(string)
	}
	if d.HasChange("room") {
		formData["room"] = d.Get("room").(string)
	}
	if d.HasChange("room_id") {
		if v := d.Get("room_id").(int); v > 0 {
			formData["room_id"] = strconv.Itoa(v)
		}
	}
	if d.HasChange("manufacturer") {
		formData["manufacturer"] = d.Get("manufacturer").(string)
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("row") {
		formData["row"] = d.Get("row").(string)
	}
	if d.HasChange("orientation") {
		formData["orientation"] = d.Get("orientation").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}
	log.Printf("[DEBUG] resourceDevice42RackUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/racks/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42RackUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42RackUpdate - Result: %#v", r)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42RackUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/rack/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			if resp.IsError() || r.Code != 0 {
				return fmt.Errorf("resourceDevice42RackUpdate - custom field %q update failed (HTTP %d, API code %d)", k, resp.StatusCode(), r.Code)
			}
			log.Printf("[DEBUG] resourceDevice42RackUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42RackRead(d, m)
}

func resourceDevice42RackDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting rack %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/racks/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42RackDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42RackDelete - Result: %#v", r)
	return nil
}

func getRackIdFromName(client *resty.Client, name string) (string, error) {
	var apiRackLookUp struct {
		Racks []struct {
			RackID int64  `json:"rack_id"`
			Name   string `json:"name"`
		} `json:"racks"`
	}
	url := fmt.Sprintf("/1.0/racks/?name=%s", url.QueryEscape(name))
	log.Printf("[DEBUG] Query for rack name on URL : %s", url)
	_, err := client.R().
		SetResult(&apiRackLookUp).
		Get(url)

	log.Printf("[DEBUG] Data from rack lookup : %#v", apiRackLookUp)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	if len(apiRackLookUp.Racks) == 0 {
		return "", fmt.Errorf("no rack found with name: %s", name)
	}

	if len(apiRackLookUp.Racks) > 1 {
		return "", fmt.Errorf("more than one rack found with name: %s", name)
	}

	return strconv.FormatInt(apiRackLookUp.Racks[0].RackID, 10), nil
}
