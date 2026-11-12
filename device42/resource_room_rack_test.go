package device42

import (
	"encoding/json"
	"testing"
)

func TestRoomLookupResponseUnmarshal(t *testing.T) {
	body := `{"rooms":[{"room_id":35,"building_id":180,"building":"JETBI_QA","name":"tf-test-room","notes":"","tags":[],"custom_fields":[],"horizontal_grid_numbering":"Numeric","vertical_grid_numbering":"Alphabetic","grid_rows":0,"grid_cols":0,"reverse_xaxis":"no","reverse_yaxis":"no"}]}`

	var lookup apiRoomLookupResponse
	if err := json.Unmarshal([]byte(body), &lookup); err != nil {
		t.Fatalf("unmarshal room lookup: %s", err)
	}

	if len(lookup.Rooms) != 1 {
		t.Fatalf("expected exactly 1 room, got %d", len(lookup.Rooms))
	}

	r := lookup.Rooms[0]
	if r.RoomID != 35 {
		t.Fatalf("room_id=%d, want 35", r.RoomID)
	}
	if r.BuildingID != 180 {
		t.Fatalf("building_id=%d, want 180", r.BuildingID)
	}
	if r.Name != "tf-test-room" {
		t.Fatalf("name=%q, want \"tf-test-room\"", r.Name)
	}
	if r.Building != "JETBI_QA" {
		t.Fatalf("building=%q, want \"JETBI_QA\"", r.Building)
	}
	if r.GridRows != 0 {
		t.Fatalf("grid_rows=%d, want 0", r.GridRows)
	}
	if r.GridCols != 0 {
		t.Fatalf("grid_cols=%d, want 0", r.GridCols)
	}
}

func TestRackLookupResponseUnmarshal(t *testing.T) {
	body := `{"racks":[{"rack_id":131,"name":"tf-test-rack","size":42,"manufacturer":"","row":"","orientation":"left","notes":"","tags":[],"custom_fields":[]}]}`

	var lookup apiRackLookupResponse
	if err := json.Unmarshal([]byte(body), &lookup); err != nil {
		t.Fatalf("unmarshal rack lookup: %s", err)
	}

	if len(lookup.Racks) != 1 {
		t.Fatalf("expected exactly 1 rack, got %d", len(lookup.Racks))
	}

	r := lookup.Racks[0]
	if r.RackID != 131 {
		t.Fatalf("rack_id=%d, want 131", r.RackID)
	}
	if r.Name != "tf-test-rack" {
		t.Fatalf("name=%q, want \"tf-test-rack\"", r.Name)
	}
	if r.Size != 42 {
		t.Fatalf("size=%d, want 42", r.Size)
	}
	if r.Orientation != "left" {
		t.Fatalf("orientation=%q, want \"left\"", r.Orientation)
	}
}
