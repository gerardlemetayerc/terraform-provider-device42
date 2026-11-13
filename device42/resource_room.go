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

type apiRoomReadResponse struct {
	RoomID                  int64         `json:"room_id"`
	BuildingID              int64         `json:"building_id"`
	Building                string        `json:"building"`
	Name                    string        `json:"name"`
	Notes                   string        `json:"notes"`
	Tags                    []interface{} `json:"tags"`
	HorizontalGridNumbering string        `json:"horizontal_grid_numbering"`
	VerticalGridNumbering   string        `json:"vertical_grid_numbering"`
	HorizontalGridStart     string        `json:"horizontal_grid_start"`
	VerticalGridStart       string        `json:"vertical_grid_start"`
	GridRows                int           `json:"grid_rows"`
	GridCols                int           `json:"grid_cols"`
	ReverseXaxis            string        `json:"reverse_xaxis"`
	ReverseYaxis            string        `json:"reverse_yaxis"`
	CustomFields            []customField `json:"custom_fields"`
}

func resourceD42Room() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42RoomCreate,
		Read:   resourceDevice42RoomRead,
		Update: resourceDevice42RoomUpdate,
		Delete: resourceDevice42RoomDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the room.",
			},
			"building_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the building that contains this room.",
			},
			"building": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the building that contains this room (alternative to building_id).",
			},
			"room_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the room.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"horizontal_grid_numbering": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Horizontal grid numbering: Numeric, Alphabetic or Alphabetic_doubled.",
			},
			"vertical_grid_numbering": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vertical grid numbering: Numeric, Alphabetic or Alphabetic_doubled.",
			},
			"horizontal_grid_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Horizontal grid start value.",
			},
			"vertical_grid_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vertical grid start value.",
			},
			"grid_rows": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of rows in the room grid.",
			},
			"grid_cols": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of columns in the room grid.",
			},
			"reverse_xaxis": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Reverse x-axis numbering (yes/no).",
			},
			"reverse_yaxis": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Reverse y-axis numbering (yes/no).",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a room (comma separated).",
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

				roomId, err := getRoomIdFromName(client, name)
				if err != nil {
					return nil, err
				}

				d.SetId(roomId)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42RoomCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	buildingID := d.Get("building_id").(int)
	if buildingID > 0 {
		formData["building_id"] = strconv.Itoa(buildingID)
	} else if v := d.Get("building").(string); v != "" {
		formData["building"] = v
	} else {
		return fmt.Errorf("either building_id or building must be set")
	}

	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("horizontal_grid_numbering").(string); v != "" {
		formData["horizontal_grid_numbering"] = v
	}
	if v := d.Get("vertical_grid_numbering").(string); v != "" {
		formData["vertical_grid_numbering"] = v
	}
	if v := d.Get("horizontal_grid_start").(string); v != "" {
		formData["horizontal_grid_start"] = v
	}
	if v := d.Get("vertical_grid_start").(string); v != "" {
		formData["vertical_grid_start"] = v
	}
	if v := d.Get("grid_rows").(int); v > 0 {
		formData["grid_rows"] = strconv.Itoa(v)
	}
	if v := d.Get("grid_cols").(int); v > 0 {
		formData["grid_cols"] = strconv.Itoa(v)
	}
	if v := d.Get("reverse_xaxis").(string); v != "" {
		formData["reverse_xaxis"] = v
	}
	if v := d.Get("reverse_yaxis").(string); v != "" {
		formData["reverse_yaxis"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/rooms/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42RoomCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42RoomCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42RoomCreate - unable to extract room id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("room_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42RoomCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/room/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42RoomCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42RoomRead(d, m)
}

func resourceDevice42RoomRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/rooms/%s/", d.Id())
	log.Printf("[DEBUG] resourceDevice42RoomRead - reading room on target URL: %s", d42Url)

	var resp apiRoomReadResponse
	_, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No room found: %s", d.Id())
		d.SetId("")
		return err
	}

	d.Set("name", resp.Name)
	d.Set("building_id", resp.BuildingID)
	d.Set("building", resp.Building)
	d.Set("room_id", resp.RoomID)
	d.Set("notes", resp.Notes)
	d.Set("horizontal_grid_numbering", resp.HorizontalGridNumbering)
	d.Set("vertical_grid_numbering", resp.VerticalGridNumbering)
	d.Set("horizontal_grid_start", resp.HorizontalGridStart)
	d.Set("vertical_grid_start", resp.VerticalGridStart)
	d.Set("grid_rows", resp.GridRows)
	d.Set("grid_cols", resp.GridCols)
	d.Set("reverse_xaxis", resp.ReverseXaxis)
	d.Set("reverse_yaxis", resp.ReverseYaxis)
	d.Set("tags", joinBuildingList(resp.Tags))
	d.Set("custom_fields", flattenCustomFields(resp.CustomFields))

	// The rooms single-GET endpoint omits custom_fields (unlike /rooms/?name=),
	// so enrich them from the list endpoint by name, matching on room_id.
	if name := d.Get("name").(string); name != "" {
		var lookup apiRoomLookupResponse
		if _, err := client.R().
			SetResult(&lookup).
			Get(fmt.Sprintf("/1.0/rooms/?name=%s", url.QueryEscape(name))); err == nil {
			for _, r := range lookup.Rooms {
				if r.RoomID == resp.RoomID {
					d.Set("custom_fields", flattenCustomFields(r.CustomFields))
					break
				}
			}
		}
	}

	return nil
}

