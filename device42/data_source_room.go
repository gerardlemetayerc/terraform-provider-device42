package device42

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiRoomLookupResponse is the shape returned by GET /1.0/rooms/?name=
// which is a bare array under the "rooms" key (no total_count/offset).
type apiRoomLookupResponse struct {
	Rooms []apiRoomReadResponse `json:"rooms"`
}

func datasourceD42Room() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42RoomRead,
		Description: "Read a Device42 room.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the room.",
			},
			"room_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the room.",
			},
			"building_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"building": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"horizontal_grid_numbering": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vertical_grid_numbering": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"horizontal_grid_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vertical_grid_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"grid_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"grid_cols": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"reverse_xaxis": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse_yaxis": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func datasourceD42RoomRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	roomID := d.Get("room_id").(int)

	if name != "" {
		url := fmt.Sprintf("/1.0/rooms/?name=%s", url.QueryEscape(name))
		log.Printf("[DEBUG] datasourceD42RoomRead - Query rooms on URL : %s", url)

		var lookup apiRoomLookupResponse
		_, err := client.R().
			SetResult(&lookup).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No room found: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Rooms) == 0 {
			log.Printf("[ERROR] No room found with name: %s", name)
			d.SetId("")
			return nil
		}

		if len(lookup.Rooms) > 1 {
			log.Printf("[ERROR] More than one room found: %d", len(lookup.Rooms))
			d.SetId("")
			return nil
		}

		r := lookup.Rooms[0]
		d.SetId(strconv.FormatInt(r.RoomID, 10))
		d.Set("room_id", r.RoomID)
		d.Set("name", r.Name)
		d.Set("building_id", r.BuildingID)
		d.Set("building", r.Building)
		d.Set("notes", r.Notes)
		d.Set("horizontal_grid_numbering", r.HorizontalGridNumbering)
		d.Set("vertical_grid_numbering", r.VerticalGridNumbering)
		d.Set("horizontal_grid_start", r.HorizontalGridStart)
		d.Set("vertical_grid_start", r.VerticalGridStart)
		d.Set("grid_rows", r.GridRows)
		d.Set("grid_cols", r.GridCols)
		d.Set("reverse_xaxis", r.ReverseXaxis)
		d.Set("reverse_yaxis", r.ReverseYaxis)
		d.Set("tags", joinBuildingList(r.Tags))
		d.Set("custom_fields", flattenCustomFields(r.CustomFields))
		return nil
	}

	if roomID > 0 {
		url := fmt.Sprintf("/1.0/rooms/%d/", roomID)
		log.Printf("[DEBUG] datasourceD42RoomRead - Query room on URL : %s", url)

		var resp apiRoomReadResponse
		_, err := client.R().
			SetResult(&resp).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No room found: %d", roomID)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.Itoa(roomID))
		d.Set("room_id", resp.RoomID)
		d.Set("name", resp.Name)
		d.Set("building_id", resp.BuildingID)
		d.Set("building", resp.Building)
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
		return nil
	}

	return fmt.Errorf("either 'name' or 'room_id' must be set to look up a Device42 room")
}
