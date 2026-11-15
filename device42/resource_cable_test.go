package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCableFormDataCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42Cable().Schema, map[string]interface{}{
		"cable_id":           "myCable",
		"origin_type":        "netport",
		"origin_id":          "1016",
		"end_type":           "netport",
		"end_id":             "1017",
		"cable_length":       7,
		"cable_length_units": "m",
		"notes":              "a note",
		"vendor":             "Cisco",
		"tags":               "prod,eth",
	})

	got := cableFormData(d, "")

	expected := map[string]string{
		"cable_id":           "myCable",
		"origin_type":        "netport",
		"origin_id":          "1016",
		"end_type":           "netport",
		"end_id":             "1017",
		"cable_length":       "7",
		"cable_length_units": "m",
		"notes":              "a note",
		"vendor":             "Cisco",
		"tags":               "prod,eth",
	}

	if len(got) != len(expected) {
		t.Fatalf("cableFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("cableFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}

	if _, ok := got["id"]; ok {
		t.Errorf("cableFormData() on create must not include an id")
	}
}

func TestCableFormDataUpdateIncludesID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42Cable().Schema, map[string]interface{}{
		"cable_id":     "myCable",
		"cable_length": 12,
	})

	got := cableFormData(d, "42")

	if got["id"] != "42" {
		t.Errorf("cableFormData() on update must include id, got %q", got["id"])
	}
	if got["cable_length"] != "12" {
		t.Errorf("cableFormData() cable_length = %q, want %q", got["cable_length"], "12")
	}
}

func TestCableFormDataOmitsEmptyOptionals(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42Cable().Schema, map[string]interface{}{
		"cable_id": "myCable",
		"notes":    "",
		"vendor":   "",
	})

	got := cableFormData(d, "")

	if _, ok := got["notes"]; ok {
		t.Errorf("cableFormData() should omit empty notes")
	}
	if _, ok := got["vendor"]; ok {
		t.Errorf("cableFormData() should omit empty vendor")
	}
	if got["cable_id"] != "myCable" {
		t.Errorf("cableFormData() cable_id = %q, want %q", got["cable_id"], "myCable")
	}
}