func resourceDevice42RoomUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if d.HasChange("name") {
		// name is ForceNew so this should not occur, but guard anyway.
		formData["new_name"] = name
	}
	if d.HasChange("building_id") || d.HasChange("building") {
		if v := d.Get("building_id").(int); v > 0 {
			formData["building_id"] = strconv.Itoa(v)
		} else if v := d.Get("building").(string); v != "" {
			formData["building"] = v
		}
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("horizontal_grid_numbering") {
		formData["horizontal_grid_numbering"] = d.Get("horizontal_grid_numbering").(string)
	}
	if d.HasChange("vertical_grid_numbering") {
		formData["vertical_grid_numbering"] = d.Get("vertical_grid_numbering").(string)
	}
	if d.HasChange("horizontal_grid_start") {
		formData["horizontal_grid_start"] = d.Get("horizontal_grid_start").(string)
	}
	if d.HasChange("vertical_grid_start") {
		formData["vertical_grid_start"] = d.Get("vertical_grid_start").(string)
	}
	if d.HasChange("grid_rows") {
		formData["grid_rows"] = strconv.Itoa(d.Get("grid_rows").(int))
	}
	if d.HasChange("grid_cols") {
		formData["grid_cols"] = strconv.Itoa(d.Get("grid_cols").(int))
	}
	if d.HasChange("reverse_xaxis") {
		formData["reverse_xaxis"] = d.Get("reverse_xaxis").(string)
	}
	if d.HasChange("reverse_yaxis") {
		formData["reverse_yaxis"] = d.Get("reverse_yaxis").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}
	log.Printf("[DEBUG] resourceDevice42RoomUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/rooms/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42RoomUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42RoomUpdate - Result: %#v", r)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42RoomUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/room/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] resourceDevice42RoomUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42RoomRead(d, m)
}

func resourceDevice42RoomDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting room %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/rooms/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42RoomDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42RoomDelete - Result: %#v", r)
	return nil
}

func getRoomIdFromName(client *resty.Client, name string) (string, error) {
	var apiRoomLookUp struct {
		Rooms []struct {
			RoomID int64  `json:"room_id"`
			Name   string `json:"name"`
		} `json:"rooms"`
	}
	url := fmt.Sprintf("/1.0/rooms/?name=%s", url.QueryEscape(name))
	log.Printf("[DEBUG] Query for room name on URL : %s", url)
	_, err := client.R().
		SetResult(&apiRoomLookUp).
		Get(url)

	log.Printf("[DEBUG] Data from room lookup : %#v", apiRoomLookUp)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	if len(apiRoomLookUp.Rooms) == 0 {
		return "", fmt.Errorf("no room found with name: %s", name)
	}

	if len(apiRoomLookUp.Rooms) > 1 {
		return "", fmt.Errorf("more than one room found with name: %s", name)
	}

	return strconv.FormatInt(apiRoomLookUp.Rooms[0].RoomID, 10), nil
}
