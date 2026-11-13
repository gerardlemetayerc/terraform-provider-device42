package device42

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiBuildingReadResponse struct {
	BuildingID   int64             `json:"building_id"`
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	Address      string            `json:"address"`
	ContactName  string            `json:"contact_name"`
	ContactPhone string            `json:"contact_phone"`
	Latitude     apiStringOrNumber `json:"latitude"`
	Longitude    apiStringOrNumber `json:"longitude"`
	Notes        string            `json:"notes"`
	Tags         []interface{}     `json:"tags"`
	Groups       []interface{}     `json:"groups"`
	CustomFields []customField     `json:"custom_fields"`
}

// apiStringOrNumber tolerates Device42 returning a field as a JSON string, a
// JSON number, or null (e.g. building latitude/longitude, which the API emits
// as numbers like 48.85 or null). The value is stored canonically as a string.
type apiStringOrNumber string

func (s *apiStringOrNumber) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*s = ""
		return nil
	}
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*s = apiStringOrNumber(str)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(b, &num); err == nil {
		*s = apiStringOrNumber(num.String())
		return nil
	}
	return fmt.Errorf("cannot unmarshal value as string or number: %s", string(b))
}

// apiDeleteResponse matches the shape returned by DELETE endpoints such as
// /1.0/buildings/{id}/ which return {"deleted":"true","id":"114"} instead of
// the {code,msg} shape used by the rest of the API.
type apiDeleteResponse struct {
	Deleted string `json:"deleted"`
	ID      string `json:"id"`
}

func resourceD42Building() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42BuildingCreate,
		Read:   resourceDevice42BuildingRead,
		Update: resourceDevice42BuildingUpdate,
		Delete: resourceDevice42BuildingDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the building.",
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the building.",
			},
			"building_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the building.",
			},
			"contact_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact name of the building.",
			},
			"contact_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact phone of the building.",
			},
			"latitude": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "North-south position.",
			},
			"longitude": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "East-west position.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a building (comma separated).",
			},
			"groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Admin groups that have access to this building (multitenancy), e.g. Prod_East:no,Corp:yes.",
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

				buildingId, err := getBuildingIdFromName(client, name)
				if err != nil {
					return nil, err
				}

				d.SetId(buildingId)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42BuildingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if v := d.Get("address").(string); v != "" {
		formData["address"] = v
	}
	if v := d.Get("contact_name").(string); v != "" {
		formData["contact_name"] = v
	}
	if v := d.Get("contact_phone").(string); v != "" {
		formData["contact_phone"] = v
	}
	if v := d.Get("latitude").(string); v != "" {
		formData["latitude"] = v
	}
	if v := d.Get("longitude").(string); v != "" {
		formData["longitude"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}
	if v := d.Get("groups").(string); v != "" {
		formData["groups"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/buildings/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42BuildingCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42BuildingCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42BuildingCreate - unable to extract building id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("building_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42BuildingCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/building/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42BuildingCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42BuildingRead(d, m)
}

func resourceDevice42BuildingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/buildings/%s/", d.Id())
	log.Printf("[DEBUG] resourceDevice42BuildingRead - reading building on target URL: %s", d42Url)

	var resp apiBuildingReadResponse
	_, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No building found: %s", d.Id())
		d.SetId("")
		return err
	}

	d.Set("name", resp.Name)
	d.Set("building_id", resp.BuildingID)
	d.Set("address", resp.Address)
	d.Set("contact_name", resp.ContactName)
	d.Set("contact_phone", resp.ContactPhone)
	d.Set("latitude", string(resp.Latitude))
	d.Set("longitude", string(resp.Longitude))
	d.Set("notes", resp.Notes)
	d.Set("tags", joinBuildingList(resp.Tags))
	d.Set("groups", joinBuildingList(resp.Groups))
	d.Set("custom_fields", flattenCustomFields(resp.CustomFields))

	return nil
}

func resourceDevice42BuildingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if d.HasChange("name") {
		// name is ForceNew so this should not occur, but guard anyway.
		formData["new_name"] = name
	}
	if d.HasChange("address") {
		formData["address"] = d.Get("address").(string)
	}
	if d.HasChange("contact_name") {
		formData["contact_name"] = d.Get("contact_name").(string)
	}
	if d.HasChange("contact_phone") {
		formData["contact_phone"] = d.Get("contact_phone").(string)
	}
	if d.HasChange("latitude") {
		formData["latitude"] = d.Get("latitude").(string)
	}
	if d.HasChange("longitude") {
		formData["longitude"] = d.Get("longitude").(string)
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}
	if d.HasChange("groups") {
		formData["groups"] = d.Get("groups").(string)
	}
	log.Printf("[DEBUG] resourceDevice42BuildingUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/buildings/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42BuildingUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42BuildingUpdate - Result: %#v", r)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42BuildingUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/building/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] resourceDevice42BuildingUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42BuildingRead(d, m)
}

func resourceDevice42BuildingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting building %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/buildings/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42BuildingDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42BuildingDelete - Result: %#v", r)
	return nil
}

func getBuildingIdFromName(client *resty.Client, name string) (string, error) {
	var apiBuildingLookUp struct {
		Buildings []struct {
			BuildingID int64  `json:"building_id"`
			Name       string `json:"name"`
		} `json:"buildings"`
	}
	url := fmt.Sprintf("/1.0/buildings/?name=%s", url.QueryEscape(name))
	log.Printf("[DEBUG] Query for building name on URL : %s", url)
	_, err := client.R().
		SetResult(&apiBuildingLookUp).
		Get(url)

	log.Printf("[DEBUG] Data from building lookup : %#v", apiBuildingLookUp)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	if len(apiBuildingLookUp.Buildings) == 0 {
		return "", fmt.Errorf("no building found with name: %s", name)
	}

	if len(apiBuildingLookUp.Buildings) > 1 {
		return "", fmt.Errorf("more than one building found with name: %s", name)
	}

	return strconv.FormatInt(apiBuildingLookUp.Buildings[0].BuildingID, 10), nil
}

// msgID extracts the new object ID from a Device42 {code,msg} response. The ID
// is expected at msg index 1 but elements can be decoded as float64 (default)
// or json.Number (when UseNumber is enabled), so a bare r.Msg[1].(float64)
// would panic. This helper is panic-safe.
func msgID(msg []interface{}) (int, error) {
	if len(msg) < 2 {
		return 0, fmt.Errorf("unexpected msg length: %d", len(msg))
	}
	switch v := msg[1].(type) {
	case float64:
		return int(v), nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("cannot parse building id %q: %s", v, err)
		}
		return int(n), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("unexpected type for building id: %T (%v)", msg[1], msg[1])
	}
}

// joinBuildingList renders a Device42 tags/groups array (e.g. ["a","b"]) as a
// comma-separated string suitable for the TypeString schema fields.
func joinBuildingList(in []interface{}) string {
	if len(in) == 0 {
		return ""
	}
	parts := make([]string, 0, len(in))
	for _, v := range in {
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, ",")
}
